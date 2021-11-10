package models

import (
	"time"

	"github.com/Journal-IFES/graphql-service/internal/postgres"
	"github.com/graphql-go/graphql"
)

type ArticleStruct struct {
	Id           int       `json:"id"`
	Author       int       `json:"author"`
	Publish_date time.Time `json:"publish_date"`
	Content      int       `json:"content"`
	Rating       float64   `json:"rating"`
	Published    bool      `json:"published"`
}

func GetArticleById(id int) (*ArticleStruct, error) {

	db := postgres.GetPostgresDB()

	rows, err := db.Query(`
	SELECT * 
	FROM article 
	WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	r := ArticleStruct{}

	rows.Next()
	err = rows.Scan(
		&r.Id,
		&r.Author,
		&r.Publish_date,
		&r.Content,
		&r.Rating,
		&r.Published,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

var ArticleType graphql.Object = *graphql.NewObject(graphql.ObjectConfig{
	Name: "Article",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Name: "id",
			Type: graphql.NewNonNull(graphql.Int),
		},
		"author": &graphql.Field{
			Name: "author",
			Type: graphql.NewNonNull(&UserType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Source.(*ArticleStruct).Author

				return GetUserById(id)
			},
		},
		"publish_date": &graphql.Field{
			Name: "publish_date",
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"content": &graphql.Field{
			Name: "content",
			Type: graphql.NewNonNull(&ContentType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id := p.Source.(*ArticleStruct).Content

				return GetContentById(id)
			},
		},
		"rating": &graphql.Field{
			Name: "rating",
			Type: graphql.NewNonNull(graphql.Float),
		},
		"published": &graphql.Field{
			Name: "published",
			Type: graphql.NewNonNull(graphql.Boolean),
		},
	},
})

var ArticleField graphql.Field = graphql.Field{
	Name: "article",
	Type: &ArticleType,
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(int)

		return GetArticleById(id)
	},
}
