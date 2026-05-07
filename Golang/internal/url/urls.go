package url

import (
	"encoding/json"
	"fmt"
	"golang/db"
	model "golang/db/model"
	"golang/internal/util"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

type UrlCache struct {
	Rw  sync.RWMutex
	Url map[string]Info
}

type Info struct {
	Original string
	Accesses int
}

func InitUrlsCache(database *db.Database) (*UrlCache, error) {
	trace := util.CreateErrorContext("urls.InitUrlsCache")

	urls, urlsErr := database.GetUrls()
	if urlsErr != nil {
		return nil, trace.Apply(urlsErr)
	}

	urlCache := make(map[string]Info)

	for _, u := range urls {

		urlCache[u.Id] = Info{
			Original: u.Original,
			Accesses: u.Accesses,
		}
	}
	return &UrlCache{Url: urlCache}, nil
}

func GetUrls(urls *UrlCache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trace := util.CreateErrorContext("urls.GetUrls")

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(urls.Url); err != nil {
			util.SendHttpError(w, 500, trace.Apply(err))
			return
		}
	}
}

func CreateUrl(urls *UrlCache, db *db.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trace := util.CreateErrorContext("urls.CreateUrl")

		type urlPayload struct {
			Url string `json:"url"`
		}

		body, bodyErr := util.JsonDecodeFromReader[urlPayload](r.Body)
		if bodyErr != nil {
			util.SendHttpError(w, 400, trace.Apply(bodyErr))
			return
		}

		if _, err := url.ParseRequestURI(body.Url); err != nil {
			util.SendHttpError(w, 400, trace.Apply(err))
			return
		}

		existingUrl, _ := db.GetUrlByUrl(body.Url)

		if existingUrl.Id != "" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"id": existingUrl.Id, "url": body.Url})
			return
		}

		urls.Rw.RLock()
		shortID, shortIDErr := util.GenerateShortID(8, func(id string) bool {
			_, exists := urls.Url[id]
			return exists
		})
		urls.Rw.RUnlock()

		if shortIDErr != nil {
			util.SendHttpError(w, 500, trace.Apply(shortIDErr))
			return
		}

		newUrl := model.Url{
			Id:       shortID,
			Original: body.Url,
			Accesses: 0,
		}

		if err := db.SaveUrl(newUrl); err != nil {
			util.SendHttpError(w, 500, trace.Apply(err))
			return
		}

		urls.Rw.Lock()
		urls.Url[shortID] = Info{
			Original: body.Url,
			Accesses: 0,
		}
		urls.Rw.Unlock()

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": shortID, "url": body.Url})
	}
}

func GetUrl(urls *UrlCache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trace := util.CreateErrorContext("urls.GetUrl")

		id := mux.Vars(r)["id"]

		urls.Rw.RLock()
		if info, exists := urls.Url[id]; !exists {
			util.SendHttpError(w, 400, trace.Apply(fmt.Errorf("")))
			return
		} else {
			if err := json.NewEncoder(w).Encode(info); err != nil {
				util.SendHttpError(w, 500, trace.Apply(err))
				urls.Rw.RUnlock()
				return
			}
		}
		urls.Rw.RUnlock()
	}
}

func DeleteUrl(urls *UrlCache, db *db.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trace := util.CreateErrorContext(fmt.Sprintf("urls.DeleteUrl"))

		id := mux.Vars(r)["id"]

		urls.Rw.RLock()
		if _, exists := urls.Url[id]; !exists {
			util.SendHttpError(w, 400, trace.Apply(fmt.Errorf("")))
			return
		}
		urls.Rw.RUnlock()

		err := db.DeleteUrl(id)
		if err != nil {
			util.SendHttpError(w, 500, trace.Apply(err))
			return
		}

		urls.Rw.Lock()
		delete(urls.Url, id)
		urls.Rw.Unlock()
	}
}

func GetOriginalUrl(urls *UrlCache, db *db.Database) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		trace := util.CreateErrorContext("urls.GetOriginalUrl")

		id := mux.Vars(r)["id"]

		urls.Rw.RLock()
		info, exists := urls.Url[id]
		urls.Rw.RUnlock()

		if !exists {
			util.SendHttpError(w, 404, trace.Apply(fmt.Errorf("URL not found")))
			return
		}

		info.Accesses++

		urls.Rw.Lock()
		urls.Url[id] = info
		urls.Rw.Unlock()

		if err := db.UpdateAccesses(id, info.Accesses); err != nil {
			util.SendHttpError(w, 500, trace.Apply(err))
			return
		}

		http.Redirect(w, r, info.Original, http.StatusFound)
	}
}
