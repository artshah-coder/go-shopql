package item

import "fmt"

var (
	ErrItemNotExist    = fmt.Errorf("item does not exist")
	ErrIncorrectBounds = fmt.Errorf("incorrect limit or offset")
	ErrNotEnoughStock  = fmt.Errorf("not enough stock")
)

type ItemsStorage interface {
	Add(*Item) (*Item, error)
	Populate([]*Item) ([]*Item, error)
	GetAll(uint32, uint32) ([]*Item, error)
	GetById(uint32) (*Item, error)
	GetBySellerId(uint32) ([]*Item, error)
	Move(uint32, uint32) (*Item, error)
	//	Return
	DeleteById(uint32) error
}

type ItemStMem struct {
	Items []*Item
}

func NewItemStMem() *ItemStMem {
	return &ItemStMem{
		Items: make([]*Item, 0, 10),
	}
}

func (iSt *ItemStMem) Add(item *Item) (*Item, error) {
	for _, i := range iSt.Items {
		if i.Name == item.Name && i.SellerID == item.SellerID {
			i.InStock += item.InStock
			return i, nil
		}
	}

	iSt.Items = append(iSt.Items, item)
	return item, nil
}

func (iSt *ItemStMem) Populate(items []*Item) ([]*Item, error) {
	iSt.Items = append(iSt.Items, items...)
	return iSt.Items, nil
}

func (iSt *ItemStMem) GetAll(limit uint32, offset uint32) ([]*Item, error) {
	if offset > uint32(len(iSt.Items)) {
		return nil, ErrIncorrectBounds
	}
	if limit >= uint32(len(iSt.Items))-offset {
		return iSt.Items[offset:], nil
	}
	return iSt.Items[offset : offset+limit], nil
}

func (iSt *ItemStMem) GetById(itemID uint32) (*Item, error) {
	for _, item := range iSt.Items {
		if item.ID == itemID {
			return item, nil
		}
	}
	return nil, ErrItemNotExist
}

func (iSt *ItemStMem) GetBySellerId(sellerID uint32) ([]*Item, error) {
	items := make([]*Item, 0, 10)
	for _, item := range iSt.Items {
		if item.SellerID == sellerID {
			items = append(items, item)
		}
	}
	return items, nil
}

func (iSt *ItemStMem) Move(itemID uint32, quantity uint32) (*Item, error) {
	item, err := iSt.GetById(itemID)
	if err != nil {
		return nil, err
	}
	if item.InStock < quantity {
		return nil, ErrNotEnoughStock
	}
	item.InStock -= quantity
	return item, nil
}

func (iSt *ItemStMem) DeleteById(itemID uint32) error {
	for i := range iSt.Items {
		if iSt.Items[i].ID == itemID {
			if i == len(iSt.Items)-1 {
				iSt.Items = iSt.Items[:i]
			} else {
				iSt.Items = append(iSt.Items[:i], iSt.Items[i+1:]...)
			}
			return nil
		}
	}
	return ErrItemNotExist
}
