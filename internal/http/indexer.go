package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/autobrr/autobrr/internal/domain"

	"github.com/go-chi/chi"
)

type indexerService interface {
	Store(indexer domain.Indexer) (*domain.Indexer, error)
	Update(indexer domain.Indexer) (*domain.Indexer, error)
	List() ([]domain.Indexer, error)
	GetAll() ([]*domain.IndexerDefinition, error)
	GetTemplates() ([]domain.IndexerDefinition, error)
	Delete(id int) error
}

type indexerHandler struct {
	encoder encoder
	service indexerService
	ircSvc  ircService
}

func newIndexerHandler(encoder encoder, service indexerService, ircSvc ircService) *indexerHandler {
	return &indexerHandler{
		encoder: encoder,
		service: service,
		ircSvc:  ircSvc,
	}
}

func (h indexerHandler) Routes(r chi.Router) {
	r.Get("/schema", h.getSchema)
	r.Post("/", h.store)
	r.Put("/", h.update)
	r.Get("/", h.getAll)
	r.Get("/options", h.list)
	r.Delete("/{indexerID}", h.delete)
}

func (h indexerHandler) getSchema(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	indexers, err := h.service.GetTemplates()
	if err != nil {
		//
	}

	h.encoder.StatusResponse(ctx, w, indexers, http.StatusOK)
}

func (h indexerHandler) store(w http.ResponseWriter, r *http.Request) {
	var data domain.Indexer

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return
	}

	indexer, err := h.service.Store(data)
	if err != nil {
		//
		h.encoder.StatusResponse(r.Context(), w, nil, http.StatusBadRequest)
		return
	}

	h.encoder.StatusResponse(r.Context(), w, indexer, http.StatusCreated)
}

func (h indexerHandler) update(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		data domain.Indexer
	)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return
	}

	indexer, err := h.service.Update(data)
	if err != nil {
		//
	}

	h.encoder.StatusResponse(ctx, w, indexer, http.StatusOK)
}

func (h indexerHandler) delete(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		idParam = chi.URLParam(r, "indexerID")
	)

	id, _ := strconv.Atoi(idParam)

	if err := h.service.Delete(id); err != nil {
		// return err
	}

	h.encoder.StatusResponse(ctx, w, nil, http.StatusNoContent)
}

func (h indexerHandler) getAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	indexers, err := h.service.GetAll()
	if err != nil {
		//
	}

	h.encoder.StatusResponse(ctx, w, indexers, http.StatusOK)
}

func (h indexerHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	indexers, err := h.service.List()
	if err != nil {
		//
	}

	h.encoder.StatusResponse(ctx, w, indexers, http.StatusOK)
}