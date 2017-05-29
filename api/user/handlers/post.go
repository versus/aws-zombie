package handlers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"reflect"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (h *handler) Post(rq map[string]interface{}) (string, error) {
	if err := h.validateInput(rq); err != nil {
		return "", err
	}

	if err := h.isExistingUser(fmt.Sprint(rq[`phone`])); err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write([]byte(fmt.Sprint(rq[`password`])))

	var specs, groups []*string

	switch reflect.TypeOf(rq[`specs`]).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(rq[`specs`])

		for i := 0; i < s.Len(); i++ {
			specs = append(specs, aws.String(fmt.Sprint(s.Index(i))))
		}
	}

	if v, ok := rq[`groups`]; ok && v != nil {
		switch reflect.TypeOf(rq[`groups`]).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(rq[`groups`])

			for i := 0; i < s.Len(); i++ {
				groups = append(groups, aws.String(fmt.Sprint(s.Index(i))))
			}
		}
	}

	item := dynamodb.PutItemInput{
		TableName: aws.String(`user`),
		Item: map[string]*dynamodb.AttributeValue{
			`phone`:      {S: aws.String(fmt.Sprint(rq[`phone`]))},
			`first_name`: {S: aws.String(fmt.Sprint(rq[`first_name`]))},
			`last_name`:  {S: aws.String(fmt.Sprint(rq[`last_name`]))},
			`password`:   {S: aws.String(fmt.Sprintf(`%x`, hash.Sum(nil)))},
			`specs`:      {SS: specs},
			`created_at`: {S: aws.String(fmt.Sprint(time.Now()))},
		},
	}

	if v, ok := rq[`push_id`]; ok && v != nil && v != "" {
		item.Item[`push_id`] = &dynamodb.AttributeValue{S: aws.String(fmt.Sprint(rq[`push_id`]))}
	}
	if v, ok := rq[`lat`]; ok && v != nil && v != "" {
		item.Item[`lat`] = &dynamodb.AttributeValue{S: aws.String(fmt.Sprint(rq[`lat`]))}
	}
	if v, ok := rq[`long`]; ok && v != nil && v != "" {
		item.Item[`long`] = &dynamodb.AttributeValue{S: aws.String(fmt.Sprint(rq[`long`]))}
	}
	if v, ok := rq[`city`]; ok && v != nil && v != "" {
		item.Item[`city`] = &dynamodb.AttributeValue{S: aws.String(fmt.Sprint(rq[`city`]))}
	}
	if v, ok := rq[`country`]; ok && v != nil && v != "" {
		item.Item[`country`] = &dynamodb.AttributeValue{S: aws.String(fmt.Sprint(rq[`country`]))}
	}

	if len(groups) > 0 {
		item.Item[`groups`] = &dynamodb.AttributeValue{SS: groups}
	}

	_, err := h.db.PutItem(&item)

	return "", err
}

func (h *handler) isExistingUser(phone string) error {
	u, err := h.getUser(phone)
	if err != nil {
		return err
	}

	if _, ok := u.Item[`phone`]; ok {
		return errors.New("User already exists")
	}

	return nil
}

func (*handler) validateInput(data map[string]interface{}) error {
	if _, ok := data[`phone`]; !ok {
		return errors.New(`Empty phone`)
	}

	if _, ok := data[`first_name`]; !ok {
		return errors.New(`Empty first name`)
	}

	if _, ok := data[`last_name`]; !ok {
		return errors.New(`Empty last name`)
	}

	if _, ok := data[`password`]; !ok {
		return errors.New(`Empty password`)
	}

	if _, ok := data[`specs`]; !ok {
		return errors.New(`Empty specs`)
	}

	return nil
}
