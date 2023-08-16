package ginutils

import (
	"errors"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type ApiScope uint

const (
	// scopes for "/admins/*"
	API_SCOPE_ADMIN ApiScope = iota
	// scopes for "/dashboard/*"
	API_SCOPE_DASHBOARD
	// scopes for open apis
	API_SCOPE_OPENAPI
	API_SCOPE_UNKNOWN
)

const (
	OpenJwtAppUserIdKey = "AppUserId"
	JwtIdentityKey      = "id"
)

func GetUserId(c *gin.Context) (uint, error) {
	apiScope := GetApiScope(c)
	switch apiScope {
	case API_SCOPE_OPENAPI:
		claims := jwt.ExtractClaims(c)
		userIdFloat, ok := claims[OpenJwtAppUserIdKey]
		if !ok {
			return 0, errors.New("not found user id")
		}
		userId := uint(userIdFloat.(float64))
		return userId, nil
	case API_SCOPE_DASHBOARD:
		return c.GetUint(JwtIdentityKey), nil
	case API_SCOPE_ADMIN:
		return c.GetUint(JwtIdentityKey), nil
	}
	return 0, errors.New("unkown api scope")
}

func MustGetUserId(c *gin.Context) uint {
	id, err := GetUserId(c)
	if err != nil {
		panic(err)
	}
	return id
}

func GetApiScope(c *gin.Context) ApiScope {
	if c.FullPath()[:len("/v1")] == "/v1" {
		return API_SCOPE_OPENAPI
	}
	if c.FullPath()[:len("/dashboard")] == "/dashboard" {
		return API_SCOPE_DASHBOARD
	}
	if c.FullPath()[:len("/admin")] == "/admin" {
		return API_SCOPE_ADMIN
	}
	return API_SCOPE_UNKNOWN
}
