package internal

import "errors"

type Product struct {
	Id    int
	Name  string
	Qty   int
	Price float64
}

type Storage interface {
	List() ([]Product, error)
	Sale(int, int) (int, error)
}

func TranslateQtyToStr(qty int) string {
	var quantity string
	switch {
	case qty > 0 && qty < 2:
		quantity = "few"
	case qty >= 2 && qty < 4:
		quantity = "several"
	default:
		quantity = "many"
	}
	return quantity
}

func ValidateSaleQty(qty int) error {
	switch {
	case qty < 1:
		return errors.New("qty cant be less than 1")
	}
	return nil
}
