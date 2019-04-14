package main

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Env struct {
	db *gorm.DB
}

func (e *Env) migrateDB() {
	e.db.AutoMigrate(&Rate{})
	e.db.AutoMigrate(&ExchangeRate{})
}

func (e *Env) clearDB() {
	e.db.Exec("DELETE FROM rates")
	e.db.Exec("ALTER SEQUENCE rates_id_seq RESTART WITH 1")
	e.db.Exec("DELETE FROM exchange_rates")
	e.db.Exec("ALTER SEQUENCE exchange_rates_id_seq RESTART WITH 1")
}

func (e *Env) run() {
	r := gin.Default()
	r.POST("/graphql", graphqlHandler(e))
	r.GET("/graphql/playground", playgroundHandler())
	r.Run()
}

func graphqlHandler(e *Env) gin.HandlerFunc {
	h := handler.GraphQL(NewExecutableSchema(
		Config{
			Resolvers: &Resolver{
				Env: *e,
			},
		},
	))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
