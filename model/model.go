package model

import jwt "github.com/dgrijalva/jwt-go"

type CabLocation struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	CabType string  `json:"cabtype"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}

type CostDetail struct {
	ID        int     `json:"id"`
	UserID    int     `json:"uid"`
	CabID     int     `json:"cid"`
	Distance  float64 `json:"distance"`
	TripTime  float64 `json:"triptime"`
	FinalCost float64 `json:"tripcost"`
}

type LoginClaims struct {
	ID       int    `json:"id"`
	Password string `json:"password"`

	jwt.StandardClaims
}

type UserClaims struct {
	CabID  int `json:"cid"`
	UserID int `json:"uid"`

	jwt.StandardClaims
}

type Distanc struct {
	CabID    int
	Distance float64
}
