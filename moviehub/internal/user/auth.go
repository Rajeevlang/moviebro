package user

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Customclaim struct {
	jwt.RegisteredClaims
	ID       string `json:"id"`
	Username string `json:"username"`
}

func Newjwttoken(id string, username string, seceretkey []byte) (string, error) {

	srvclaim := Customclaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "server",
		},
		ID:       id,
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, srvclaim)
	jwttoken, err := token.SignedString(seceretkey)
	if err != nil {
		return "", err
	}

	return jwttoken, nil

}

func Refreshtoken(id string, username string, seceretkey []byte) (string, error) {

	srvclaim := Customclaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "server",
		},
		ID:       id,
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, srvclaim)
	jwttoken, err := token.SignedString(seceretkey)
	if err != nil {
		return "", err
	}

	return jwttoken, nil
}

func Validaccesstoken(token string, seceretkey []byte) (*Customclaim, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Customclaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return seceretkey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*Customclaim); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func Validrefreshtoken(token string, seceretkey []byte) (*Customclaim, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Customclaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return seceretkey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*Customclaim); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
