package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

func main() {
	ctx := runtime.Context{}
	Handle("", &ctx)
}

func Handle(evt interface{}, ctx *runtime.Context) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))
	svc := dynamodb.New(sess)

	item := dynamodb.PutItemInput{
		TableName: aws.String(`zTestTable`),
		Item: map[string]*dynamodb.AttributeValue{
			`name`:   {S: aws.String(`John Doe`)},
			`geo`:    {S: aws.String(strconv.Itoa(rand.Intn(10)))},
			`phone`:  {S: aws.String(`555-55-55`)},
			`skills`: {SS: []*string{aws.String(`medic`), aws.String(`hunter`)}},
		},
	}
	resp, err := svc.PutItem(&item)

	return fmt.Sprint(resp), err
}
