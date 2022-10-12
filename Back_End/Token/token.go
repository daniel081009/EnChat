package Token

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthTokenClaims struct {
	TokenUUID string `json:"tid"`
	UserID    int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"mail"`

	jwt.StandardClaims
}

func CreateJWT(UserId int, Name, Email string) (string, error) {
	at := AuthTokenClaims{
		TokenUUID: uuid.NewString(),
		UserID:    UserId,
		Name:      Name,
		Email:     Email,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "EnChat",
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)),
		},
	}

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, err := atoken.SignedString([]byte(os.Getenv("SECRET_KEY")))

	return signedAuthToken, err
}

func CheckJWT(Token string) (AuthTokenClaims, error) {
	claims := AuthTokenClaims{}
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected Signing Method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	}
	_, err := jwt.ParseWithClaims(Token, &claims, key)

	return claims, err
}
func AuthMiddleWare(c *gin.Context) {
	token, err := c.Request.Cookie("Token")
	if err != nil {
		c.JSON(200,
			gin.H{"code": 0, "error": "Authentication failed"})
		c.Abort()
		return
	}
	_, err = CheckJWT(token.Value)
	if err != nil {
		c.JSON(200,
			gin.H{"code": 0, "error": "Token Not available"})
		c.Abort()
	}
	c.Next()
}
