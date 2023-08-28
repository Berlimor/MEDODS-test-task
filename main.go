package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"MEDODS-test-task/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// A shcema for the database instance
type User struct {
	ID				primitive.ObjectID		`bson:"_id"`
	GUID			string					`bson:"GUID"`
	Firstname		string					`bson:"Firstname"`
	Lastname		string					`bson:"Lastname"`
	Access_token	string					`bson:"Access_token"`
	Refresh_token	string					`bson:"Refresh_token"`
}

// Struct for returning a token pair 
type TokenPair struct {
	Access_token	string
	Refresh_token	string
}


func main() {
	//	Add a new test user to the database
	client, err := database.DBGetClient()
	if err != nil {
		panic(err)
	}
	
	firstUser := User{
		ID: primitive.NewObjectID(),
		GUID: "a6bcd248-c9ce-4475-9caf-b3313af3f14c",
		Firstname: "El",
		Lastname: "Yusufov",
		Access_token: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJHVUlEIjoiYTZiY2QyNDgtYzljZS00NDc1LTljYWYtYjMzMTNhZjNmMTRjIiwic3ViIjoiYWNjZXNzX3Rva2VuIiwiZXhwIjoxNjkzMTQ0NDExLCJpYXQiOjE2OTMxNDQ0MTF9.tE_izf1PJcOsqgi5DRI3hL0_V4Sf1HfES68WEJKsRU8MOjW3X52n-HA4RsJMYgxg9sgigmvpJtd-69YY5QyXzg",
		Refresh_token: "YWJjZGVBQkNERTEyMyFAIyQl",
	}
	database.AddNewDocumentForTest(client, firstUser)
	database.DBDisconnect(client)
	
	router := gin.Default()
	router.GET("/auth/token/:guid", getTokenPair)
	router.GET("/auth/refresh/:refresh", getRefresh)

	router.Run("localhost:8080")
}


func getTokenPair(c* gin.Context) {
	client, err := database.DBGetClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Failed rto connect to db")
	}
	defer database.DBDisconnect(client)
	
	guid := c.Param("guid")
	collection := client.Database("user-tokens").Collection("JWT")
	filter := bson.D{{"GUID", guid}}
	opt := options.FindOne().SetProjection(bson.D{{"Access_token", 1}, {"Refresh_token", 1}})
	cursor := collection.FindOne(context.TODO(), filter, opt)
	if cursor == nil {
		c.IndentedJSON(http.StatusNotFound, "This access token does not exist")
		return
	}

	
	var tokens TokenPair
	if err = cursor.Decode(&tokens); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, tokens)
}


func getRefresh(c* gin.Context) {
	client, err := database.DBGetClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	defer database.DBDisconnect(client)

	old_refresh := c.Param("refresh")
	
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	secretKey := os.Getenv("SECRET_KEY")
	var hmacSampleSecret []byte = []byte(secretKey)
	
	collection := client.Database("user-tokens").Collection("JWT")
	filter := bson.D{{"Refresh_token", old_refresh}}
	cursor := collection.FindOne(context.TODO(), filter)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, "User with this key not found")
		return
	}

	var user User
	if err = cursor.Decode(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Could not decode")
	}
	
	type CustomClaims struct {
		GUID			string		`bson:"GUID"`
		firstname		string		`bson:"firstname"`
		lastname		string		`bson:"lastname"`
		jwt.RegisteredClaims
	}
	claims := CustomClaims{
		user.GUID,
		user.Firstname,
		user.Lastname,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10000)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject: "access_token",
		},
	}
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Converting the token to a string
	newAccessTokenString, err := newAccessToken.SignedString(hmacSampleSecret)
	if err != nil {
		panic(err)
	}

	// Data variable should be a random string for security purposes
	data := "abcdeABCDE123!@#$%"
	newRefreshToken := base64.StdEncoding.EncodeToString([]byte(data))

	// Creating the return value for IndentedJSON
	newTokenPair := TokenPair{newAccessTokenString, newRefreshToken}

	update := bson.D{{
		"$set", 
		bson.D{{"Access_token", newAccessTokenString}, 
			   {"Refresh_token", newRefreshToken}},
	}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, newTokenPair)
}