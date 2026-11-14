package request

import (
	"errors"
	"pet/pkg/constant"

	"github.com/gin-gonic/gin"
)

func GetUserIdFromRequest(c *gin.Context) (int, error) {
	value, ok := c.Get(constant.UserId)
	if !ok {
		return 0, errors.New("user id not found")
	}

	userId, ok := value.(int)
	if !ok {
		return 0, errors.New("invalid type of user id")
	}

	return userId, nil
}
