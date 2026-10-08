package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// dist holds the production build of the React client. The build copies
// frontend/dist here before compiling; in a plain checkout it holds only
// .gitkeep, and the server says so instead of failing to compile.
//
//go:embed all:dist
var embedded embed.FS

const notBuiltPage = `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>Client not built</title></head>
<body style="font-family: system-ui, sans-serif; max-width: 40rem; margin: 3rem auto; line-height: 1.5">
<h1>The client has not been built into this server</h1>
<p>The API is running at <code>/api</code>. To use the interface either start the
development client with <code>npm run dev</code> in <code>frontend/</code>, or build
it and rebuild the server as described in the README.</p>
</body></html>`

func Handler() http.Handler {
	dist, err := fs.Sub(embedded, "dist")
	if err != nil {
		panic(err)
	}
	if _, err := fs.Stat(dist, "index.html"); err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(notBuiltPage))
		})
	}

	files := http.FileServer(http.FS(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name != "" && name != "index.html" {
			if _, err := fs.Stat(dist, name); err == nil {
				if strings.HasPrefix(name, "assets/") {
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				}
				files.ServeHTTP(w, r)
				return
			}
		}

		// Any other path is a client-side route such as /incidents/4, so the
		// page shell is returned and React Router takes it from there.
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFileFS(w, r, dist, "index.html")
	})
}
