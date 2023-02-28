package auth

import (
	"context"
	"time"

	"github.com/LSDXXX/libs/constant"
	"github.com/LSDXXX/libs/model"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/LSDXXX/libs/repo"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
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

func IdentityHandler(c *gin.Context, identityKey string) IdentityInfo {
	claims := jwt.ExtractClaims(c)
	info := claims[identityKey].(map[string]interface{})
	return IdentityInfo{
		Id:       cast.ToInt(info["id"]),
		UserName: cast.ToString(info["user_name"]),
		WSKey:    cast.ToString(info["ws_key"]),
	}
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
			log.WithContext(context.Background()).Debugf("payload : %+v", data)
			if v, ok := data.(model.User); ok {
				return jwt.MapClaims{
					identityKey: map[string]interface{}{
						"id":        v.Id,
						"ws_key":    uuid.NewString(),
						"user_name": v.UserName,
					},
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			return IdentityHandler(c, identityKey)
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				log.WithContext(c).Errorf("bind login val error: %s", err.Error())
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
			log.WithContext(context.Background()).Debugf("authorizator : %+v", data)
			if _, ok := data.(IdentityInfo); ok {
				return true
			}
			// if v, ok := data.(*User); ok && v.UserName == "admin" {
			// 	return true
			// }

			return false
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
