package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (c *Claims) Deadline() (deadline time.Time, ok bool) {
	panic("unimplemented")
}

func (c *Claims) Done() <-chan struct{} {
	panic("unimplemented")
}

func (c *Claims) Err() error {
	panic("unimplemented")
}

func (c *Claims) Value(key any) any {
	panic("unimplemented")
}

func CreateToken(userID string, role string, email string) (string, error) {
	claims := Claims{
		userID,
		role,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("jwt-secret-key"))
}

func KeyFunction(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.ErrTokenSignatureInvalid
	}
	return []byte("jwt-secret-key"), nil
}

func ParseJWT(token string, c *Claims) (*Claims, error) {
	parse, err := jwt.ParseWithClaims(token, c, KeyFunction)

	if err != nil {
		return nil, err
	}
	if !parse.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return c, err
}
