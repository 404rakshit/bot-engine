package users

import (
	userModels "bot-engine/models/mongo/users"
	botServices "bot-engine/services/bot"
	userServices "bot-engine/services/users"
	"bot-engine/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userService = userServices.UserService
type botService = botServices.BotService // Added this contract
type userModel = userModels.User

// CreateBotRequest DTO ensures we only accept exactly what we need
type CreateBotRequest struct {
	Token string `json:"token" binding:"required"`
}

type UserHandler struct {
	UserService userService
	BotService  botService // Inject this via Wire
}

func NewUserHandler(userService userService, botService botService) *UserHandler {
	return &UserHandler{
		UserService: userService,
		BotService:  botService,
	}
}

func (h *UserHandler) List(c *gin.Context) {

	users, err := h.UserService.List(c)

	if err != nil {
		response := utils.ErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if users == nil {
		users = make([]userModel, 0)
	}

	response := utils.SuccessResponse(users, "Successfully fetched Users")
	c.JSON(http.StatusOK, response)

}

func (h *UserHandler) Create(c *gin.Context) {

	var user userModel

	if err := c.ShouldBindJSON(&user); err != nil {
		response := utils.ErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := h.UserService.Create(c, &user); err != nil {
		response := utils.ErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.SuccessResponse(user, "Created Successfully")
	c.JSON(http.StatusOK, response)

}

func (h *UserHandler) ListUserBots(c *gin.Context) {
	// 1. Identify the tenant.
	// Assuming an Auth Middleware sets "user_id" in the Gin context.
	// If you use URL params (e.g., /users/:id/bots), use c.Param("id") instead.
	ownerID := c.Param("user_id")

	if ownerID == "" {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized user"))
		return
	}

	// 2. Fetch only this specific user's bots
	bots, err := h.BotService.GetBotsByOwnerID(c.Request.Context(), ownerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	// 3. Ensure we return an empty array `[]` in JSON instead of `null`
	if bots == nil {
		// Assuming botServices.Bot is your returned entity type
		bots = make([]botServices.Bot, 0)
	}

	response := utils.SuccessResponse(bots, "Successfully fetched user bots")
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) CreateUserBot(c *gin.Context) {
	// 1. Identify the tenant making the request
	ownerID := c.Param("user_id")
	if ownerID == "" {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized user"))
		return
	}

	// 2. Bind and validate the JSON payload using our DTO
	var req CreateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response := utils.ErrorResponse("Invalid payload: telegram token is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 3. Pass the data to the domain layer (orchestrating verification & encryption)
	createdBot, err := h.BotService.ConnectNewBot(c.Request.Context(), ownerID, req.Token)

	if err != nil {
		// You might want to parse specific domain errors here (e.g., 400 for invalid token vs 500 for db failure)
		response := utils.ErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 4. Return the successfully created bot entity (without the raw token)
	response := utils.SuccessResponse(createdBot, "Bot successfully connected")
	c.JSON(http.StatusCreated, response) // 201 Created is the standard here
}
