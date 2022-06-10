package internal

type Product struct {
	Id   int
	Name string
}

type Storage interface {
	List() ([]Product, error)
}
