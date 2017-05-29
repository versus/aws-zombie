package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"log"
	"strings"
)

func main() {
	ctx := runtime.Context{}
	evt := &apigatewayproxyevt.Event{
		HTTPMethod:            "GET",
		QueryStringParameters: map[string]string{"keyword": "1"},
	}

	res, err := Handle(evt, &ctx)
	log.Println(res, err)
}

func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (string, error) {
	keyword := strings.TrimSpace(evt.QueryStringParameters["keyword"])

	log.Printf("query: %+v \n", evt.QueryStringParameters)

	if keyword == "" {
		return "Empty KeyWord", errors.New("Invalid Keyword")
	}

	switch evt.HTTPMethod {
	case "GET":
		return findUser(keyword)
	}
	return "Err", nil
}

func getSession() *session.Session {
	cred := credentials.NewStaticCredentials(
		"xxxxxxx",
		"xxxxxxxxxxx",
		"")
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: cred,
	}))

	return sess
}

func findUser(keyword string) (string, error) {
	sess := getSession()
	db := dynamodb.New(sess)

	//Get User
	params := dynamodb.QueryInput{
		TableName:                 aws.String(`user`),
		KeyConditionExpression:    aws.String(`contains(first_name, :p1)`),
		FilterExpression:          aws.String(),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{`:p1`: &dynamodb.AttributeValue{S: &keyword}},
	}

	qO, err := db.Query(&params)
	log.Println(err)
	for _, Item := range qO.Items {
		log.Println(Item)
	}

	//Update record
	//upQ := dynamodb.UpdateItemInput{
	//	TableName: aws.String(`user`),
	//	Key: map[string]*dynamodb.AttributeValue{
	//		`phone`: {S: aws.String(phone)},
	//	},
	//	UpdateExpression: aws.String(`set contacts = :p1`),
	//	ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
	//		`:p1`: itemOut.Item["contacts"],
	//	},
	//}
	//

	return "Ok", nil
}

func contains(s []*string, e string) bool {
	for _, a := range s {
		if *a == e {
			return true
		}
	}
	return false
}
