package auth

import (
	"time"

	"github.com/LSDXXX/libs/constant"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/LSDXXX/libs/repo"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AuthHandler struct {
	jwtMiddleware *jwt.GinJWTMiddleware
	userMapper    repo.UserMapper `container:"type"`
}

func NewAuthHandler() (*AuthHandler, error) {
	out := AuthHandler{}
	util.PanicWhenError(container.Fill(&out))
	auth, err := newJWTAuth(constant.JWTIdentityKey, out.userMapper)
	if err != nil {
		return nil, err
	}
	out.jwtMiddleware = auth
	return &out, nil
}

func (a *AuthHandler) Middleware() gin.HandlerFunc {
	return a.jwtMiddleware.MiddlewareFunc()
}

func (a *AuthHandler) Use(e *gin.Engine) {
	e.POST("/login", a.jwtMiddleware.LoginHandler)
	e.POST("/logout", a.jwtMiddleware.LogoutHandler)
}

func newJWTAuth(identityKey string, mapper repo.UserMapper) (*jwt.GinJWTMiddleware, error) {
	middle, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "login",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		SendCookie:  true,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: IdentityInfo{
						Id:    v.Id,
						WSKey: uuid.NewString(),
					},
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return claims[identityKey]
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			info, err := mapper.GetByUserName(userID)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			if info.Password != password {
				return nil, jwt.ErrFailedAuthentication
			}

			return info, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*User); ok && v.UserName == "admin" {
			// 	return true
			// }

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		return nil, errors.Wrap(err, "create jwt middleware")
	}

	err = middle.MiddlewareInit()
	if err != nil {
		return nil, errors.Wrap(err, "init jwt middleware")
	}
	return middle, nil
}
