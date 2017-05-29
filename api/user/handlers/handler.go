package handlers

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Handler interface {
	Post(data map[string]interface{}) (string, error)
	Get(data map[string]interface{}) (string, error)
	Put(data map[string]interface{}, phone string) (string, error)
}

type handler struct {
	db *dynamodb.DynamoDB
}

type user struct {
	Pk          string    `json:"pk"`
	Phone       string    `json:"phone"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Password    string    `json:"password,omitempty"`
	Specs       []string  `json:"specs"`
	Lat         string    `json:"lat"`
	Long        string    `json:"long"`
	CreatedAt   time.Time `json:"created_at"`
	Groups      []string  `json:"groups"`
	City        string    `json:"string"`
	Country     string    `json:"string"`
	PushId      string    `json:"push_id"`
	Contacts    []string  `json:"contacts"`
	ContactList []user    `json:"contact_list"`
}

func NewHandler() Handler {
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))

	return &handler{db: dynamodb.New(sess)}
}
