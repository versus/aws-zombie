package sqs

import (
	"fmt"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func New() *sqs.SQS {
	creds := credentials.NewStaticCredentials(
		"public-key",
		"private-key",
		"",
	)

	config := &aws.Config{
		Region: aws.String("eu-west-1"),
		//Endpoint:    aws.String("s3.amazonaws.com"),
		Credentials: creds,
		MaxRetries:  aws.Int(5),
	}

	//sess := session.Must(session.NewSession(s3config))
	sess := session.New(config)

	// Create the service's client with the session.
	return sqs.New(sess)
}

func Listen(conn *sqs.SQS, callback func(rq *SendRequest)) {
	url := "https://sqs.eu-west-1.amazonaws.com/messages"

	params := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(url), // Required
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout: aws.Int64(1),
		WaitTimeSeconds:   aws.Int64(1),
	}

	go func() {
		for {
			//fmt.Println("READ!")
			resp, err := conn.ReceiveMessage(params)
			if err != nil {
				//fmt.Println(err.Error())
				continue
			}

			if len(resp.Messages) > 0 {
				fmt.Println(resp)

				for _, message := range resp.Messages {

					var sr SendRequest
					if err := json.Unmarshal([]byte(*message.Body), &sr); err == nil {
						go func() {
							callback(&sr)
						}()
					} else {
						fmt.Println(err)
					}

					delete_params := &sqs.DeleteMessageInput{
						QueueUrl:      aws.String(url),       // Required
						ReceiptHandle: message.ReceiptHandle, // Required

					}
					_, err := conn.DeleteMessage(delete_params) // No response returned when successed.
					if err != nil {
						fmt.Println(err)
					}
					fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted.\n\n", *message.MessageId)
				}
			}
		}
	}()
}
