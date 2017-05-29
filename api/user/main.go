package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/apigatewayproxyevt"
	"github.com/ivch/aws-zombie/api/user/handlers"
	"github.com/pkg/errors"
)

func Handle(evt *apigatewayproxyevt.Event, _ *runtime.Context) (interface{}, error) {
	h := handlers.NewHandler()

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(evt.Body), &m); err != nil {
		return "", err
	}

	switch evt.HTTPMethod {
	case "GET":
		return h.Get(m)
	case "PUT":
		return h.Put(m, evt.QueryStringParameters[`phone`])
	case "POST":
		return h.Post(m)
	}

	return "", errors.New(`Unsupported request method`)
}

type user struct {
	Phone     string    `json:"phone"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Specs     []string  `json:"specs"`
	Lat       string    `json:"lat"`
	Long      string    `json:"long"`
	CreatedAt time.Time `json:"created_at"`
	Groups    []string  `json:"groups"`
	City      string    `json:"string"`
	Country   string    `json:"string"`
}

func main() {
	u := user{
		Phone: `380636637354`,
		//		//FirstName: `John`,
		//		//LastName:  `Doe`,
		//		//Password:  `p@$$w0rd`,
		//		//Specs:     []string{`medic`, `hunter`},
		//		//Lat:       `63.860036`,
		//		//Long:      `-75.761719`,
		//		//CreatedAt: time.Now(),
		//		//Groups:    []string{`LosSolomas`, `POH`},
	}

	ctx := runtime.Context{}

	j, err := json.Marshal(u)
	if err != nil {
		log.Fatal()
	}

	//fmt.Println(j)

	e := apigatewayproxyevt.Event{
		HTTPMethod: "GET",
		Body:       string(j),
	}

	d, err := Handle(&e, &ctx)

	fmt.Println(d)
	fmt.Println(err)
}
