package http

import (
	"encoding/json"
	"net/http"

	"github.com/ffo32167/flowershop/internal/storage"
	"go.uber.org/zap"
)

//	zap.NewNop() - создаёт дамми-логгер, для тестов
type listHandler struct {
	storage storage.StorageProduct
	log     *zap.Logger
}

type listResponse struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func newListHandler(storage storage.StorageProduct, log *zap.Logger) listHandler {
	return listHandler{storage: storage, log: log}
}

func (l listHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	productList, err := l.storage.List(req.Context())
	if err != nil {
		l.log.Error("listHandler storage error:", zap.Error(err))
	}
	listResponse := make([]listResponse, len(productList))
	for i := range productList {
		listResponse[i].Id = productList[i].Id
		listResponse[i].Name = productList[i].Name
		listResponse[i].Price = productList[i].Price
	}

	err = json.NewEncoder(res).Encode(listResponse)
	if err != nil {
		l.log.Error("listHandler encoder error:", zap.Error(err))
	}
}
