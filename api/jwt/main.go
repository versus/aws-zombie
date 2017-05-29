package main

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	"github.com/pkg/errors"
)

func Handle(evt map[string]string, ctx *runtime.Context) (interface{}, error) {
	token, _ := jwt.Parse(evt[`authorizationToken`], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(`covabunga`), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		p := policy{
			PrincipalId: `user`,
			PolicyDocument: policyDoc{
				Version: `2012-10-17`,
				Statement: []policyStatement{
					{
						Action:   `execute-api:Invoke`,
						Effect:   `Allow`,
						Resource: evt[`methodArn`],
					},
				},
			},
		}

		//j, _ := json.Marshal(p)

		return p, nil
	}

	return "", errors.New(`Error: Invalid token`)
}

type policy struct {
	PrincipalId    string    `json:"principalId"`
	PolicyDocument policyDoc `json:"policyDocument"`
}

type policyDoc struct {
	Version   string            `json:"Version"`
	Statement []policyStatement `json:"Statement"`
}

type policyStatement struct {
	Action   string `json:"Action"`
	Effect   string `json:"Effect"`
	Resource string `json:"Resource"`
}

func main() {
	m := map[string]string{
		"type":               "TOKEN",
		"authorizationToken": "tockerstring",
		"methodArn":          "arn:aws:execute-api:eu-west-1:xxxxxxx:xxxxxxxxxxx/null/GET/",
	}

	ctx := runtime.Context{}

	Handle(m, &ctx)
}
