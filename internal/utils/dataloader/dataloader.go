package dataloader

import (
	"encoding/json"
	"io"
	"os"

	"github.com/artshah-coder/go-shopql/internal/catalog"
	"github.com/artshah-coder/go-shopql/internal/item"
	"github.com/artshah-coder/go-shopql/internal/seller"
)

func MustDataLoad(
	fileName string, catalogs catalog.CatalogsStorage, items item.ItemsStorage,
	sellers seller.SellersStorage,
) {
	in, err := os.Open(fileName)
	if err != nil {
		panic("error while opening file: " + err.Error())
	}
	decoder := json.NewDecoder(in)
	catalog := catalog.Catalog{}
	parentIds := make([]uint32, 0, 1)
	nestedChilds := 0
	isCatalog := false
	isChilds := false
	isItems := false
	isSellers := false
LOOP:
	for {
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			panic("error while reading file")
		}
		if tokenErr == io.EOF {
			break
		}
		switch tok := tok.(type) {
		case string:
			switch tok {
			case "catalog":
				isCatalog = true
				continue
			case "childs":
				if isChilds {
					nestedChilds++
					continue
				}
				isChilds = true
				continue
			case "items":
				isItems = true
				continue
			case "sellers":
				isSellers = true
				continue
			case "id":
				if !decoder.More() {
					panic("invalid json format")
				}
				tok, tokenErr := decoder.Token()
				if tokenErr != nil && tokenErr != io.EOF {
					panic("error while reading file")
				}
				if tokenErr == io.EOF {
					break LOOP
				}
				if len(parentIds) != 0 {
					catalog.ParentID = &parentIds[len(parentIds)-1]
				}
				catalogId := uint32(tok.(float64))
				catalog.ID = catalogId
				parentIds = append(parentIds, catalogId)
				continue
			case "name":
				if !decoder.More() {
					panic("invalid json format")
				}
				tok, tokenErr := decoder.Token()
				if tokenErr != nil && tokenErr != io.EOF {
					panic("error while reading file")
				}
				if tokenErr == io.EOF {
					break
				}
				catalog.Name = tok.(string)
				if isChilds {
					tmp := new(uint32)
					*tmp = *catalog.ParentID
					catalog.ParentID = tmp
				}
				catalogs.AddCatalog(catalog)
				continue
			}
		case json.Delim:
			switch tok.String() {
			case "[":
				if isItems {
					for decoder.More() {
						itm := new(item.Item)
						err := decoder.Decode(itm)
						if err != nil {
							panic("error while reading file")
						}
						if len(parentIds) != 0 {
							itm.CatalogID = parentIds[len(parentIds)-1]
						}
						_, err = items.Add(itm)
						if err != nil {
							panic(err)
						}
						_, err = catalogs.AddItem(itm)
						if err != nil {
							panic(err)
						}
					}
					continue
				}
				if isChilds {
					continue
				}
				if isSellers {
					for decoder.More() {
						seller := new(seller.Seller)
						err := decoder.Decode(seller)
						if err != nil {
							panic("error while reading file")
						}
						items, err := catalogs.GetItemsBySellerId(seller.ID)
						if err != nil {
							panic("error while get items by seller ID from catalog storage")
						}
						if seller.Items == nil {
							seller.Items = item.NewItemStMem()
						}
						seller.Items.Populate(items)
						sellers.Add(*seller)
					}
					continue
				}
			case "{":
				continue
			case "]":
				if isItems {
					isItems = false
					continue
				}
				if isChilds {
					if nestedChilds == 0 {
						isChilds = false
						continue
					}
					nestedChilds--
					continue
				}
				if isSellers {
					isSellers = false
					continue
				}
			case "}":
				if isChilds {
					if len(parentIds) != 0 {
						parentIds = parentIds[:len(parentIds)-1]
					}
					continue
				}
				if isCatalog {
					isCatalog = false
					continue
				}
			}
		}
	}
}
