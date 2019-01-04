package controller

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	D "../dbPool"
	M "../model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var err error

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	privKeyPath = "key/app.rsa"
	pubKeyPath  = "key/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		fmt.Println(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		fmt.Println(err)
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		fmt.Println(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		fmt.Println(err)
	}

}

func createTokenString(cid, uid int) string {
	now := time.Now()
	iat := now.Unix()
	s := strconv.FormatInt(iat, 10)
	jti := s + "fuber"

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), jwt.MapClaims{
		"cid": cid,
		"uid": uid,
		"iat": iat,
		"jti": jti,
	})

	tokenstring, err := token.SignedString(signKey)
	if err != nil {
		log.Fatalln(err)
	}
	return tokenstring
}

func RequestRide(c *gin.Context) {

	var (
		b M.CabLocation
	)

	c.BindJSON(&b)
	IsInserted, userid := UserLocation(b.Lat, b.Long, b.Name)
	if IsInserted {

		cabID := NearestCab(b.Lat, b.Long)

		token := createTokenString(cabID, userid)

		c.JSON(http.StatusOK, gin.H{
			"Token": token,
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "Try Again",
		})
	}

}

func ConfirmRequest(c *gin.Context) {
	db := D.DB()

	var count = 0

	uid, cid := TokenValidator(c)

	stmt2, err := db.Prepare("insert into cost_detail (user_id,cab_id,last_updated) values(?,?,?)")

	checkErr(err)

	ress, err := stmt2.Exec(uid, cid, time.Now())
	fmt.Println(ress)
	if err != nil {
		log.Fatal(err)
	} else {
		count = 1
	}
	defer db.Close()

	if count == 1 {
		c.JSON(http.StatusOK, gin.H{
			"Status": "Now Your Ride has been Started --> Use End Trip to end the ride",
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"Status": "try again",
		})
	}

}

func TokenValidator(c *gin.Context) (int, int) {

	var (
		user M.UserClaims
	)
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	reqToken = strings.TrimSpace(splitToken[1])

	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		fmt.Println(err)
	}

	token, err = jwt.ParseWithClaims(reqToken, &user, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		fmt.Println(err)
	}

	if token.Valid == true {
		return user.CabID, user.UserID
	}
	return 0, 0
}

// func EndTrip(c *gin.Context) {

// 	//

// }

func UserLocation(lat, long float64, name string) (bool, int) {
	db := D.DB()

	var count = 0

	stmt2, err := db.Prepare("insert into user_location (name,lat,lng,last_updated) values(?,?,?,?)")

	checkErr(err)

	ress, err := stmt2.Exec(name, lat, long, time.Now())
	fmt.Println(ress)
	if err != nil {
		log.Fatal(err)
	} else {
		count = 1
	}
	id, err := ress.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if count == 1 {
		return true, int(id)
	}
	return false, 0
}

func NearestCab(ulat, ulong float64) int {

	db := D.DB()

	var (
		// counter = 0
		clat  float64
		clong float64
		b     M.Distanc
		bs    []M.Distanc
	)

	stmt, err := db.Prepare("select id, lat, lng from cab_location where ontrip=?")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query(false)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&b.CabID, &clat, &clong)

		if err != nil {

			log.Fatal(err)

		}

		b.Distance = math.Sqrt((ulat-clat)*(ulat-clat) + (ulong-clong)*(ulong-clong))
		bs = append(bs, b)
	}
	id := NearestDistance(bs)

	return id

}

func NearestDistance(d []M.Distanc) int {

	min := d[0].Distance
	var id int
	for i := 0; i < len(d); i++ {
		if d[i].Distance < min {
			min = d[i].Distance
			id = d[i].CabID
		}
	}

	fmt.Println(min)

	return id
}
