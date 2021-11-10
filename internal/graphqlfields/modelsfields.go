package graphqlfields

import (
	"github.com/Journal-IFES/graphql-service/models"
	"github.com/graphql-go/graphql"
)

func ModelsFields() *graphql.Fields {
	return &graphql.Fields{
		"User":      &models.UserField,
		"Review":    &models.ReviewField,
		"Report":    &models.ReportField,
		"Content":   &models.ContentField,
		"Article":   &models.ArticleField,
		"BasicAuth": &models.BasicAuthField,
	}
}

func ModelsMutations() *graphql.Fields {
	a := graphql.Fields(make(map[string]*graphql.Field))

	for k, v := range models.UserMutation {
		a[k] = v
	}

	return &a
}
