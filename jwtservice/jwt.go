package jwtservice

import (
	"errors"
	"fmt"
	"github.com/KonstantinGrig/goapisql/config"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

const (
	prefixToken  = "Bearer "
	keyJwtSecret = "jwt-secret"
)

func Parse(authorizationHeader string) (jwt.MapClaims, error) {
	var res jwt.MapClaims
	if !strings.HasPrefix(authorizationHeader, prefixToken) {
		return res, errors.New("Error in authorization header: should be prefix 'Bearer '")
	}
	tokenString := authorizationHeader[len(prefixToken):]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		jwtSecret := config.GetString(keyJwtSecret)
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
