package auth

import (
	"bot-engine/helper"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	authModel "bot-engine/models/mongo/auth"
	userModels "bot-engine/models/mongo/users"
)

type TelegramAuthPayload = authModel.TelegramAuthPayload

func (s *authService) ProcessTelegramLogin(ctx context.Context, tgPayload TelegramAuthPayload) (string, error, string) {
	// Convert Telegram ID to string for the provider table
	providerUserID := strconv.FormatInt(tgPayload.ID, 10)

	// 1. Check if this Telegram ID is already linked to a user
	identity, err := s.identityRepo.GetByProvider(ctx, "telegram", providerUserID)
	if err == nil && identity != nil {
		// User exists! Generate your App JWT using their master User ID
		token, err := helper.GenerateToken(identity.UserID.Hex(), "", s.jwtSecret)
		return token, err, identity.ID.Hex()
	}

	// 2. User does not exist.
	// Create the master User record first.
	// Since Telegram doesn't guarantee an email, we leave it blank or handle it as optional.
	newUser := &userModels.User{
		Name: strings.TrimSpace(tgPayload.FirstName + " " + tgPayload.LastName),
		// CreatedAt: time.Now(),
		// Email is empty, PasswordHash is empty
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return "", fmt.Errorf("failed to create user"), ""
	}

	// 3. Create the Identity linkage
	newIdentity := &userModels.Identity{
		UserID:         newUser.ID,
		Provider:       "telegram",
		ProviderUserID: providerUserID,
		LinkedAt:       time.Now(),
	}

	if err := s.identityRepo.Create(ctx, newIdentity); err != nil {
		return "", fmt.Errorf("failed to link telegram identity"), ""
	}

	token, err := helper.GenerateToken(newUser.ID.Hex(), "", s.jwtSecret)

	// 4. Generate the JWT
	return token, err, newUser.ID.Hex()
}
