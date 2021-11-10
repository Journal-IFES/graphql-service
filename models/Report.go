package models

import (
	"github.com/Journal-IFES/graphql-service/internal/postgres"
	"github.com/graphql-go/graphql"
)

type ReportStruct struct {
	Id      int    `json:"id"`
	By      int    `json:"by"`
	To      int    `json:"to"`
	Comment string `json:"comment"`
	Review  int    `json:"review"`
}

func GetReportById(id int) (*ReportStruct, error) {

	db := postgres.GetPostgresDB()

	rows, err := db.Query(`
	SELECT * 
	FROM report 
	WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	r := ReportStruct{}

	rows.Next()
	err = rows.Scan(
		&r.Id,
		&r.By,
		&r.To,
		&r.Comment,
		&r.Review,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

var ReportType graphql.Object = *graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Report",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Name: "id",
				Type: graphql.NewNonNull(graphql.Int),
			},
			"by": &graphql.Field{
				Name: "by",
				Type: graphql.NewNonNull(&UserType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(*ReportStruct).By

					return GetUserById(id)
				},
			},
			"to": &graphql.Field{
				Name: "to",
				Type: graphql.NewNonNull(&UserType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(*ReportStruct).To

					return GetUserById(id)
				},
			},
			"comment": &graphql.Field{
				Name: "comment",
				Type: graphql.NewNonNull(graphql.String),
			},
			"review": &graphql.Field{
				Name: "review",
				Type: &ReviewType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Source.(*ReportStruct).Review

					return GetReviewById(id)
				},
			},
		},
	})

var ReportField graphql.Field = graphql.Field{
	Name: "report",
	Type: &ReportType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(int)

		return GetReportById(id)
	},
}
