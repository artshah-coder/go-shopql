package graphql

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/artshah-coder/go-shopql/internal/cart"
	"github.com/artshah-coder/go-shopql/internal/catalog"
	"github.com/artshah-coder/go-shopql/internal/item"
	"github.com/artshah-coder/go-shopql/internal/seller"
	"github.com/artshah-coder/go-shopql/internal/session"
	"github.com/artshah-coder/go-shopql/internal/user"
)

var ErrInvalidQuantity = fmt.Errorf("invalid quantity: must be above 0")

type Resolver struct {
	USt    user.UsersStorage
	CSt    catalog.CatalogsStorage
	ISt    item.ItemsStorage
	SellSt seller.SellersStorage
}

// ID is the resolver for the id field.
func (r *catalogResolver) ID(ctx context.Context, obj *catalog.Catalog) (int, error) {
	return int(obj.ID), nil
}

// Parent is the resolver for the parent field.
func (r *catalogResolver) Parent(ctx context.Context, obj *catalog.Catalog) (*catalog.Catalog, error) {
	return r.CSt.GetById(obj.ID)
}

// Childs is the resolver for the childs field.
func (r *catalogResolver) Childs(ctx context.Context, obj *catalog.Catalog) ([]*catalog.Catalog, error) {
	return r.CSt.GetChilds(obj.ID)
}

// Items is the resolver for the items field.
func (r *catalogResolver) Items(ctx context.Context, obj *catalog.Catalog, limit int, offset int) ([]*item.Item, error) {
	if limit < 0 || offset < 0 {
		return nil, item.ErrIncorrectBounds
	}
	return obj.Items.GetAll(uint32(limit), uint32(offset))
}

// Quantity is the resolver for the quantity field.
func (r *inCartItemResolver) Quantity(ctx context.Context, obj *cart.InCartItem) (int, error) {
	return int(obj.Quantity), nil
}

// ID is the resolver for the id field.
func (r *itemResolver) ID(ctx context.Context, obj *item.Item) (int, error) {
	return int(obj.ID), nil
}

// Parent is the resolver for the parent field.
func (r *itemResolver) Parent(ctx context.Context, obj *item.Item) (*catalog.Catalog, error) {
	return r.CSt.GetById(obj.CatalogID)
}

// Seller is the resolver for the seller field.
func (r *itemResolver) Seller(ctx context.Context, obj *item.Item) (*seller.Seller, error) {
	return r.SellSt.GetById(obj.SellerID)
}

// InCart is the resolver for the inCart field.
func (r *itemResolver) InCart(ctx context.Context, obj *item.Item) (int, error) {
	uid, ok := session.UIDFromContext(ctx)
	if !ok {
		return 0, ErrNotAuth
	}
	user, err := r.USt.GetById(uid)
	if err != nil {
		return 0, err
	}
	item, err := user.Cart.GetByItemId(obj.ID)
	if err != nil && !errors.Is(err, cart.ErrNoInCartItem) {
		return 0, err
	}
	if errors.Is(err, cart.ErrNoInCartItem) {
		return 0, nil
	}

	return int(item.Quantity), nil
}

// InStockText is the resolver for the inStockText field.
func (r *itemResolver) InStockText(ctx context.Context, obj *item.Item) (string, error) {
	switch {
	case obj.InStock >= 2 && obj.InStock < 4:
		return "хватает", nil
	case obj.InStock >= 4:
		return "много", nil
	default:
		return "мало", nil
	}
}

// AddToCart is the resolver for the AddToCart field.
func (r *mutationResolver) AddToCart(ctx context.Context, in Bundle) ([]*cart.InCartItem, error) {
	if in.ItemID < 1 {
		return nil, item.ErrItemNotExist
	}
	if in.Quantity < 1 {
		return nil, ErrInvalidQuantity
	}

	uid, ok := session.UIDFromContext(ctx)
	if !ok {
		return nil, ErrNotAuth
	}

	item, err := r.ISt.Move(uint32(in.ItemID), uint32(in.Quantity))
	if err != nil {
		return nil, err
	}

	inCartItm := &cart.InCartItem{
		Item:     item,
		Quantity: uint32(in.Quantity),
	}

	return r.USt.AddToCart(uid, inCartItm)
}

// RemoveFromCart is the resolver for the RemoveFromCart field.
func (r *mutationResolver) RemoveFromCart(ctx context.Context, in Bundle) ([]*cart.InCartItem, error) {
	if in.ItemID < 1 {
		return nil, item.ErrItemNotExist
	}
	if in.Quantity < 1 {
		return nil, ErrInvalidQuantity
	}

	uid, ok := session.UIDFromContext(ctx)
	if !ok {
		return nil, ErrNotAuth
	}

	return r.USt.RemoveFromCart(uid, uint32(in.ItemID), uint32(in.Quantity))
}

// Catalog is the resolver for the Catalog field.
func (r *queryResolver) Catalog(ctx context.Context, id string) (*catalog.Catalog, error) {
	catalogID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if catalogID < 1 {
		return nil, catalog.ErrCatalogNotExist
	}

	return r.CSt.GetById(uint32(catalogID))
}

// Seller is the resolver for the Seller field.
func (r *queryResolver) Seller(ctx context.Context, id string) (*seller.Seller, error) {
	sellerID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if sellerID < 1 {
		return nil, seller.ErrSellerNotExist
	}

	return r.SellSt.GetById(uint32(sellerID))
}

// MyCart is the resolver for the MyCart field.
func (r *queryResolver) MyCart(ctx context.Context) ([]*cart.InCartItem, error) {
	uid, ok := session.UIDFromContext(ctx)
	if !ok {
		return nil, ErrNotAuth
	}
	user, err := r.USt.GetById(uid)
	if err != nil {
		return nil, err
	}
	c, err := user.Cart.GetAll()
	if err != nil {
		return nil, err
	}
	result, ok := c.(*cart.CartInMem)
	if !ok {
		return nil, cart.ErrNoCartInMem
	}
	return *result, nil
}

// ID is the resolver for the id field.
func (r *sellerResolver) ID(ctx context.Context, obj *seller.Seller) (int, error) {
	return int(obj.ID), nil
}

// Items is the resolver for the items field.
func (r *sellerResolver) Items(ctx context.Context, obj *seller.Seller, limit int) ([]*item.Item, error) {
	if limit < 0 {
		return nil, item.ErrIncorrectBounds
	}

	return obj.Items.GetAll(uint32(limit), 0)
}

// Catalog returns CatalogResolver implementation.
func (r *Resolver) Catalog() CatalogResolver { return &catalogResolver{r} }

// InCartItem returns InCartItemResolver implementation.
func (r *Resolver) InCartItem() InCartItemResolver { return &inCartItemResolver{r} }

// Item returns ItemResolver implementation.
func (r *Resolver) Item() ItemResolver { return &itemResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Seller returns SellerResolver implementation.
func (r *Resolver) Seller() SellerResolver { return &sellerResolver{r} }

type catalogResolver struct{ *Resolver }
type inCartItemResolver struct{ *Resolver }
type itemResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type sellerResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct {
	USt    user.UserStMem
	SessSt session.SessionStMem
	CSt    catalog.CatalogStMem
	SellSt seller.SellerStMem
}
*/
