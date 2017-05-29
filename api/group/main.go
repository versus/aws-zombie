package main
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/satori/go.uuid"
	"time"
)


func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (string, error) {
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(evt.Body), &data); err != nil {
		return ``, err
	}

	if err := validateInput(data); err != nil {
		return "", err
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))

	svc := dynamodb.New(sess)

	uID := uuid.NewV4()

	item := dynamodb.PutItemInput{
		TableName: aws.String(`group`),
		Item: map[string]*dynamodb.AttributeValue{
			`id`:        {S: aws.String(uID.String())},
			`name`:      {S: aws.String(data[`name`].(string))},
			`timestamp`: {N: aws.String(fmt.Sprintf("%d", time.Now().Unix()))},
		},
	}
	if _, err := svc.PutItem(&item); err != nil {
		return err.Error(), err
	}

	return "Ok", nil
}

func validateInput(data map[string]interface{}) error {
	if _, ok := data[`name`]; !ok {
		return errors.New(`Empty name`)
	}

	return nil
}
