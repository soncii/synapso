package middles

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"net/http"
	"synapso/enums"
	"synapso/model"
	"time"
)

var jwtSecret = []byte("Damir Gimaletdinov")

func GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 365).Unix(),
	})
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
		if len(tokenString) < 7 {
			return ctx.String(http.StatusUnauthorized, "Unauthorized")
		}
		tokenString = tokenString[7:]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return ctx.String(http.StatusUnauthorized, "Unauthorized")
		}

		claims := token.Claims.(jwt.MapClaims)
		user := model.User{
			ID:   int(claims["user_id"].(float64)),
			Role: claims["role"].(string),
		}

		ctx.Set("user", user)

		return next(ctx)
	}
}

func AuthResearcherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if GetUserRoleFromContext(ctx) != enums.RESEARCHER {
			return ctx.String(http.StatusUnauthorized, "Unauthorized: you are not a researcher")
		}
		return next(ctx)
	}
}

func AuthSubjectMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if GetUserRoleFromContext(ctx) != enums.SUBJECT {
			return ctx.String(http.StatusUnauthorized, "Unauthorized: you are not a subject")
		}
		return next(ctx)
	}
}

func getUserModelFromContext(ctx echo.Context) model.User {
	return ctx.Get("user").(model.User)
}

func GetUserIDFromContext(ctx echo.Context) int {
	return getUserModelFromContext(ctx).ID
}

func GetUserRoleFromContext(ctx echo.Context) string {
	return getUserModelFromContext(ctx).Role
}
