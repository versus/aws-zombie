package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/satori/go.uuid"
)

//func main() {
//	ctx := runtime.Context{}
//	evt := &apigatewayproxyevt.Event{
//		HTTPMethod: "POST",
//		Body: `{"text" : "dadfsdf",
//			   "from_phone": "380934235345",
//			   "to" : "354767789"
//		}`,
//	}
//
//	res, err := Handle(evt, &ctx)
//	log.Println(res, err)
//}

const (
	QueueUrl = "https://sqs.eu-west-1.amazonaws.com/ssages"
	Region   = "eu-west-1"
)

type Sender struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type SendRequest struct {
	Phone     string `json:"phone"` //to
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	From      Sender `json:"sender"`
}

func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(evt.Body), &data); err != nil {
		return ``, err
	}

	if err := validateInput(data); err != nil {
		return "", err
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(Region),
		//	Credentials: cred,
	}))

	svdb := dynamodb.New(sess)
	svsqs := sqs.New(sess)

	if _, err := saveToDynamo(data, svdb, timestamp); err != nil {
		return "", err
	}

	user, err := getUser(svdb, data[`from_phone`].(string))
	if err != nil {
		return "", err
	}

	sender, err := mapSenderData(user)
	if err != nil {
		return "", err
	}

	sendrequest := &SendRequest{
		Phone:     data[`to`].(string),
		Timestamp: timestamp,
		Message:   data[`text`].(string),
		From:      sender,
	}
	if _, err := sendToSQS(svsqs, sendrequest); err != nil {
		return "", err
	}
	return "OK", nil
}

func sendToSQS(svc *sqs.SQS, sendrequest *SendRequest) (string, error) {

	s, err := json.Marshal(sendrequest)
	if err != nil {
		return "", err
	}

	// Send message
	send_params := &sqs.SendMessageInput{
		MessageBody:  aws.String(string(s)), // Required
		QueueUrl:     aws.String(QueueUrl),  // Required
		DelaySeconds: aws.Int64(3),
	}
	send_resp, err := svc.SendMessage(send_params)
	if err != nil {
		return "", err
	}
	fmt.Printf("[Send message] \n%v \n\n", send_resp)
	return "Send", nil
}

func saveToDynamo(data map[string]interface{}, db *dynamodb.DynamoDB, timestamp string) (string, error) {

	uID := uuid.NewV4()

	item := dynamodb.PutItemInput{
		TableName: aws.String(`message`),
		Item: map[string]*dynamodb.AttributeValue{
			`id`:        {S: aws.String(uID.String())},
			`from`:      {S: aws.String(data[`from_phone`].(string))},
			`to`:        {S: aws.String(data[`to`].(string))},
			`text`:      {S: aws.String(data[`text`].(string))},
			`timestamp`: {S: aws.String(timestamp)},
		},
	}
	if _, err := db.PutItem(&item); err != nil {
		return err.Error(), err
	}
	return "Ok", nil
}

func validateInput(data map[string]interface{}) error {
	if _, ok := data[`from_phone`]; !ok {
		return errors.New(`Empty phone`)
	}

	if _, ok := data[`to`]; !ok {
		return errors.New(`Empty recepient`)
	}

	if _, ok := data[`text`]; !ok {
		return errors.New(`Empty text`)
	}

	return nil
}

func mapSenderData(data *dynamodb.GetItemOutput) (u Sender, err error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		u = Sender{}
	//		log.Println(r)
	//		err = errors.New("Error retrieving data")
	//	}
	//}()

	if _, ok := data.Item[`phone`]; ok {
		u.Phone = *data.Item[`phone`].S
	}

	if _, ok := data.Item[`first_name`]; ok {
		u.FirstName = *data.Item[`first_name`].S
	}

	if _, ok := data.Item[`last_name`]; ok {
		u.LastName = *data.Item[`last_name`].S
	}

	return u, nil
}

func getUser(db *dynamodb.DynamoDB, phone string) (*dynamodb.GetItemOutput, error) {
	params := dynamodb.GetItemInput{
		TableName: aws.String(`user`),
		Key: map[string]*dynamodb.AttributeValue{
			`phone`: {S: aws.String(phone)},
		},
	}

	res, err := db.GetItem(&params)
	if err != nil {
		return nil, err
	}

	return res, nil
}
