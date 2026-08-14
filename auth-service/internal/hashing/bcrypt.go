package hashing

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHash struct {
	cost int
}

func NewBcryptHash(cost int) *BcryptHash {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &BcryptHash{
		cost: cost,
	}
}

func (h *BcryptHash) Hash(password string) (string, error) {
	if password == "" {
		return "", ErrPasswordMismatch
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		if errors.Is(err, ErrPasswordTooLong) {
			return "", ErrPasswordTooLong
		}
		return "", err
	}
	return string(hash), nil
}
func (h *BcryptHash) CompareTo(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrPasswordMismatch
	}
	return fmt.Errorf("compare password hash: %v", err)
}
