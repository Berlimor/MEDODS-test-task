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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// A shcema for the database instance
type User struct {
	GUID			string		`bson:"GUID"`
	firstname		string		`bson:"firstname"`
	lastname		string		`bson:"lastname"`
	access_token	string		`bson:"access_token"`
	refresh_token	string		`bson:"refresh_token"`
}


func main() {
	router := gin.Default()
	router.GET("/JWT/get-token/:guid", getTokenPair)
	router.GET("/JWT/refresh-token/:refresh", getRefresh)

	router.Run("localhost:8080")
}


func getTokenPair(c* gin.Context) {
	client, err := database.DBGetClient()
	if err != nil {
		panic(err)
	}
	defer database.DBDisconnect(client)
	
	guid := c.Param("guid")
	collection := client.Database("user-tokens").Collection("JWT")
	filter := bson.M{"GUID": guid}
	opt := options.FindOne().SetProjection(bson.D{{"access_token", 1}, {"refresh_token", 1}})
	cursor := collection.FindOne(context.TODO(), filter, opt)
	if cursor == nil {
		c.IndentedJSON(http.StatusNotFound, "This access token does not exist")
		return
	}
	c.IndentedJSON(http.StatusOK, cursor)
}


func getRefresh(c* gin.Context) {
	client, err := database.DBGetClient()
	if err != nil {
		panic(err)
	}
	defer database.DBDisconnect(client)

	old_refresh := c.Param("refresh")
	
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	secretKey := os.Getenv("SECRET_KEY")
	var hmacSampleSecret []byte = []byte(secretKey)
	
	collection := client.Database("user-tokens").Collection("JWT")
	filter := bson.M{"refresh_token": old_refresh}
	cursor := collection.FindOne(context.TODO(), filter)
	if cursor == nil {
		c.IndentedJSON(http.StatusNotFound, "User with this key not found")
		return
	}

	var user User
	if err = cursor.Decode(&user); err != nil {
		panic(err)
	}
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"GUID": user.GUID,
		"firstname": user.firstname,
		"lastname": user.lastname,
		"iat": time.Now(),
	})

	newAccessTokenString, err := newAccessToken.SignedString(hmacSampleSecret)
	if err != nil {
		panic(err)
	}

	// Data should be a random string
	data := "abcABC123!@#$%"
	newRefreshToken := base64.StdEncoding.EncodeToString([]byte(data))

	user.access_token = newAccessTokenString
	user.refresh_token = newRefreshToken
	result, err := collection.ReplaceOne(context.TODO(), filter, user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}