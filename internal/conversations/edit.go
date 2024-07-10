package conversations

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Edit(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	convos := r.Context().Value("conversations").(interface {
		Find(id string) (*conversation, error)
	})

	conversation, err := convos.Find(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rw.Set("conversation", conversation)
	err = rw.RenderClean("conversations/edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
