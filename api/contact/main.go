package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"log"
	"strings"
)

//func main() {
//	ctx := runtime.Context{}
//	evt := &apigatewayproxyevt.Event{
//		HTTPMethod: "POST",
//		Body:       `{"contact_phone" : "05"}`,
//		QueryStringParameters: map[string]string{"phone": "8"},
//	}
//
//	res, err := Handle(evt, &ctx)
//	log.Println(res, err)
//}

func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (string, error) {
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(evt.Body), &data); err != nil {
		log.Printf("Error: %+v \n", err)
		return err.Error(), err
	}

	log.Printf("evt: %+v \n", evt)
	log.Printf("query: %+v \n", evt.QueryStringParameters)
	phone := strings.TrimSpace(evt.QueryStringParameters["phone"])

	if phone == "" {
		return "Empty Phone", errors.New("Invalid Phone")
	}

	if sPhone, ok := data["contact_phone"].(string); sPhone == "" || !ok {
		return "Wrong Contact Phone", errors.New("Wrong Contact Phone")
	}

	switch evt.HTTPMethod {
	case "POST":
		return addContact(data, phone)
	}
	return "Err", nil
}

func getSession() *session.Session {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
		//	Credentials: cred,
	}))

	return sess
}

func addContact(data map[string]interface{}, phone string) (string, error) {
	sess := getSession()
	db := dynamodb.New(sess)

	//Get User
	params := dynamodb.GetItemInput{
		TableName: aws.String(`user`),
		Key: map[string]*dynamodb.AttributeValue{
			`phone`: {S: aws.String(phone)},
		},
	}
	itemOut, err := db.GetItem(&params)
	if err != nil {
		return "No User", err
	}
	sPhone, _ := data["contact_phone"].(string)
	if _, ok := itemOut.Item["contacts"]; !ok {
		var s []*string
		itemOut.Item["contacts"] = &dynamodb.AttributeValue{SS: s}
	}

	if contains(itemOut.Item["contacts"].SS, sPhone) {
		return "duplicate", nil
	}

	//Adding contact
	if contactList, ok := itemOut.Item["contacts"]; ok != false {
		itemOut.Item["contacts"].SS = append(contactList.SS, &sPhone)
	}

	//Update record
	upQ := dynamodb.UpdateItemInput{
		TableName: aws.String(`user`),
		Key: map[string]*dynamodb.AttributeValue{
			`phone`: {S: aws.String(phone)},
		},
		UpdateExpression: aws.String(`set contacts = :p1`),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			`:p1`: itemOut.Item["contacts"],
		},
	}

	if _, err := db.UpdateItem(&upQ); err != nil {
		return err.Error(), err
	}

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
