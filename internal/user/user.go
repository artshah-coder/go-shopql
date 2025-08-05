package user

import "github.com/artshah-coder/go-shopql/internal/cart"

type User struct {
	ID       uint32 `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	password []byte `json:"-"`
	Cart     cart.Cart
}
