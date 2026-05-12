package engine

import (
	"context"
	"fmt"
	"log"
	"sync"

	engineModels "bot-engine/models/mongo/engine"
	engineRepos "bot-engine/repositories/engine"
	"bot-engine/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EngineService interface {
	ProcessIncomingMessage(ctx context.Context, botID bson.ObjectID, chatID int64, message string) ([]string, error)
}

type engineService struct {
	repo          engineRepos.EngineRepository
	workflowCache sync.Map // In-memory workflow cache (workflowID -> *Workflow)
	maxSteps      int      // Infinite loop circuit breaker
}

func NewEngineService(repo engineRepos.EngineRepository) EngineService {
	return &engineService{
		repo:     repo,
		maxSteps: 20,
	}
}

func (s *engineService) getCachedWorkflow(ctx context.Context, id bson.ObjectID) (*engineModels.Workflow, error) {
	if val, ok := s.workflowCache.Load(id); ok {
		return val.(*engineModels.Workflow), nil
	}

	// Cache Miss: Query MongoDB
	wf, err := s.repo.GetWorkflow(ctx, id)
	if err != nil {
		return nil, err
	}

	s.workflowCache.Store(id, wf)
	return wf, nil
}

func (s *engineService) ProcessIncomingMessage(ctx context.Context, botID bson.ObjectID, chatID int64, message string) ([]string, error) {
	// 1. Retrieve session from MongoDB
	session, err := s.repo.GetSession(ctx, botID, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user session: %v", err)
	}

	// 2. If no session exists, initialize a new default session
	// (Assumes a system-level mapping connects BotID to a Master WorkflowID)
	if session == nil {
		// Mock workflow mapping retrieval: In production, lookup the active workflow tied to this BotID
		session = &engineModels.UserSession{
			BotID:         botID,
			ChatID:        chatID,
			WorkflowID:    bson.NilObjectID, // Bind your actual Workflow ID here
			CurrentNodeID: "",               // Will initiate at StartNodeID
			Context:       make(map[string]interface{}),
		}
	}

	// 3. Fetch the Blueprint
	workflow, err := s.getCachedWorkflow(ctx, session.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve workflow blueprint: %v", err)
	}

	currentNodeID := session.CurrentNodeID
	if currentNodeID == "" {
		currentNodeID = workflow.StartNodeID
	}

	var messagesToSend []string
	currentNode := workflow.Nodes[currentNodeID]

	// 4. Process Pending Input
	if currentNode.Type == engineModels.NodeInput && currentNode.VariableName != "" {
		session.Context[currentNode.VariableName] = utils.CastInput(message)
		currentNodeID = currentNode.Next
	}

	iterations := 0
	isWaitingForInput := false

	// 5. Non-blocking FSM Loop (Hops through trigger, message, and condition nodes instantly)
	for currentNodeID != "" && !isWaitingForInput {
		iterations++
		if iterations > s.maxSteps {
			log.Printf("[FSM Guardrail Triggered] Potential infinite loop in workflow %s", session.WorkflowID.Hex())
			messagesToSend = append(messagesToSend, "Something went wrong. Please try again later.")
			currentNodeID = "" // Break out safely
			break
		}

		node, exists := workflow.Nodes[currentNodeID]
		if !exists {
			break
		}

		switch node.Type {
		case engineModels.NodeTrigger:
			currentNodeID = node.Next

		case engineModels.NodeMessage:
			if node.Content != "" {
				processedText := utils.InterpolateVariables(node.Content, session.Context)
				messagesToSend = append(messagesToSend, processedText)
			}
			currentNodeID = node.Next

		case engineModels.NodeCondition:
			if node.Expression != "" && node.Branches != nil {
				conditionMet := utils.EvaluateExpression(node.Expression, session.Context)
				if conditionMet {
					currentNodeID = node.Branches["true"]
				} else {
					currentNodeID = node.Branches["false"]
				}
			} else {
				currentNodeID = ""
			}

		case engineModels.NodeInput:
			// Hit input node, pause engine execution and wait for the next user event
			isWaitingForInput = true
			session.CurrentNodeID = currentNodeID

		default:
			currentNodeID = ""
		}
	}

	if !isWaitingForInput {
		// Flow reached an end node
		session.CurrentNodeID = ""
	}

	// 6. Write state updates back to MongoDB (A single upsert)
	if err := s.repo.SaveSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save user session: %v", err)
	}

	return messagesToSend, nil
}
