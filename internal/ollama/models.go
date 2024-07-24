package ollama

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Models(w http.ResponseWriter, r *http.Request) {
	ollama := r.Context().Value("ollamaService").(interface {
		Models() ([]string, error)
	})

	models, err := ollama.Models()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rw := render.FromCtx(r.Context())
	rw.Set("models", models)

	err = rw.RenderClean("ollama/models.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
