package messages

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Show(w http.ResponseWriter, r *http.Request) {
	reader := r.Context().Value("messages").(interface {
		ContentOf(id string) (string, error)
		IsComplete(id string) (bool, error)
	})

	message, err := reader.ContentOf(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	complete, err := reader.IsComplete(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rw := render.FromCtx(r.Context())
	rw.Set("id", r.PathValue("id"))
	rw.Set("message", message)
	rw.Set("complete", complete)

	err = rw.RenderClean("messages/show.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
