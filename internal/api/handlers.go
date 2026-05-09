package api

import (
	"encoding/json"
	"net/http"

	"github.com/asim-sde/rag-search-service/internal/db"
	"github.com/asim-sde/rag-search-service/internal/embedding"
	"github.com/asim-sde/rag-search-service/internal/rate"
)

type Handler struct {
	db      *db.DB
	embed   *embedding.Service
	limiter *rate.Limiter
}

func NewHandler(d *db.DB, e *embedding.Service, l *rate.Limiter) *Handler {
	return &Handler{db: d, embed: e, limiter: l}
}

type DocumentRequest struct {
	Content string `json:"content"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

func (h *Handler) IngestDocument(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("X-Client-Id")
	if clientID == "" {
		clientID = "anonymous"
	}

	if !h.limiter.AllowRequest(r.Context(), clientID) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var req DocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	emb := h.embed.GenerateEmbedding(req.Content)
	id, err := h.db.SaveDocument(r.Context(), req.Content, emb)
	if err != nil {
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("X-Client-Id")
	if clientID == "" {
		clientID = "anonymous"
	}

	if !h.limiter.AllowRequest(r.Context(), clientID) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Query == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	emb := h.embed.GenerateEmbedding(req.Query)
	results, err := h.db.SearchNearest(r.Context(), emb, 5)
	if err != nil {
		http.Error(w, "Failed to search", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(results)
}
