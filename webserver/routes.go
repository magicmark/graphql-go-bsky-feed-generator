package webserver

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-bsky-feed/config"
	"go-bsky-feed/database"
	"log"
	"strconv"
)

func ServeHTTP(db *sql.DB) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("bsky feed generator powered by go-bsky-feed"))
	})

	r.Get("/.well-known/did.json", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"@context": []string{"https://www.w3.org/ns/did/v1"},
			"id":       "did:web:" + config.HOSTNAME,
			"service": []map[string]interface{}{
				{
					"id":              "#bsky_fg",
					"type":            "BskyFeedGenerator",
					"serviceEndpoint": "https://" + config.HOSTNAME,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.Get("/xrpc/app.bsky.feed.describeFeedGenerator", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"encoding": "application/json",
			"body": map[string]interface{}{
				"did": "did:web:" + config.HOSTNAME,
				"feeds": []map[string]interface{}{
					{"uri": config.FEED_URI},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	r.Get("/xrpc/app.bsky.feed.getFeedSkeleton", func(w http.ResponseWriter, r *http.Request) {
		feed := r.URL.Query().Get("feed")
		if feed != config.FEED_URI {
			http.Error(w, "bad 'feed' query parameter", http.StatusBadRequest)
			return
		}

		cursor := r.URL.Query().Get("cursor")
		limit := r.URL.Query().Get("feed")
		if limit == "" {
			limit = "20"
		}
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			log.Println(err)
			http.Error(w, "bad 'limit' query parameter", http.StatusBadRequest)
			return
		}

		feedResults, err := database.GetFeedResults(db, cursor, limitInt)
		if err != nil {
			log.Println(err)
			http.Error(w, "error fetching feed results", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feedResults)

	})

	http.ListenAndServe("0.0.0.0:8081", r)
}
