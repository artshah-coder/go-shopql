package user

import (
	"fmt"

	"github.com/artshah-coder/go-shopql/internal/cart"
)

var (
	ErrEmailExists    = fmt.Errorf("user with this e-mail already exists")
	ErrUsernameExists = fmt.Errorf("user with this username already exists")
	ErrUserNotExist   = fmt.Errorf("user does not exist")
)

type UsersStorage interface {
	Add(User) (*User, error)
	GetByEmail(string) (*User, error)
	GetById(uint32) (*User, error)
	DeleteByEmail(string) error
	DeleteById(uint32) error
	AddToCart(uint32, *cart.InCartItem) ([]*cart.InCartItem, error)
	RemoveFromCart(uint32, uint32, uint32) ([]*cart.InCartItem, error)
}

type UserStMem struct {
	Users []*User
	Count uint32
}

func NewUserStMem() *UserStMem {
	return &UserStMem{
		Users: make([]*User, 0, 10),
		Count: 0,
	}
}

func (uSt *UserStMem) Add(user User) (*User, error) {
	for _, u := range uSt.Users {
		if u.Email == user.Email {
			return nil, ErrEmailExists
		}
		if u.Username == user.Username {
			return nil, ErrUsernameExists
		}
	}

	user.ID = uSt.Count + 1
	user.Cart = cart.NewCartInMem()
	uSt.Users = append(uSt.Users, &user)
	uSt.Count++

	return &user, nil
}

func (uSt *UserStMem) GetByEmail(email string) (*User, error) {
	for _, user := range uSt.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotExist
}

func (uSt *UserStMem) GetById(uid uint32) (*User, error) {
	for _, user := range uSt.Users {
		if user.ID == uid {
			return user, nil
		}
	}
	return nil, ErrUserNotExist
}

func (uSt *UserStMem) DeleteByEmail(email string) error {
	for i, user := range uSt.Users {
		if user.Email == email {
			if i == len(uSt.Users)-1 {
				uSt.Users = uSt.Users[:i]
			} else {
				uSt.Users = append(uSt.Users[:i], uSt.Users[i+1:]...)
			}
			return nil
		}
	}
	return ErrUserNotExist
}

func (uSt *UserStMem) DeleteById(uid uint32) error {
	if uSt.Count < uid || uid == 0 {
		return ErrUserNotExist
	}
	for i, user := range uSt.Users {
		if user.ID == uid {
			if i == len(uSt.Users)-1 {
				uSt.Users = uSt.Users[:i]
			} else {
				uSt.Users = append(uSt.Users[:i], uSt.Users[i+1:]...)
			}
			return nil
		}
	}
	return ErrUserNotExist
}

func (uSt *UserStMem) AddToCart(userID uint32, item *cart.InCartItem) ([]*cart.InCartItem, error) {
	user, err := uSt.GetById(userID)
	if err != nil {
		return nil, err
	}
	if user.Cart == nil {
		user.Cart = cart.NewCartInMem()
	}
	c, err := user.Cart.Add(item)
	if err != nil {
		return nil, err
	}
	cartInMem, ok := c.(*cart.CartInMem)
	if !ok {
		return nil, cart.ErrNoCartInMem
	}
	return *cartInMem, nil
}

func (uSt *UserStMem) RemoveFromCart(
	userID uint32, itemID uint32, quantity uint32,
) ([]*cart.InCartItem, error) {
	user, err := uSt.GetById(userID)
	if err != nil {
		return nil, err
	}
	inCartMem, ok := user.Cart.(*cart.CartInMem)
	if !ok {
		return nil, cart.ErrNoCartInMem
	}
	for i, item := range *inCartMem {
		if item.Item.ID == itemID {
			if item.Quantity <= quantity {
				item.Item.InStock += item.Quantity
				if i == len(*inCartMem)-1 {
					*inCartMem = (*inCartMem)[:i]
				} else {
					*inCartMem = append((*inCartMem)[:i], (*inCartMem)[i+1:]...)
				}
			} else {
				item.Item.InStock += quantity
				item.Quantity -= quantity
			}
			break
		}
	}
	return *inCartMem, nil
}
