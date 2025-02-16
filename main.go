package main

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	_ "golang.org/x/crypto/x509roots/fallback"
)

//go:embed dist
var embedFS embed.FS

func must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}
	return a
}

//go:embed VERSION.txt
var version string

//go:embed TAG.txt
var tag string

func main() {
	version = strings.TrimSpace(version)
	tag = strings.TrimSpace(tag)

	fs := http.FileServer(http.FS(must(fs.Sub(embedFS, "dist"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		fs.ServeHTTP(w, r)
	})
	http.HandleFunc("GET /{$}", route("/", IndexData))
	http.HandleFunc("GET /about/{id...}", route("/about", AboutData))
	http.HandleFunc("GET /health", healthHandler(version, tag))

	dev := os.Getenv("DEV") != ""
	if dev {
		fmt.Println("websocket enabled")
		http.HandleFunc("/ws", wsHandler)
	}

	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	fmt.Println("listening on " + port)
	http.ListenAndServe(":"+port, nil)
}

func IndexData(r *http.Request, data RouteData) {
	data["prompt"] = "Como estas?"
}

func AboutData(r *http.Request, data RouteData) {
	data["greeting"] = "halo!"
	id := r.PathValue("id")
	if id != "" {
		data["id"] = r.PathValue("id")
	}
}

type (
	RouteData map[string]interface{}
	LoadFunc  func(r *http.Request, data RouteData)
)

const HeadTag = "<head>"

func route(path string, load LoadFunc) http.HandlerFunc {
	log.Println("register", path)
	if path == "/" {
		path = ""
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("page", r.URL.Path)
		content, err := embedFS.ReadFile("dist" + path + "/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		i := strings.Index(string(content), HeadTag)
		if i == -1 {
			http.Error(w, "head tag not found", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{}
		load(r, data)
		dataJSON, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		content = append(
			content[:i+len(HeadTag)],
			append([]byte(
				"\n"+
					strings.Repeat(" ", 8)+
					`<script>window.__DATA__ = `+string(dataJSON)+`;</script>`),
				content[i+len(HeadTag):]...,
			)...)
		if os.Getenv("DEV") != "" {
			content = append(content, reloader...)
		}
		w.Write(content)
	}
}

func healthHandler(version, tag string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := struct {
			Version string `json:"version"`
			Tag     string `json:"tag"`
		}{Version: version, Tag: tag}
		err := json.NewEncoder(w).Encode(health)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.NotFound(w, r)
	}
	defer func() {
		cerr := ws.Close()
		if cerr != nil {
			log.Println("close websocket:", cerr)
		}
	}()
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
		case <-r.Context().Done():
			return
		}
	}
}

const reloader = `
<script>
	(function () {
		const { host } = document.location;
		const url = "ws://" + host + "/ws";
		console.log("ws/reloader", url);
		let disconnected = false;
		function monitor() {
			let ws = new WebSocket(url);
			ws.onopen = () => {
				console.log("ws: open");
				if (disconnected) location.reload();
			};
			ws.onclose = () => {
				console.error("ws: close");
				disconnected = true;
				ws = null;
				setTimeout(monitor, 1000);
			};
		}
		monitor();
	}());
</script>
`
