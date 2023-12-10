package main

import (
	//"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hedenface/catbuttbonanza/packages/log"
	"github.com/hedenface/catbuttbonanza/packages/session"
	"github.com/hedenface/catbuttbonanza/packages/cbbhttp"
)

const (
	defaultPort = ":8081"
)

var (
	sessions = make(map[string]session.Session)
	mu sync.Mutex
)

func main() {
	log.Setup("session", log.DebugLevel)

	http.HandleFunc("/session/create", createHandler)
	http.HandleFunc("/session/read/", readHandler)
	http.HandleFunc("/session/update", updateHandler)
	http.HandleFunc("/session/delete/", deleteHandler)
	fmt.Println(http.ListenAndServe(defaultPort, nil))
}

func cleanupAllOldSessions() {
	log.Debug("Sessions: %+v\n", sessions)
	for key, val := range sessions {
		if val.Authenticated == true && time.Since(val.LoggedIn).Seconds() > session.MaxSessionAge {
			delete(sessions, key)
		}
	}
}

func getSessionFromPath(w http.ResponseWriter, r *http.Request, path string, method string) (session.Session, error) {
	cleanupAllOldSessions()

	if r.Method != method {
		cbbhttp.Error(w, http.StatusMethodNotAllowed)
		return session.Session{}, errors.New("invalid method")
	}

	id := strings.TrimPrefix(r.URL.Path, path)

	s, ok := sessions[id]
	if ok != true {
		cbbhttp.Error(w, http.StatusNotFound)
		return session.Session{}, errors.New("session not found")
	}

	return s, nil
}

func getSessionFromBody(w http.ResponseWriter, r *http.Request, method string) (session.Session, error) {
	cleanupAllOldSessions()
	var s session.Session

	r.ParseForm()

	if r.Method != method {
		cbbhttp.Error(w, http.StatusMethodNotAllowed)
		return session.Session{}, errors.New("invalid method")
	}

	err := cbbhttp.GetBody(w, r, &s)
	if err != nil {
		cbbhttp.Error(w, http.StatusBadRequest)
		return session.Session{}, err
	}

	return s, nil
}

func getUniqueSessionID() string {
	for {
		sessionID := uuid.New().String()
		_, ok := sessions[sessionID]

		// if ok isn't true, that means our sessionID wasn't in the map; we have a unique sessionID
		if ok != true {
			return sessionID
		}
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	s, err := getSessionFromBody(w, r, "PUT")
	if err != nil {
		log.Error("createHandler: getSessionFromBody failed (%v)\n", err)
		return
	}

	s.ID = getUniqueSessionID()
	s.LastActivity = time.Now()

	sessions[s.ID] = s
	cbbhttp.ReturnObject(w, s)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	s, err := getSessionFromPath(w, r, "/session/read/", "GET")
	if err != nil {
		log.Error("readHandler: getSessionFromPath failed (%v)\n", err)
		return
	}

	cbbhttp.ReturnObject(w, s)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	s, err := getSessionFromBody(w, r, "PATCH")
	if err != nil {
		log.Error("updateHandler: getSessionFromBody failed (%v)\n", err)
		return
	}

	if s.ID == "" {
		cbbhttp.Error(w, http.StatusBadRequest)
		return
	}

	s.LastActivity = time.Now()
	sessions[s.ID] = s

	cbbhttp.Message(w, "session updated")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	s, err := getSessionFromPath(w, r, "/session/delete/", "DELETE")
	if err != nil {
		log.Error("deleteHandler: getSessionFromPath failed (%v)\n", err)
		return
	}

	delete(sessions, s.ID)
	cbbhttp.Message(w, "session deleted")
}
