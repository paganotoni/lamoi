package conversations

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/leapkit/leapkit/core/render"
)

var idRe = regexp.MustCompile(`.*\/conversations\/(.*)$`)

func List(w http.ResponseWriter, r *http.Request) {
	var id string
	rouM := idRe.FindStringSubmatch(r.Header.Get("Hx-Current-Url"))
	if len(rouM) >= 2 {
		id = rouM[1]
	}

	regexp.MustCompile(`^/conversations$`)
	lister := r.Context().Value("conversations").(interface {
		List() ([]conversation, error)
	})

	conversations, err := lister.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rw := render.FromCtx(r.Context())
	rw.Set("conversations", conversations)
	rw.Set("id", id)

	err = rw.RenderClean("conversations/list.html")
	if err != nil {
		fmt.Println("Error rendering template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
