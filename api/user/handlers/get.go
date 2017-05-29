package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func (h *handler) Get(rq map[string]interface{}) (string, error) {

	user, err := h.getUser(fmt.Sprint(rq[`phone`]))
	if err != nil {
		return "", err
	}

	d, err := h.mapUserData(user)
	if err != nil {
		return "", err
	}

	h.getUserContacts(&d)

	fmt.Println(d)

	s, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func (h *handler) mapUserData(data *dynamodb.GetItemOutput) (u user, err error) {
	defer func() {
		if r := recover(); r != nil {
			u = user{}
			err = errors.New("Error retrieving data")
		}
	}()

	u = user{
		Phone:     *data.Item[`phone`].S,
		FirstName: *data.Item[`first_name`].S,
		LastName:  *data.Item[`last_name`].S,
		Specs:     []string{},
	}

	if _, ok := data.Item[`lat`]; ok {
		u.Lat = *data.Item[`lat`].S
	}

	if _, ok := data.Item[`long`]; ok {
		u.Long = *data.Item[`long`].S
	}

	if _, ok := data.Item[`city`]; ok {
		u.City = *data.Item[`city`].S
	}
	if _, ok := data.Item[`country`]; ok {
		u.Country = *data.Item[`country`].S
	}
	if _, ok := data.Item[`push_id`]; ok {
		u.PushId = *data.Item[`push_id`].S
	}

	if _, ok := data.Item[`groups`]; ok {
		for _, v := range data.Item[`groups`].SS {
			u.Groups = append(u.Groups, *v)
		}
	}

	if _, ok := data.Item[`contacts`]; ok {
		for _, v := range data.Item[`contacts`].SS {
			u.Contacts = append(u.Contacts, *v)
		}
	}

	for _, v := range data.Item[`specs`].SS {
		u.Specs = append(u.Specs, *v)
	}

	return u, nil
}

func (h *handler) getUserContacts(u *user) {
	if len(u.Contacts) > 0 {
		for _, c := range u.Contacts {
			user, err := h.getUser(c)
			if err != nil {
				continue
			}

			d, err := h.mapUserData(user)
			if err != nil {
				continue
			}

			u.ContactList = append(u.ContactList, d)
		}
	}
}

func (h *handler) getUser(phone string) (*dynamodb.GetItemOutput, error) {
	params := dynamodb.GetItemInput{
		TableName: aws.String(`user`),
		Key: map[string]*dynamodb.AttributeValue{
			`phone`: {S: aws.String(phone)},
		},
	}

	res, err := h.db.GetItem(&params)
	if err != nil {
		return nil, err
	}

	return res, nil
}
