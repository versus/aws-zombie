package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"log"
)

type message struct {
	phone_from string
	phone_to   string
	text       string
}

//func main() {
//	ctx := runtime.Context{}
//	evt := &apigatewayproxyevt.Event{
//		HTTPMethod: "POST",
//		Body: `{"phone_from" : "380636637354",
//			   "phone_to": "380636637355"
//		}`,
//	}
//
//	res, err := Handle(evt, &ctx)
//	log.Println(res, err)
//}

func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (interface{}, error) {
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(evt.Body), &data); err != nil {
		return "", err
	}

	log.Printf("evt: %+v \n", evt)
	log.Printf("query: %+v \n", evt.QueryStringParameters)

	phoneS, ok := data["phone_from"].(string)
	if !ok {
		return "Invalid phone", errors.New("Invalid phone")
	}

	phoneTo, ok := data["phone_to"].(string)
	if !ok {
		return "Invalid phone_to", errors.New("Invalid phone_to")
	}

	sess := getSession()
	db := dynamodb.New(sess)

	params := dynamodb.QueryInput{
		TableName: aws.String(`message`),
		//Limit:     aws.Int64(20),
		KeyConditionExpression: aws.String("#from = :p1 "),
		FilterExpression:       aws.String("#to = :p2 "),
		//IndexName:              aws.String("timestamp"),
		//ScanIndexForward: aws.Bool(true),

		ExpressionAttributeNames: map[string]*string{
			"#from": aws.String("from"),
			"#to":   aws.String("to"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":p1": {S: aws.String(phoneS)},
			":p2": {S: aws.String(phoneTo)},
		},
	}

	//if timestamp, ok := data["timestamp"].(string); ok {
	//	phoneFrom, ok :=  data["phone_from"].(string); ok
	//	params.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
	//		`from`:      {S: aws.String(data["phone_from"]).()},
	//		`timestamp`: {S: aws.String(timestamp).(string)},
	//	}
	//}

	res, err := db.Query(&params)
	if err != nil {
		return "", err
	}

	var (
		list []map[string]string
	)
	var line map[string]string
	for _, v := range res.Items {
		line = map[string]string{
			"from":      *v["from"].S,
			"text":      *v["text"].S,
			"to":        *v["to"].S,
			"timestamp": *v["timestamp"].S,
		}
		list = append(list, line)

	}

	//b, err := json.Marshal(list)
	//bS := string(b)
	return list, err
}

func mapMessages(output *dynamodb.QueryOutput) map[string]interface{} {
	//res := struct{"last_key" : output.LastEvaluatedKey["timestamp"].S}
	//return res
	return nil
}

func getSession() *session.Session {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
		//	Credentials: cred,
	}))

	return sess
}
