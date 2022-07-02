package internal

import "errors"

type Product struct {
	Id    int
	Name  string
	Price float64
}

type Storage interface {
	List() ([]Product, error)
	Sale(int, int) (int, error)
}

func ValidateSaleQty(qty int) error {
	switch {
	case qty < 1:
		return errors.New("qty cant be less than 1")
	}
	return nil
}
