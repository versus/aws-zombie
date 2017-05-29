package main

import (
	"encoding/json"
	"errors"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
)

func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(evt.Body), &m); err != nil {
		return "", err
	}

	switch evt.HTTPMethod {
	case "GET":
		return get(m)
	}

	return "", errors.New(`Unsupported request method`)
}

func get(rq map[string]interface{}) (interface{}, error) {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))

	db := dynamodb.New(sess)

	q := dynamodb.ScanInput{
		TableName: aws.String(`user`),
	}

	res, err := db.Scan(&q)
	if err != nil {
		return "", err
	}

	var items []user
	for _, item := range res.Items {
		if i, err := mapUserData(item); err == nil {
			items = append(items, i)
		}
	}

	return items, nil
}

func mapUserData(data map[string]*dynamodb.AttributeValue) (u user, err error) {
	defer func() {
		if r := recover(); r != nil {
			u = user{}
			err = errors.New("Error retrieving data")
		}
	}()

	u = user{
		Phone:     *data[`phone`].S,
		FirstName: *data[`first_name`].S,
		LastName:  *data[`last_name`].S,
		Specs:     []string{},
	}

	if _, ok := data[`lat`]; ok {
		u.Lat = *data[`lat`].S
	}

	if _, ok := data[`long`]; ok {
		u.Long = *data[`long`].S
	}

	if _, ok := data[`city`]; ok {
		u.City = *data[`city`].S
	}
	if _, ok := data[`country`]; ok {
		u.Country = *data[`country`].S
	}
	if _, ok := data[`push_id`]; ok {
		u.Lat = *data[`push_id`].S
	}

	if _, ok := data[`groups`]; ok {
		for _, v := range data[`groups`].SS {
			u.Groups = append(u.Groups, *v)
		}
	}

	for _, v := range data[`specs`].SS {
		u.Specs = append(u.Specs, *v)
	}

	return u, nil
}

//func main() {
//	ctx := runtime.Context{}
//
//	e := apigatewayproxyevt.Event{
//		HTTPMethod: "GET",
//		Body:       string("{}"),
//	}
//
//	d, err := Handle(&e, &ctx)
//
//	fmt.Println(d)
//	fmt.Println(err)
//}

type user struct {
	Phone     string    `json:"phone"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password,omitempty"`
	Specs     []string  `json:"specs"`
	Lat       string    `json:"lat"`
	Long      string    `json:"long"`
	CreatedAt time.Time `json:"created_at"`
	Groups    []string  `json:"groups"`
	City      string    `json:"string"`
	Country   string    `json:"string"`
	PushId    string    `json:"push_id"`
}
