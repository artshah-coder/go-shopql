package catalog

import (
	"errors"
	"fmt"

	"github.com/artshah-coder/go-shopql/internal/item"
)

var (
	ErrCatalogExists   = fmt.Errorf("catalog already exists")
	ErrCatalogNotExist = fmt.Errorf("catalog does not exist")
	ErrParentNotExist  = fmt.Errorf("parent catalog does not exist")
	ErrChildAdd        = fmt.Errorf("error while adding child catalog to parent")
	ErrNoChilds        = fmt.Errorf("catalog hasn't childs")
)

type CatalogsStorage interface {
	AddCatalog(Catalog) (*Catalog, error)
	GetAll() ([]*Catalog, error)
	GetChilds(uint32) ([]*Catalog, error)
	GetById(uint32) (*Catalog, error)
	DeleteById(uint32) error
	AddItem(*item.Item) (*item.Item, error)
	GetItemById(uint32) (*item.Item, error)
	GetItemsBySellerId(uint32) ([]*item.Item, error)
}

type CatalogStMem struct {
	Catalogs []*Catalog
}

func NewCatalogStMem() *CatalogStMem {
	return &CatalogStMem{
		Catalogs: make([]*Catalog, 0, 10),
	}
}

func (cSt *CatalogStMem) AddCatalog(catalog Catalog) (*Catalog, error) {
	cat, err := cSt.GetById(catalog.ID)
	if err != nil && !errors.Is(err, ErrCatalogNotExist) {
		return nil, err
	}
	if err == nil && cat != nil {
		return nil, ErrCatalogExists
	}

	if catalog.Childs == nil {
		catalog.Childs = NewCatalogStMem()
	}
	if catalog.Items == nil {
		catalog.Items = item.NewItemStMem()
	}

	if catalog.ParentID != nil {
		parent, err := cSt.GetById(*catalog.ParentID)
		if err != nil {
			return &catalog, ErrChildAdd
		}
		if parent.Childs == nil {
			parent.Childs = NewCatalogStMem()
		}
		parent.Childs.(*CatalogStMem).Catalogs = append(
			parent.Childs.(*CatalogStMem).Catalogs, &catalog,
		)
	} else {
		cSt.Catalogs = append(cSt.Catalogs, &catalog)
	}
	return &catalog, nil
}

func (cSt *CatalogStMem) GetAll() ([]*Catalog, error) {
	return cSt.Catalogs, nil
}

func (cSt *CatalogStMem) GetChilds(parentID uint32) ([]*Catalog, error) {
	parent, err := cSt.GetById(parentID)
	if err != nil {
		return nil, err
	}
	if parent.Childs != nil {
		return parent.Childs.GetAll()
	}
	return nil, ErrNoChilds
}

func (cSt *CatalogStMem) GetById(catalogID uint32) (*Catalog, error) {
	for _, catalog := range cSt.Catalogs {
		if catalog.ID == catalogID {
			return catalog, nil
		}
		if catalog.Childs != nil {
			cat, err := catalog.Childs.GetById(catalogID)
			if err != nil && errors.Is(err, ErrCatalogNotExist) {
				continue
			}
			if err != nil {
				return nil, err
			}
			return cat, nil
		}
	}
	return nil, ErrCatalogNotExist
}

func (cSt *CatalogStMem) DeleteById(catalogID uint32) error {
	for i, catalog := range cSt.Catalogs {
		if catalog.ID == catalogID {
			if i == len(cSt.Catalogs)-1 {
				cSt.Catalogs = cSt.Catalogs[:i]
			} else {
				cSt.Catalogs = append(cSt.Catalogs[:i], cSt.Catalogs[i+1:]...)
			}
			return nil
		}
	}
	return ErrCatalogNotExist
}

func (cSt *CatalogStMem) AddItem(itm *item.Item) (*item.Item, error) {
	catalog, err := cSt.GetById(itm.CatalogID)
	if err != nil {
		return itm, err
	}
	if catalog.Items == nil {
		catalog.Items = item.NewItemStMem()
	}
	return catalog.Items.Add(itm)
}

func (cSt *CatalogStMem) GetItemById(itemID uint32) (*item.Item, error) {
	var item *item.Item
	var err error
	for _, catalog := range cSt.Catalogs {
		if catalog.Items != nil {
			item, err = catalog.Items.GetById(itemID)
			if item != nil {
				return item, err
			}
		}
		if catalog.Childs != nil {
			for _, child := range catalog.Childs.(*CatalogStMem).Catalogs {
				if child.Items != nil {
					item, err = child.Items.GetById(itemID)
					if item != nil {
						return item, err
					}
				}
			}
		}
	}
	return item, err
}

func (cSt *CatalogStMem) GetItemsBySellerId(sellerID uint32) ([]*item.Item, error) {
	items := make([]*item.Item, 0, 10)
	for _, catalog := range cSt.Catalogs {
		if catalog.Items != nil {
			for _, item := range catalog.Items.(*item.ItemStMem).Items {
				if item.SellerID == sellerID {
					items = append(items, item)
				}
			}
		}
		if catalog.Childs != nil {
			childItems, err := catalog.Childs.GetItemsBySellerId(sellerID)
			if err != nil {
				return nil, err
			}
			items = append(items, childItems...)
		}
	}
	return items, nil
}
