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

		cabID := NearestCab(b.Lat, b.Long, b.CabType)

		if cabID != 0 {
			token := createTokenString(cabID, userid)

			c.JSON(http.StatusOK, gin.H{
				"Token": token,
			})

		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"Status": "Select Correct Cab Type",
			})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "Try Again",
		})
	}

}

func ConfirmRequest(c *gin.Context) {
	db := D.DB()

	var count = 0

	cid, uid := TokenValidator(c)

	stmt2, err := db.Prepare("insert into cost_detail (user_id,cab_id,last_updated) values(?,?,?)")

	checkErr(err)

	ress, err := stmt2.Exec(uid, cid, time.Now())
	fmt.Println(ress)
	if err != nil {
		log.Fatal(err)
	} else {
		count = 1
		stmt2, err := db.Prepare("update cab_location set ontrip=?, last_updated=? where id=?")

		checkErr(err)

		ress, err := stmt2.Exec(true, time.Now(), cid)
		fmt.Println(ress)
		if err != nil {
			log.Fatal(err)
		}

	}
	defer db.Close()

	if count == 1 {
		isUpdated := UserTimeUpdate(uid)
		if isUpdated {
			c.JSON(http.StatusOK, gin.H{
				"Status": "Now Your Ride has been Started --> Use End Trip to endride api",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status": "Now Your Ride has Not been Started --> Please Try Again",
			})
		}

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

func TimeCalulator(cTime string) int {

	t := time.Now()

	val := t.Format("20060102:150405")

	s1 := strings.Split(cTime, ":")
	s2 := strings.Split(val, ":")

	a, err := strconv.Atoi(s1[1])
	checkErr(err)
	b, err := strconv.Atoi(s2[1])
	checkErr(err)
	c := b - a
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	return c
}

func EndTrip(c *gin.Context) {

	db := D.DB()

	var (
		count    = 0
		distance float64
		cost     float64
		b        M.CostDetail
		bs       []M.CostDetail
	)
	c.BindJSON(&b)

	cid, uid := TokenValidator(c)
	cType, isRide := FindCab(cid)
	if isRide {
		clat, clng, cTime := FindUser(uid)

		distance = Distance(b.Lat, b.Long, clat, clng)

		totalDuration := TimeCalulator(cTime)

		if cType == "pink" {
			cost = (distance/1000)*2 + float64(totalDuration)/60 + 5
		} else {
			cost = (distance/1000)*2 + float64(totalDuration)/60
		}

		stmt2, err := db.Prepare("update cost_detail set distance=?, minute_travel=?, final_cost=?, last_updated=? where user_id=? and cab_id=?")

		checkErr(err)

		ress, err := stmt2.Exec(distance, totalDuration, cost, time.Now(), uid, cid)
		fmt.Println(ress)
		if err != nil {
			log.Fatal(err)
		} else {
			count = 1
			stmt2, err := db.Prepare("update cab_location set lat=?, lng=?, ontrip=?, last_updated=? where id=?")

			if err != nil {
				fmt.Println(err)
			}

			ress, err := stmt2.Exec(b.Lat, b.Long, false, time.Now(), cid)
			fmt.Println(ress)
			if err != nil {
				log.Fatal(err)
			} else {
				b.UserID = uid
				b.CabID = cid
				b.Distance = distance / 1000
				b.TripTime = totalDuration
				b.FinalCost = cost

			}

			bs = append(bs, b)
		}

		if count == 1 {
			c.JSON(http.StatusOK, gin.H{
				"Status": bs,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"Status": "Not Ok",
			})
		}

	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"Status": "You Haven't Started The Ride",
		})
	}

}

func FindUser(id int) (float64, float64, string) {
	db := D.DB()
	var (
		lat   float64
		lng   float64
		t     string
		count = 0
	)
	stmt, err := db.Prepare("select lat, lng, time_start from user_location where id=?")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&lat, &lng, &t)

		if err != nil {

			log.Fatal(err)

		}

		count = 1
	}

	if count == 1 {
		return lat, lng, t
	}
	return 0, 0, "time"

}
func FindCab(id int) (string, bool) {
	db := D.DB()
	var (
		cabtype string
		ontrip  bool
		count   = 0
	)
	stmt, err := db.Prepare("select cabtype, ontrip from cab_location where id=?")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&cabtype, &ontrip)

		if err != nil {

			log.Fatal(err)

		}

		count = 1
	}

	if count == 1 {
		return cabtype, ontrip
	}
	return "N", ontrip

}

func UserTimeUpdate(id int) bool {
	db := D.DB()

	var count = 0

	t := time.Now()

	val := t.Format("20060102:150405")

	stmt2, err := db.Prepare("update user_location set time_start=? where id=?")

	checkErr(err)

	ress, err := stmt2.Exec(val, id)
	fmt.Println(ress)
	if err != nil {
		log.Fatal(err)
	} else {
		count = 1
	}

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if count == 1 {
		return true
	}
	return false
}

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

func NearestCab(ulat, ulong float64, cab string) int {

	db := D.DB()

	var (
		count = 0
		clat  float64
		clong float64
		b     M.Distanc
		bs    []M.Distanc
	)

	stmt, err := db.Prepare("select id, lat, lng from cab_location where ontrip=? and cabtype=?")
	if err != nil {
		fmt.Println(err)
	}

	rows, err := stmt.Query(false, cab)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&b.CabID, &clat, &clong)

		if err != nil {

			log.Fatal(err)

		}

		b.Distance = Distance(ulat, ulong, clat, clong)

		bs = append(bs, b)

		count = 1
	}

	if count == 1 {
		id := NearestDistance(bs)

		return id

	}
	return 0

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

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {

	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
