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
	http.HandleFunc("GET /ip", ipHandler)

	dev := os.Getenv("DEV") != ""
	if dev {
		fmt.Println("websocket enabled")
		http.HandleFunc("/ws", wsHandler)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("listening on " + port)
	http.ListenAndServe(":"+port, nil)
}

func IndexData(r *http.Request) interface{} {
	return struct {
		Prompt string `json:"prompt"`
	}{Prompt: "Como estas?"}
}

func AboutData(r *http.Request) interface{} {
	data := struct {
		Greeting string `json:"greeting"`
		ID       string `json:"id"`
	}{Greeting: "halo!"}

	id := r.PathValue("id")
	if id != "" {
		data.ID = r.PathValue("id")
	}
	return data
}

type loadFunc func(r *http.Request) interface{}

const headTag = "<head>"

func route(path string, load loadFunc) http.HandlerFunc {
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

		data, err := json.Marshal(load(r))
		if err != nil {
			msg := fmt.Sprintf("error loading page %q data: %v", path, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		injection := headTag + "\n" + strings.Repeat(" ", 6) + `<script>window.__DATA__ = ` + string(data) + `;</script>`
		strings.NewReplacer(headTag, injection).WriteString(w, string(content))

		if os.Getenv("DEV") != "" {
			w.Write(reloader)
		}
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

type IP struct {
	IP string `json:"ip"`
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://api.myip.com"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var ip IP
	err = json.NewDecoder(resp.Body).Decode(&ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(ip.IP))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

var reloader = []byte(`
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
`)
