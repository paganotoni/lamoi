package internal

import (
	"cmp"
	"embed"
	"net/http"
	"os"
	"time"

	"lamoi/internal/conversations"
	"lamoi/internal/messages"
	"lamoi/internal/ollama"
	"lamoi/public"

	"github.com/dustin/go-humanize"
	"github.com/leapkit/leapkit/core/assets"
	"github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server"
)

var (
	//go:embed **/*.html **/*.html *.html
	tmpls embed.FS

	// DB is the database connection builder function
	// that will be used by the application based on the driver and
	// connection string.
	DB = db.ConnectionFn(
		cmp.Or(os.Getenv("DATABASE_URL"), "leapkit.db"),
		db.WithDriver("sqlite3"),
	)
)

// Server interface exposes the methods
// needed to start the server in the cmd/app package
type Server interface {
	Addr() string
	Handler() http.Handler
}

func New() Server {
	assetsManager := assets.NewManager(public.Files)
	renderMW := render.Middleware(
		render.TemplateFS(tmpls, "internal"),

		render.WithDefaultLayout("layout.html"),
		render.WithHelpers(render.AllHelpers),
		render.WithHelpers(map[string]any{
			"assetPath": assetsManager.PathFor,
			"timeSince": func(t time.Time) string {
				return humanize.Time(t)
			},
		}),
	)

	// Creating a new server instance with the
	// default host and port values.
	r := server.New(
		server.WithHost(cmp.Or(os.Getenv("HOST"), "0.0.0.0")),
		server.WithPort(cmp.Or(os.Getenv("PORT"), "3000")),
		server.WithSession(
			cmp.Or(os.Getenv("SESSION_SECRET"), "d720c059-9664-4980-8169-1158e167ae57"),
			cmp.Or(os.Getenv("SESSION_NAME"), "leapkit_session"),
		),
	)

	r.Use(renderMW)
	r.Use(server.InCtxMiddleware("ollamaService", ollama.NewService()))
	r.Use(server.InCtxMiddleware("conversations", conversations.NewService(DB)))
	r.Use(server.InCtxMiddleware("messages", messages.NewService(DB)))

	r.HandleFunc("GET /{$}", conversations.New)
	r.HandleFunc("POST /conversations/{$}", conversations.Send)
	r.HandleFunc("GET /conversations/{id}", conversations.Show)
	r.HandleFunc("GET /conversations/{$}", conversations.List)
	r.HandleFunc("GET /messages/{id}", messages.Show)

	r.HandleFunc("GET /ollama/status", ollama.Status)
	// Mounting the assets manager at the end of the routes
	// so that it can serve the public assets.
	r.HandleFunc(assetsManager.HandlerPattern(), assetsManager.HandlerFn)

	return r
}
