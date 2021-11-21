package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Journal-IFES/graphql-service/internal/postgres"
	"github.com/graphql-go/graphql"
)

type UserStruct struct {
	Id        int       `json:"id"`
	Firstname string    `json:"firstname"`
	Surname   string    `json:"surname"`
	Join_date time.Time `json:"join_date"`
	Hierarchy bool      `json:"hierarchy"`
	Activated bool      `json:"activated"`
}

func GetUserById(id int) (*UserStruct, error) {
	db := postgres.GetPostgresDB()

	rows, err := db.Query(`
	SELECT * 
	FROM userr 
	WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	r := UserStruct{}

	rows.Next()
	err = rows.Scan(
		&r.Id,
		&r.Firstname,
		&r.Surname,
		&r.Join_date,
		&r.Hierarchy,
		&r.Activated,
	)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

var UserType graphql.Object = *graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "User",
		Description: "Definição do tipo usuário no banco de dados.",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Name:        "id",
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Identificador único do usuário no banco de dados.",
			},
			"firstname": &graphql.Field{
				Name:        "name",
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Primeiro nome do usuário.",
			},
			"surname": &graphql.Field{
				Name:        "surname",
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Sobrenome do usuário.",
			},
			"join_date": &graphql.Field{
				Name:        "join_date",
				Type:        graphql.NewNonNull(graphql.DateTime),
				Description: "Data de registro do usuário.",
			},
			"hierarchy": &graphql.Field{
				Name:        "hierarchy",
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Hierarquia do usuário. Define as permissões do mesmo.",
			},
			"activated": &graphql.Field{
				Name:        "activated",
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Define se o usuário está ativo ou não para utilização do sistema.",
			},
		},
	},
)

var UserField graphql.Field = graphql.Field{
	Name:        "User",
	Type:        &UserType,
	Description: "Busca usuário no banco de dados.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Identificador único do usuário no banco de dados.",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id := p.Args["id"].(int)

		return GetUserById(id)
	},
}

var UserMutation graphql.Fields = graphql.Fields{
	"newUser": &graphql.Field{
		Name:        "NewUser",
		Description: "Cria usuário no banco de dados.",
		Type:        &UserType,
		Args: graphql.FieldConfigArgument{
			"firstname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Primeiro nome do usuário.",
			},
			"surname": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Sobrenome do usuário.",
			},
			"hierarchy": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Hierarquia do usuário. Define as permissões do mesmo.",
			},
			"activated": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "Define se o usuário está ativo ou não para utilização do sistema.",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			db := postgres.GetPostgresDB()

			rows, err := db.Query(`
			INSERT INTO userr (
				firstname,
				surname,
				join_date,
				hierarchy,
				activated
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			)
			RETURNING *
			`,
				p.Args["firstname"].(string),
				p.Args["surname"].(string),
				time.Now().Format("2006-01-02"),
				p.Args["hierarchy"].(int),
				p.Args["activated"].(bool),
			)
			if err != nil {
				return nil, err
			}

			r := UserStruct{}

			rows.Next()
			rows.Scan(
				&r.Id,
				&r.Firstname,
				&r.Surname,
				&r.Join_date,
				&r.Hierarchy,
				&r.Activated,
			)

			return &r, nil
		},
	},
	"deleteUser": &graphql.Field{
		Name:        "DeleteUser",
		Description: "Deleta usuário no banco de dados.",
		Type:        &UserType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			db := postgres.GetPostgresDB()

			rows, err := db.Query(`
			DELETE FROM userr
			WHERE id = $1
			RETURNING *
			`,
				p.Args["id"].(int),
			)
			if err != nil {
				return nil, err
			}

			r := UserStruct{}

			rows.Next()
			rows.Scan(
				&r.Id,
				&r.Firstname,
				&r.Surname,
				&r.Join_date,
				&r.Hierarchy,
				&r.Activated,
			)

			return &r, nil
		},
	},
	"updateUser": &graphql.Field{
		Name:        "UpdateUser",
		Description: "Atualiza usuário no banco de dados.",
		Type:        &UserType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Identificador único do usuário no banco de dados.",
			},
			"firstname": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Primeiro nome do usuário.",
			},
			"surname": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Sobrenome do usuário.",
			},
			"hierarchy": &graphql.ArgumentConfig{
				Type:        graphql.Int,
				Description: "Hierarquia do usuário. Define as permissões do mesmo.",
			},
			"activated": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "Define se o usuário está ativo ou não para utilização do sistema.",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			db := postgres.GetPostgresDB()

			query := "UPDATE userr\nSET "
			for k, v := range p.Args {
				if k != "id" {
					switch v := v.(type) {
					case string:
						query += fmt.Sprintf("%s = '%s'\n", k, v)
					case int:
						query += fmt.Sprintf("%s = '%d'\n", k, v)
					case bool:
						query += fmt.Sprintf("%s = '%s'\n", k, strconv.FormatBool(v))
					}
				}
			}

			query += "WHERE id = $1\nRETURNING *"

			rows, err := db.Query(query,
				p.Args["id"].(int),
			)
			if err != nil {
				return nil, err
			}

			r := UserStruct{}

			rows.Next()
			rows.Scan(
				&r.Id,
				&r.Firstname,
				&r.Surname,
				&r.Join_date,
				&r.Hierarchy,
				&r.Activated,
			)

			return &r, nil
		},
	},
	"listUsers": &graphql.Field{
		Name:        "listUsers",
		Description: "Lista os usuários no banco de dados.",
		Type:        graphql.NewList(&UserType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			db := postgres.GetPostgresDB()

			rows, err := db.Query(`
			SELECT *
			FROM userr
			`)
			if err != nil {
				return nil, err
			}

			lst := make([]*UserStruct, 0)

			for rows.Next() {
				r := UserStruct{}

				rows.Scan(
					&r.Id,
					&r.Firstname,
					&r.Surname,
					&r.Join_date,
					&r.Hierarchy,
					&r.Activated,
				)

				lst = append(lst, &r)
			}

			return &lst, nil
		},
	},
}
