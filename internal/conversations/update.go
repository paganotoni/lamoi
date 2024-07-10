package conversations

import (
	"fmt"
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Update(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())
	convos := r.Context().Value("conversations").(interface {
		Find(id string) (*conversation, error)
		UpdateName(id string, name string) error
	})

	conversation, err := convos.Find(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = convos.UpdateName(conversation.ID, r.FormValue("Name"))
	if err != nil {
		fmt.Println("Error updating conversation name", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	conversation.Name = r.FormValue("Name")

	rw.Set("conversation", conversation)
	err = rw.RenderClean("conversations/update.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
