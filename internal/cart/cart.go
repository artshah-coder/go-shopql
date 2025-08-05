package cart

import (
	"fmt"
	"reflect"

	"github.com/artshah-coder/go-shopql/internal/item"
)

var (
	ErrNoInCartItem = fmt.Errorf("no item in cart")
	ErrNoCartInMem  = fmt.Errorf("types mismatch")
)

type InCartItem struct {
	Item     *item.Item
	Quantity uint32
}

type Cart interface {
	Add(*InCartItem) (Cart, error)
	GetByItemId(uint32) (*InCartItem, error)
	GetAll() (Cart, error)
	Delete(InCartItem) (Cart, error)
}

type CartInMem []*InCartItem

func NewCartInMem() *CartInMem {
	cartInMem := CartInMem(make([]*InCartItem, 0, 10))
	return &cartInMem
}

func (cart *CartInMem) Add(item *InCartItem) (Cart, error) {
	for _, itm := range *cart {
		if itm.Item.ID == item.Item.ID {
			// itm.Item = item.Item
			itm.Quantity += item.Quantity
			return cart.GetAll()
		}
	}
	*cart = append(*cart, item)
	return cart.GetAll()
}

func (cart *CartInMem) GetByItemId(itemID uint32) (*InCartItem, error) {
	for _, item := range *cart {
		if item.Item.ID == itemID {
			return item, nil
		}
	}
	return nil, ErrNoInCartItem
}

func (cart *CartInMem) GetAll() (Cart, error) {
	return cart, nil
}

func (cart *CartInMem) Delete(item InCartItem) (Cart, error) {
	for i, itm := range *cart {
		if reflect.DeepEqual(itm, &item) {
			if i == len(*cart)-1 {
				*cart = (*cart)[:i]
			} else {
				*cart = append((*cart)[:i], (*cart)[i+1:]...)
			}
		}
	}
	return cart.GetAll()
}
