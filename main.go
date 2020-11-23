package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/Andreeey/basketball/game"
	"github.com/gorilla/websocket"
)

const appAddr = "127.0.0.1:8080"

var (
	upgrader   = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	games      []*game.Game
	gamesMutex sync.Mutex
)

func main() {
	log.Printf("App is starting please go to http://%s/ for UI", appAddr)
	topScorerChan := make(chan *game.Game, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "index.html") })
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { serveWs(w, r, topScorerChan) })
	http.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) { serveGames(w, r, topScorerChan) })
	err := http.ListenAndServe(appAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveWs(w http.ResponseWriter, r *http.Request, topScorerChan chan *game.Game) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws upgrade error", err)
		return
	}

	defer ws.Close() // nolint

	for {
		// FIXME this will work only with single ws instance
		p := <-topScorerChan
		err = ws.WriteJSON(p)
		if err != nil {
			log.Println("ws write error", err)
			break
		}
	}
}

func serveGames(w http.ResponseWriter, r *http.Request, topScorerChan chan *game.Game) {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	if r.Method == http.MethodPost {
		var req []createGamesRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, g := range games {
			g.End()
		}
		games = make([]*game.Game, len(req))
		for i, g := range req {
			g := game.New(g.Name, topScorerChan)
			g.Start()
			games[i] = g
		}

		w.WriteHeader(http.StatusCreated)
	}

	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(games)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type createGamesRequest struct {
	Name string
}
