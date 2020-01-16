package proxy

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var hmacSampleSecret = []byte("my_secret_key")

func verifyUserName(user string) bool {
	return user != ""
}

func verifyUserNameAndPassword(user, pass string) bool {
	return user != "" && pass != ""
}

func signToken(user string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": user})
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		log.Errorf("signToken Error: %v", err)
		return ""
	}
	return tokenString
}

func verifyToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil || token == nil {
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Tracef("verifyToken: %v", claims)
		return verifyUserName(claims["user"].(string))
	}
	return false
}
