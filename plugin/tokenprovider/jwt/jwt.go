package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hoangtk0100/social-todo-list/common"
	"github.com/hoangtk0100/social-todo-list/component/tokenprovider"
)

type jwtProvider struct {
	prefix string
	secret string
}

func NewJWTProvider(prefix string, secret string) *jwtProvider {
	return &jwtProvider{prefix: prefix, secret: secret}
}

type myClaims struct {
	Payload common.TokenPayload `json:"payload"`
	jwt.StandardClaims
}

type token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

func (t *token) GetToken() string {
	return t.Token
}

func (j *jwtProvider) SecretKey() string {
	return j.secret
}

func (j *jwtProvider) Generate(data tokenprovider.TokenPayload, expiry int) (tokenprovider.Token, error) {
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		common.TokenPayload{
			UId:   data.UserId(),
			URole: data.Role(),
		},
		jwt.StandardClaims{
			ExpiresAt: now.Local().Add(time.Second * time.Duration(expiry)).Unix(),
			IssuedAt:  now.Local().Unix(),
			Id:        fmt.Sprintf("%d", now.UnixNano()),
		},
	})

	myToken, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &token{
		Token:   myToken,
		Expiry:  expiry,
		Created: now,
	}, nil
}

func (j *jwtProvider) Validate(mytoken string) (tokenprovider.TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, tokenprovider.ErrInvalidToken
		}
		return []byte(j.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(mytoken, &myClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, tokenprovider.ErrExpiredToken) {
			return nil, tokenprovider.ErrExpiredToken
		}

		return nil, tokenprovider.ErrInvalidToken
	}

	if !jwtToken.Valid {
		return nil, tokenprovider.ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*myClaims)
	if !ok {
		return nil, tokenprovider.ErrInvalidToken
	}

	return claims.Payload, nil
}
