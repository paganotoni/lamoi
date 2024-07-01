package conversations

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func New(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	rw.Set("endpoint", "/conversations/create")
	rw.Set("conversation", conversation{
		Name: "New Conversation",
	})

	err := rw.Render("conversations/base.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
