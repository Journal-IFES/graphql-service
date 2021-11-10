package models

import (
	"github.com/Journal-IFES/graphql-service/internal/postgres"
	"github.com/graphql-go/graphql"
)

type BasicAuthStruct struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	User     int    `json:"user"`
}

func GetBasicAuthById(id int) (*BasicAuthStruct, error) {

	db := postgres.GetPostgresDB()

	rows, err := db.Query(`
	SELECT * 
	FROM basic_auth 
	WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	r := BasicAuthStruct{}

	rows.Next()
	err = rows.Scan(
		&r.Id,
		&r.Username,
		&r.Password,
		&r.User,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

var BasicAuthType graphql.Object = *graphql.NewObject(graphql.ObjectConfig{
	Name: "BasicAuth",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "id",
			Type: graphql.NewNonNull(graphql.Int),
		},
		"username": &graphql.Field{
			Name: "username",
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.Field{
			Name: "password",
			Type: graphql.NewNonNull(graphql.String),
		},
		"user": &graphql.Field{
			Name: "user",
			Type: graphql.NewNonNull(&UserType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Source.(*BasicAuthStruct).User

				return GetUserById(id)
			},
		},
	},
})

var BasicAuthField graphql.Field = graphql.Field{
	Name: "basicAuth",
	Type: &BasicAuthType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(int)

		return GetBasicAuthById(id)
	},
}
