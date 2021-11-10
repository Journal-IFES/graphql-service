package models

import (
	"github.com/Journal-IFES/graphql-service/internal/postgres"
	"github.com/graphql-go/graphql"
)

type ContentStruct struct {
	Id       int         `json:"id"`
	Title    string      `json:"title"`
	Subtitle interface{} `json:"subtitle"`
	Body     string      `json:"body"`
	Author   int         `json:"author"`
}

func GetContentById(id int) (*ContentStruct, error) {

	db := postgres.GetPostgresDB()

	rows, err := db.Query(`
	SELECT * 
	FROM content 
	WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	r := ContentStruct{}

	rows.Next()
	err = rows.Scan(
		&r.Id,
		&r.Title,
		&r.Subtitle,
		&r.Body,
		&r.Author,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

var ContentType graphql.Object = *graphql.NewObject(graphql.ObjectConfig{
	Name: "Content",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "id",
			Type: graphql.NewNonNull(graphql.Int),
		},
		"title": &graphql.Field{
			Name: "title",
			Type: graphql.NewNonNull(graphql.String),
		},
		"subtitle": &graphql.Field{
			Name: "subtitle",
			Type: graphql.NewNonNull(graphql.String),
		},
		"body": &graphql.Field{
			Name: "body",
			Type: graphql.NewNonNull(graphql.String),
		},
		"author": &graphql.Field{
			Name: "author",
			Type: graphql.NewNonNull(&UserType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Source.(*ContentStruct).Author

				return GetUserById(id)
			},
		},
	},
})

var ContentField graphql.Field = graphql.Field{
	Name: "content",
	Type: &ContentType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(int)

		return GetContentById(id)
	},
}
