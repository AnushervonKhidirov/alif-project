package middleware

import (
	"errors"
	"net/http"
	"pet/internal/model"
	"pet/internal/service"
	"pet/pkg/constant"
	"pet/pkg/response"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const Authorization = "Authorization"

func Auth(tokenService service.TokenService, permittedRoles []model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromHeader(c)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, err)
			return
		}

		tokenClaims, err := tokenService.VerifyAccess(token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, err)
			return
		}

		permittedRequest := isPermitted(permittedRoles, tokenClaims.UserRole)
		if !permittedRequest {
			response.Fail(c, http.StatusForbidden, nil)
			return
		}

		c.Set(constant.UserId, tokenClaims.UserId)

		c.Next()
	}
}

func getTokenFromHeader(c *gin.Context) (string, error) {
	header := c.Request.Header

	if authorization, ok := header[Authorization]; ok {
		token := strings.ReplaceAll(authorization[0], "Bearer ", "")
		return token, nil
	}

	return "", errors.New("access token not found")
}

func isPermitted(permittedRoles []model.Role, role model.Role) bool {
	return permittedRoles == nil || (len(permittedRoles) == 0 || slices.Contains(permittedRoles, role))
}
