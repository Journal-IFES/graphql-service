package models

import (
	"github.com/Journal-IFES/graphql-service/internal/postgres"
	"github.com/graphql-go/graphql"
)

type ReviewStruct struct {
	Id      int    `json:"id"`
	By      int    `json:"by"`
	Comment string `json:"comment"`
}

func GetReviewById(id int) (*ReviewStruct, error) {

	db := postgres.GetPostgresDB()

	rows, err := db.Query(`
	SELECT * 
	FROM review 
	WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	r := ReviewStruct{}

	rows.Next()
	err = rows.Scan(
		&r.Id,
		&r.By,
		&r.Comment,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

var ReviewType graphql.Object = *graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Review",
		Description: "Definição do tipo revisão no banco de dados.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Name:        "id",
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Identificador único da revisão no banco de dados.",
			},
			"by": &graphql.Field{
				Name:        "by",
				Type:        graphql.NewNonNull(&UserType),
				Description: "Identificador único do usuário que fez a revisão.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(*ReviewStruct).By

					return GetUserById(id)
				},
			},
			"comment": &graphql.Field{
				Name:        "comment",
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Comentário referente ao que foi revisado.",
			},
		},
	})

var ReviewField graphql.Field = graphql.Field{
	Name:        "review",
	Type:        &ReviewType,
	Description: "Busca revisão no banco de dados.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(int)

		return GetReviewById(id)
	},
}
