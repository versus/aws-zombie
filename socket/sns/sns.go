package sns

import (
	"encoding/json"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func New() *sns.SNS {
	creds := credentials.NewStaticCredentials(
		"piblic-key",
		"private key",
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

	return sns.New(sess)
}

type Data struct {
	Message string `json:"message"`
}

type GCM struct {
	Data Data `json:"data"`
}

type Message struct {
	GCM string `json:"GCM"`
}

func SendPush(conn *sns.SNS, pushID string) {
	resp, err := conn.CreatePlatformEndpoint(&sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String("arn"),
		Token: aws.String(pushID),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//m, err := newMessageJSON(data)
	//if err != nil {
	//	return
	//}

	bbb, _ := json.Marshal(GCM{
		Data: Data{
			Message: "You have new message!",
		},
	})

	msg := Message{
		GCM: string(bbb),
	}

	bts, _ := json.Marshal(msg)

	fmt.Println(string(bts))

	input := &sns.PublishInput{
		Message:          aws.String(string(bts)),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(*resp.EndpointArn),
	}
	_, err = conn.Publish(input)
	if err != nil {
		fmt.Println(err)
	}

	return
}
