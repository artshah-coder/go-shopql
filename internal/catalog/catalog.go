package catalog

import (
	"github.com/artshah-coder/go-shopql/internal/item"
)

type Catalog struct {
	ID       uint32            `json:"id"`
	Name     string            `json:"name"`
	ParentID *uint32           `json:"parent_id"`
	Childs   CatalogsStorage   `json:"childs"`
	Items    item.ItemsStorage `json:"items"`
}

// func (catalog *Catalog) Id() string {
// 	return strconv.Itoa(int(catalog.ID))
// }
