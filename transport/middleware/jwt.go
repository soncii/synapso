package middles

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"synapso/model"
	"time"
)

var jwtSecret = []byte("your-secret-key")

func GenerateToken(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthAndExtractUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenString := ctx.Request().Header.Get("Authorization")
		if tokenString == "" {
			return ctx.String(http.StatusUnauthorized, "Unauthorized")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return ctx.String(http.StatusUnauthorized, "Unauthorized")
		}

		claims := token.Claims.(jwt.MapClaims)
		user := model.User{
			ID:   uint(claims["user_id"].(float64)),
			Role: claims["role"].(string),
		}

		ctx.Set("user", user)

		return next(ctx)
	}
}
