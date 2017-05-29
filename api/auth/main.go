package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dgrijalva/jwt-go"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
)

func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (string, error) {
	rq := make(map[string]interface{})
	if err := json.Unmarshal([]byte(evt.Body), &rq); err != nil {
		return "", err
	}

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))

	db := dynamodb.New(sess)

	params := dynamodb.GetItemInput{
		TableName: aws.String(`user`),
		Key: map[string]*dynamodb.AttributeValue{
			`phone`: {S: aws.String(fmt.Sprint(rq[`phone`]))},
		},
	}

	res, err := db.GetItem(&params)
	if err != nil {
		return "", err
	}

	if _, ok := res.Item[`phone`]; !ok {
		return "", errors.New("Error: User does not exists")
	}

	hash := sha256.New()
	hash.Write([]byte(fmt.Sprint(rq[`password`])))

	if _, ok := res.Item[`password`]; !ok {
		return "", errors.New("Error: User does not exists")
	}

	if fmt.Sprintf(`%x`, hash.Sum(nil)) != *res.Item[`password`].S {
		return "", errors.New("Error: Wrong password")
	}

	token := getJWT(fmt.Sprint(*res.Item[`phone`].S))

	if err := setUserToken(fmt.Sprint(*res.Item[`phone`].S), token, db); err != nil {
		return "", errors.New(fmt.Sprint("Error: ", err))
	}

	return token, nil
}

func setUserToken(phone, token string, db *dynamodb.DynamoDB) error {
	q := dynamodb.UpdateItemInput{
		TableName: aws.String(`user`),
		Key: map[string]*dynamodb.AttributeValue{
			`phone`: {S: aws.String(phone)},
		},
		UpdateExpression: aws.String(`set jwt = :p1`),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			`:p1`: {S: aws.String(token)},
		},
	}

	_, err := db.UpdateItem(&q)

	return err
}

func getJWT(phone string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": phone,
		"gdate": time.Now().Unix(),
		"ttl":   3600,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte(`covabunga`))

	return tokenString
}
