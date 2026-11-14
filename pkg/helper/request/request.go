package request

import (
	"errors"
	"pet/pkg/constant"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit  int
	Offset int
}

func GetPaginationQuery(c *gin.Context) (Pagination, error) {
	pagination := Pagination{}

	limit, err := strconv.Atoi(c.DefaultQuery("take", "15"))
	if err != nil {
		return pagination, errors.New("query 'take' should be a number")
	}
	if limit <= 0 {
		return pagination, errors.New("query 'take' should be greater then 0")
	}

	offset, err := strconv.Atoi(c.DefaultQuery("skip", "0"))
	if err != nil {
		return pagination, errors.New("query 'skip' should be a number")
	}
	if offset <= 0 {
		return pagination, errors.New("query 'skip' must not be negative")
	}

	pagination.Limit = limit
	pagination.Offset = offset

	return pagination, nil
}

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
