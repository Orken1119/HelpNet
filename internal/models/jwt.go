package models

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	RoleID uint `json:"role"`
	ID     uint `json:"id"`
	jwt.RegisteredClaims
}

type JwtUser interface {
	GetID() uint
	GetRoleID() uint
}

type JwtRefreshClaims struct {
	RoleID uint `json:"role"`
	ID     uint `json:"id"`
	jwt.RegisteredClaims
}
