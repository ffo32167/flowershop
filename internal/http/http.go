package http

import (
	"net/http"

	"github.com/ffo32167/flowershop/internal/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ApiServer struct {
	storage storage.StorageProducts
	port    string
	log     *zap.Logger
}

func New(storage storage.StorageProducts, port string, log *zap.Logger) ApiServer {
	return ApiServer{storage: storage, port: port, log: log}
}

func (as ApiServer) Run() error {
	saleHandler := newSaleHandler(as.storage, as.log)
	listHandler := newListHandler(as.storage, as.log)

	router := mux.NewRouter()
	router.Handle("/list", listHandler).Methods("GET")
	router.Handle("/sale/{id:[0-9]+}/{cnt:[0-9]+}", saleHandler).Methods("GET")

	err := http.ListenAndServe(as.port, router)
	if err != nil {
		return err
	}
	return nil
}
