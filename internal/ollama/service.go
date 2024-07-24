package ollama

import (
	"cmp"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Builds a new service
func NewService() *service {
	return &service{
		url:    cmp.Or(os.Getenv("OLLAMA_URL"), "http://localhost:11434/api"),
		client: http.Client{},
	}
}

type service struct {
	url    string
	client http.Client
}

type Response struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Context  []int  `json:"context"`
}

func (r Response) EncodedContext() string {
	data, _ := json.Marshal(&r.Context)
	return base64.StdEncoding.EncodeToString(data)
}

// Generate the response to a message in a conversation
// return an error if something happened.
func (s service) Generate(id, model, message, context string, updateFn func(content, context string)) error {
	message = strings.TrimSpace(message)
	message = strings.ReplaceAll(message, "\n", "")

	cson := []byte("[]")
	if context != "" {
		var err error
		cson, err = base64.StdEncoding.DecodeString(context) // I expect this to be []int
		if err != nil {
			return err
		}
	}

	payload := fmt.Sprintf(
		`{"model": "%s","prompt": "%v", "context": %s}`,
		model, message, cson,
	)

	resp, err := s.client.Post(s.url+"/generate", "", strings.NewReader(payload))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("ollama service returned status", resp.StatusCode)
		bb := make([]byte, 1024)
		n, _ := resp.Body.Read(bb)
		fmt.Println("Error response:", string(bb[:n]))

		return fmt.Errorf("ollama service returned status %d", resp.StatusCode)
	}

	// Read the response body
	var response string
	var ncontext string
	for {
		buf := make([]byte, 1024*1024)
		_, err := resp.Body.Read(buf)
		buf = []byte(strings.Trim(string(buf), "\x00"))
		if string(buf) == "" {
			break
		}

		partial := struct {
			Model    string `json:"model"`
			Response string `json:"response"`
			Context  []int  `json:"context"`
		}{}

		err = json.Unmarshal(buf, &partial)
		if err != nil {
			log.Fatal("Error unmarshalling JSON: ", err.Error(), string(buf))
			break
		}

		response += partial.Response
		if len(partial.Context) > 0 {
			data, _ := json.Marshal(&partial.Context)
			ncontext = base64.StdEncoding.EncodeToString(data)
		}

		updateFn(response, ncontext)
	}

	// Call the update function with the response
	// and the new context so it can mark the message
	// as completed
	updateFn(response, ncontext)

	return err
}

// IsOnline checks if the service is online
// by making a request to the status endpoint
func (s service) IsOnline() bool {
	resp, err := s.client.Get(s.url + "/ps")
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (s service) Models() ([]string, error) {
	resp, err := s.client.Get(s.url + "/tags")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama service returned status %d", resp.StatusCode)
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type mms struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	mx := mms{}
	err = json.Unmarshal(bb, &mx)
	if err != nil {
		return nil, err
	}

	models := make([]string, len(mx.Models))
	for i, m := range mx.Models {
		models[i] = strings.Split(m.Name, ":")[0]
	}

	return models, nil
}
