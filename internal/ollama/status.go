package ollama

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Status(w http.ResponseWriter, r *http.Request) {
	statuser := r.Context().Value("ollamaService").(interface {
		IsOnline() bool
	})

	rw := render.FromCtx(r.Context())
	rw.Set("online", statuser.IsOnline())

	err := rw.RenderClean("ollama/status.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
