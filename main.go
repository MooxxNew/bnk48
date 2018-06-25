package main

import (	
	"github.com/pallat/tis620"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"flag"
	"time"
	"github.com/dgrijalva/jwt-go"
	"crypto/sha256"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/globalsign/mgo"
)

const (
	signature = "drowssap"
)

var (
	secret *string
)

type accessToken struct {
	Token string `json:"accessToken"`
	ExpiresIn int64 `json:"expiresIn"`
}

type signupName struct{
	Username string
	Password string
	Email string
}

// type users struct {
// 	Users []userDetail
// }

type userDetail struct {
	UserID int `json:"userId"`
	ID int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func signup(c echo.Context) error {
	var m signupName
	err := c.Bind(&m)
	if err != nil {
		c.Error(err)
	}

	if m.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message" : "email is require",
		})
	}

	sum := sha256.Sum256([]byte(m.Password))
	m.Password = fmt.Sprintf("%x", sum)

	fmt.Println(m)

	err = mgoInsert(m)
	if err != nil {
		c.Error(err)
	}

	token , err := genareteToken()
	if err != nil {
		c.Error(err)
	}

	return c.JSON(http.StatusOK, accessToken{
		Token: token,
		ExpiresIn: int64(time.Hour.Minutes()),
	})
}

func mgoInsert(data signupName) error {
	url := "mongodb://localhost:27017"
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	col := session.DB("odds").C("credential")
	_, err = col.Upsert(bson.M{"email": data.Email}, data)
	return err
}

func genareteToken() (string, error) {
	mySigningKey := []byte(*secret)

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Issuer: "odds",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func init() {
	secret = flag.String("secret", signature, "-secret=yourpassword")
}

func getData(c echo.Context) error{
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(err)
	}
	var user,result []userDetail
	resp, err := http.Get("http://jsonplaceholder.typicode.com/posts")
	if err != nil {
		c.Error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &user)

	if id == 0 { return c.JSON(http.StatusOK, user)}
	for _, item := range user {
		if item.ID == id {
			result = append(result, item)
		}
	}
	

	return c.JSON(http.StatusOK, result)
}

func getThai(c echo.Context) error{
	url := "http://192.168.100.3:8080/thai"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "abe128b1-afb8-43d8-a034-72cc6fc10b33")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(tis620.ToUTF8(string(body)))

	return c.JSON(http.StatusOK, string(body))
}

func main() {
	flag.Parse()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", hello)
	e.POST("/signup", signup)
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}, middleware.JWT([]byte(*secret)))


	e.GET("/post/:id", getData)
	e.GET("/post", getData)
	e.GET("/thai", getThai)

	e.Logger.Fatal(e.Start(":1323"))
}

