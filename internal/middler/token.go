package middler

import (
	"workflow/config"
	"workflow/internal/common"
	"workflow/internal/constants"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 获取AT中的用户信息
func SetUserIdToCtxHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		tk_str, _ := c.Cookie(constants.ACCESSTOKENCOOKIEKEY)
		if tk_str != "" {
			if token, err := tokenParse(tk_str); err == nil {
				claims := token.Claims.(*common.CustomClaims)
				c.Set(constants.ACCESSTOKENUSERIDKEY, claims.Id)
			}
		}
		c.Next()
	}
}

func tokenParse(tk_str string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tk_str, &common.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return jwt.ErrInvalidKey, nil
		}
		return []byte(config.Conf.Jwt.Secret), nil
	})
	return token, err
}
