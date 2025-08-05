package seller

import "github.com/artshah-coder/go-shopql/internal/item"

type Seller struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Deals uint32 `json:"deals"`
	Items item.ItemsStorage
}

// func (seller *Seller) Id() string {
// 	return strconv.Itoa(int(seller.ID))
// }
