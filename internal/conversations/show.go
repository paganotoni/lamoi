package conversations

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Show(w http.ResponseWriter, r *http.Request) {
	finder := r.Context().Value("conversations").(interface {
		Find(id string) (*conversation, error)
	})

	id := r.PathValue("id")
	conversation, err := finder.Find(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if conversation == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	rw := render.FromCtx(r.Context())
	rw.Set("endpoint", "/conversations/"+id)
	rw.Set("htmx", r.Header.Get("HX-Request") == "true")
	rw.Set("conversation", conversation)

	renderFn := rw.Render
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Add("HX-Push", "/conversations/"+id)
		renderFn = rw.RenderClean
	}

	err = renderFn("conversations/base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
