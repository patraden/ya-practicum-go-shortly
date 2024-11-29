package handler

import (
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"

	"github.com/patraden/ya-practicum-go-shortly/internal/app/dto"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/remover"
)

type DeleteHandler struct {
	log     *zerolog.Logger
	remover remover.URLRemover
}

func NewDeleteHandler(remover remover.URLRemover, log *zerolog.Logger) *DeleteHandler {
	return &DeleteHandler{
		remover: remover,
		log:     log,
	}
}

func (h *DeleteHandler) HandleDelUserURLs(w http.ResponseWriter, r *http.Request) {
	var userSlugs dto.UserSlugBatch

	if err := easyjson.UnmarshalFromReader(r.Body, &userSlugs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	err := h.remover.RemoveUserSlugs(r.Context(), userSlugs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set(ContentType, ContentTypeText)
	w.WriteHeader(http.StatusAccepted)
}