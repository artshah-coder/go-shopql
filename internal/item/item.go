package item

type Item struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	CatalogID uint32 `json:"catalog_id"`
	SellerID  uint32 `json:"seller_id"`
	InStock   uint32 `json:"in_stock"`
}

// func (item *Item) Id() string {
// 	return strconv.Itoa(int(item.ID))
// }
