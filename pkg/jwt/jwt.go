package jwt

import (
	"strconv"
	"time"

	"github.com/VariableSan/gia-sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(
	user models.User,
	jwtSecret string,
	duration time.Duration,
) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = strconv.Itoa(int(user.ID))
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
