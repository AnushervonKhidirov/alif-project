package helper

import (
	"errors"
	"pet/internal/constant"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromRequest(c *gin.Context) (int, error) {
	value, ok := c.Get(constant.UserIdKey)
	if !ok {
		return 0, errors.New("user id not found")
	}

	userId, ok := value.(int)
	if !ok {
		return 0, errors.New("invalid type of user id")
	}

	return userId, nil
}
