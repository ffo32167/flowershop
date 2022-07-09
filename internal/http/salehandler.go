package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ffo32167/flowershop/internal/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type saleHandler struct {
	storage storage.StorageProducts
	log     *zap.Logger
}

type SaleResponse struct {
	Message string
}

func newSaleHandler(storage storage.StorageProducts, log *zap.Logger) saleHandler {
	return saleHandler{storage: storage, log: log}
}

func (s saleHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		s.log.Error("cant convert id:", zap.Error(err))
	}
	cnt, err := strconv.Atoi(mux.Vars(req)["cnt"])
	if err != nil {
		s.log.Error("cant convert cnt:", zap.Error(err))
	}

	_, err = s.storage.Sale(req.Context(), id, cnt)
	if err != nil {
		s.log.Error("storage handler error:", zap.Error(err))
		saleResponse := SaleResponse{Message: err.Error()}
		json.NewEncoder(res).Encode(saleResponse)
		return
	}
	saleResponse := SaleResponse{Message: "successful sale!"}
	err = json.NewEncoder(res).Encode(saleResponse)
	if err != nil {
		s.log.Error("rate handler encoder error:", zap.Error(err))
	}
}
