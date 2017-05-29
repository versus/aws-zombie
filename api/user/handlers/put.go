package handlers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (h *handler) Put(data map[string]interface{}, phone string) (string, error) {
	q := dynamodb.UpdateItemInput{
		TableName: aws.String(`user`),
		Key: map[string]*dynamodb.AttributeValue{
			`phone`: {S: aws.String(phone)},
		},
		UpdateExpression: aws.String(`set push_id = :p1`),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			`:p1`: {S: aws.String(fmt.Sprint(data[`push_id`]))},
		},
	}

	_, err := h.db.UpdateItem(&q)

	return "", err
}
