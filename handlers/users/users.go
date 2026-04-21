package users

import (
	userModels "di/models/mongo/users"
	userServices "di/services/users"
	"di/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userService = userServices.UserService
type userModel = userModels.User

type UserHandler struct {
	UserService userService
}

func NewUserHandler(userService userService) *UserHandler {
	return &UserHandler{
		UserService: userService,
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
