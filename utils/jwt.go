package utils

import (
	"crontab/global"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	ID          uint   `json:"id"`
	Nickname    string `json:"nickname"`
	AuthorityId uint   `json:"authority_id"`
	jwt.StandardClaims
}

type JWT struct {
	SigningKey []byte
}

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token ")
)

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.Settings.JwtInfo.Key),
	}
}

//创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(j.SigningKey)
}

//解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, ErrTokenInvalid

	} else {
		return nil, ErrTokenInvalid

	}
}

//更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", ErrTokenInvalid
}

//登陆获取token
func LoginCreateToken(Id int, Nickname string, Role int) (string, error) {
	//生成token信息
	j := NewJWT()
	claims := CustomClaims{
		ID:          uint(Id),
		Nickname:    Nickname,
		AuthorityId: uint(Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			// TODO 设置token过期时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //token -->30天过期
			Issuer:    "Ezreal-Rao",
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return "", errors.New("token 生成失败")
	}
	return token, nil
}
