package graphql

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/artshah-coder/go-shopql/internal/catalog"
	"github.com/artshah-coder/go-shopql/internal/item"
	"github.com/artshah-coder/go-shopql/internal/seller"
	"github.com/artshah-coder/go-shopql/internal/session"
	"github.com/artshah-coder/go-shopql/internal/user"
	"github.com/gin-gonic/gin"
)

var ErrNotAuth = fmt.Errorf("user not authorized")

func GraphqlHandler(
	users user.UsersStorage, catalogs catalog.CatalogsStorage, items item.ItemsStorage,
	sellers seller.SellersStorage,
) gin.HandlerFunc {
	cfg := Config{Resolvers: &Resolver{
		USt:    users,
		ISt:    items,
		CSt:    catalogs,
		SellSt: sellers,
	}}

	cfg.Directives.Shopql_authorized = func(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
		uid, ok := session.UIDFromContext(ctx)
		if !ok {
			return nil, ErrNotAuth
		}
		_, err = users.GetById(uid)
		if err != nil {
			return nil, err
		}
		return next(ctx)
	}

	gh := handler.New(NewExecutableSchema(cfg))
	gh.AddTransport(transport.POST{})
	gh.Use(extension.FixedComplexityLimit(100))

	return func(c *gin.Context) {
		gh.ServeHTTP(c.Writer, c.Request)
	}
}
