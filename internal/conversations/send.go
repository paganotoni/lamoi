package conversations

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/leapkit/leapkit/core/render"
)

func Send(w http.ResponseWriter, r *http.Request) {

	convos := r.Context().Value("conversations").(interface {
		Create(message string) (string, error)
		ContextFor(id string) (string, error)
		Update(messageID, content, context string) error
	})

	messages := r.Context().Value("messages").(interface {
		AppendTo(id, message, kind, context string) (string, error)
		AppendPendingTo(string) (string, error)
	})

	conversationID := r.FormValue("ID")
	message := strings.TrimSpace(r.FormValue("message"))
	isNew := conversationID == ""
	if isNew {
		var err error
		conversationID, err = convos.Create(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// If conversation already exists then we append the user message
		// to it
		_, err := messages.AppendTo(conversationID, message, "user", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	context, err := convos.ContextFor(conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// response.Response
	messageID, err := messages.AppendPendingTo(conversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updater := func(content, context string) {
		err := convos.Update(messageID, content, context)
		if err != nil {
			fmt.Println("Error updating message", err)
		}
	}

	ollama := r.Context().Value("ollamaService").(interface {
		Generate(id, message, context string, updater func(content, context string)) error
	})

	go func() {
		err := ollama.Generate(messageID, message, context, updater)
		if err != nil {
			fmt.Println("Error generating response:", err)
		}
	}()

	rw := render.FromCtx(r.Context())
	rw.Set("question", struct {
		Message string
		Kind    string
	}{
		Message: message,
		Kind:    "user",
	})

	rw.Set("title", message)
	rw.Set("pendingID", messageID)
	rw.Set("conversationID", conversationID)
	rw.Set("isNew", isNew)

	w.Header().Set("HX-Push", "/conversations/"+conversationID)

	err = rw.RenderClean("conversations/send.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
