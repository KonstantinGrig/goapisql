package jwtservice

import (
	"errors"
	"fmt"
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

//Parse Authorization Header and retrieves jwt.MapClaims
func Parse(authorizationHeader string) (jwt.MapClaims, error) {
	var res jwt.MapClaims
	if !strings.HasPrefix(authorizationHeader, config.PREFIX_TOKEN) {
		return res, errors.New("Error in authorization header: should be prefix 'Bearer '")
	}
	tokenString := authorizationHeader[len(config.PREFIX_TOKEN):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecret := config.GetString(config.KEY_JWT_SECRET)
		hmacSampleSecret := []byte(jwtSecret)
		return hmacSampleSecret, nil
	})

	if token == nil {
		return res, errors.New("Error in authorization token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		res = claims
	} else {
		return res, err
	}

	return res, nil
}

func CreateToken(data jwt.MapClaims) (token string, err error) {
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	jwtSecret := config.GetString(config.KEY_JWT_SECRET)
	return tokenStruct.SignedString([]byte(jwtSecret))
}
