package app

import (
	"net/http"

	"github.com/artshah-coder/go-shopql/internal/catalog"
	"github.com/artshah-coder/go-shopql/internal/graphql"
	"github.com/artshah-coder/go-shopql/internal/item"
	"github.com/artshah-coder/go-shopql/internal/seller"
	"github.com/artshah-coder/go-shopql/internal/session"
	"github.com/artshah-coder/go-shopql/internal/user"
	"github.com/artshah-coder/go-shopql/internal/utils/dataloader"
	"github.com/gin-gonic/gin"
)

func GetApp() http.Handler {
	app := gin.Default()

	uStor := user.NewUserStMem()
	sessStor := session.NewSessionStMem()
	catStor := catalog.NewCatalogStMem()
	sellStor := seller.NewSellerStMem()
	itmsStor := item.NewItemStMem()
	dataloader.MustDataLoad("../test/data/testdata.json", catStor, itmsStor, sellStor)

	uh := &user.UserHandler{
		Users:    uStor,
		Sessions: sessStor,
	}
	app.POST("/register", uh.Register)
	app.Use(session.UIDToContext(sessStor))
	app.POST("/query", graphql.GraphqlHandler(uStor, catStor, itmsStor, sellStor))
	return app
}
