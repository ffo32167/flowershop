package http

import (
	"encoding/json"
	"net/http"

	"github.com/ffo32167/flowershop/internal/storage"
	"go.uber.org/zap"
)

type listHandler struct {
	storage storage.StorageProducts
	log     *zap.Logger
}

func newListHandler(storage storage.StorageProducts, log *zap.Logger) listHandler {
	return listHandler{storage: storage, log: log}
}

func (l listHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	data, err := l.storage.List(req.Context())
	if err != nil {
		l.log.Error("storage handler error:", zap.Error(err))
	}

	err = json.NewEncoder(res).Encode(data)
	if err != nil {
		l.log.Error("rate handler encoder error:", zap.Error(err))
	}
}
