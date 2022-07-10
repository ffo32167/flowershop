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

type saleResponse struct {
	Message string `json:"message"`
}

func newSaleHandler(storage storage.StorageProducts, log *zap.Logger) saleHandler {
	return saleHandler{storage: storage, log: log}
}

func (s saleHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		s.log.Error("saleHandler: cant convert id:", zap.Error(err))
	}
	cnt, err := strconv.Atoi(mux.Vars(req)["cnt"])
	if err != nil {
		s.log.Error("saleHandler: cant convert cnt:", zap.Error(err))
	}

	err = s.storage.Sale(req.Context(), id, cnt)
	if err != nil {
		s.log.Error("saleHandler storage error:", zap.Error(err))
		saleResponse := saleResponse{Message: err.Error()}
		json.NewEncoder(res).Encode(saleResponse)
		return
	}
	saleResponse := saleResponse{Message: "successful sale!"}
	err = json.NewEncoder(res).Encode(saleResponse)
	if err != nil {
		s.log.Error("saleHandler encoder error:", zap.Error(err))
	}
}
