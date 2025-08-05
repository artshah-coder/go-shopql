package seller

import (
	"fmt"

	"github.com/artshah-coder/go-shopql/internal/item"
)

var (
	ErrSellerExists   = fmt.Errorf("seller already exists")
	ErrSellerNotExist = fmt.Errorf("seller does not exist")
)

type SellersStorage interface {
	Add(Seller) (*Seller, error)
	GetAll() ([]*Seller, error)
	GetById(uint32) (*Seller, error)
	DeleteById(uint32) error
}

type SellerStMem struct {
	Sellers []*Seller
}

func NewSellerStMem() *SellerStMem {
	return &SellerStMem{
		Sellers: make([]*Seller, 0, 10),
	}
}

func (sSt *SellerStMem) Add(seller Seller) (*Seller, error) {
	for _, s := range sSt.Sellers {
		if s.ID == seller.ID || s.Name == seller.Name {
			return nil, ErrSellerExists
		}
	}
	if seller.Items == nil {
		seller.Items = item.NewItemStMem()
	}
	sSt.Sellers = append(sSt.Sellers, &seller)
	return &seller, nil
}

func (sSt *SellerStMem) GetAll() ([]*Seller, error) {
	return sSt.Sellers, nil
}

func (sSt *SellerStMem) GetById(sellerID uint32) (*Seller, error) {
	for _, seller := range sSt.Sellers {
		if seller.ID == sellerID {
			return seller, nil
		}
	}
	return nil, ErrSellerNotExist
}

func (sSt *SellerStMem) DeleteById(sellerID uint32) error {
	for i, seller := range sSt.Sellers {
		if seller.ID == sellerID {
			if i == len(sSt.Sellers)-1 {
				sSt.Sellers = sSt.Sellers[:i]
			} else {
				sSt.Sellers = append(sSt.Sellers[:i], sSt.Sellers[i+1:]...)
			}
			return nil
		}
	}
	return ErrSellerNotExist
}
