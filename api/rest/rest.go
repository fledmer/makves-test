package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"items-service/model"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"
)

type ItemStorage interface {
	ItemsById(context.Context, []string) ([]map[string]string, error)
}

type Server struct {
	ctx     context.Context
	storage ItemStorage
}

func New(ctx context.Context, storage ItemStorage) *Server {
	return &Server{
		ctx:     ctx,
		storage: storage,
	}
}

func (s *Server) Run(port int) error {
	slog.Info("starting REST api")
	mux := http.NewServeMux()
	mux.HandleFunc("/get-items", s.getItemsHandler)
	return http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(port)), mux)
}

func (s *Server) getItemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	idsQuery := r.URL.Query()["ids"]
	if len(idsQuery) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ids := strings.Split(r.URL.Query()["ids"][0], ",")
	items, err := s.storage.ItemsById(s.ctx, ids)
	if errors.Is(err, model.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(jsonError(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(items)
	return
}

func jsonError(message any) map[string]any {
	m := map[string]any{
		"error": message,
	}
	return m
}
