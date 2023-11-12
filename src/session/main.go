package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hedenface/catbuttbonanza/packages/session"
	"github.com/hedenface/catbuttbonanza/packages/cbbhttp"
)

const (
	defaultPort = ":8081"
)

var (
	sessions = make(map[string]session.Session)
	lock sync.Mutex
)

func main() {
	http.HandleFunc("/session/set", setHandler)
	http.HandleFunc("/session/get", getHandler)
	http.HandleFunc("/session/delete", deleteHandler)
	fmt.Println(http.ListenAndServe(defaultPort, nil))
}

func cleanupAllOldSessions() {
	fmt.Printf("Sessions: %+v\n", sessions)
	for key, val := range sessions {
		if val.Authenticated == true && time.Since(val.LoggedIn).Seconds() > session.MaxSessionAge {
			fmt.Println("deleting : %d\n", time.Since(val.LoggedIn).Seconds())
			delete(sessions, key)
		}
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var s session.Session

	cleanupAllOldSessions()
	r.ParseForm()

	if r.Method != "POST" {
		cbbhttp.Error(w, http.StatusMethodNotAllowed)
		return
	}

	err := cbbhttp.GetBody(w, r, &s)
	if err != nil {
		fmt.Printf("getHandler: GetBody failed: %v\n", err)
		return
	}

	if s.ID == "" {
		cbbhttp.Error(w, http.StatusBadRequest)
		return
	}

	s, ok := sessions[s.ID]
	if ok == true {
		cbbhttp.ReturnObject(w, s)
	} else {
		cbbhttp.Error(w, http.StatusNotFound)
	}
}


// set is a bit tricky
// since the session data should only ever change
// when someone first visits the ui
// when the user logs in or out
// after a prolonged period of time
func setHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var s session.Session

	cleanupAllOldSessions()
	r.ParseForm()

	fmt.Printf("Sessions: %+v\n", sessions)

	if r.Method != "POST" {
		cbbhttp.Error(w, http.StatusMethodNotAllowed)
		return
	}

	err := cbbhttp.GetBody(w, r, &s)
	if err != nil {
		fmt.Printf("setHandler: GetBody failed: %v\n", err)
		return
	}

	// this means it's a new session
	// so we'll need to make sure it has a unique ID
	if s.ID == "" {
		for {
			sessionID := uuid.New().String()
			_, ok := sessions[sessionID]

			// if ok isn't true, that means our sessionID wasn't in the map; we have a unique sessionID
			if ok != true {
				s.ID = sessionID
				break
			}
		}
	}

	existingSession, ok := sessions[s.ID]

	// this means its a new session
	if ok != true {
		s.LastActivity = time.Now()

		sessions[s.ID] = s
		cbbhttp.ReturnObject(w, s)
		return
	}

	// something fishy is happening here
	if existingSession.ReqRemoteAddr != s.ReqRemoteAddr || existingSession.ReqHeaderXForwardedFor != s.ReqHeaderXForwardedFor {
		delete(sessions, s.ID)
		cbbhttp.Message(w, "session deleted due to suspicious activity")
		return
	}

	// if the user logged out, we just clear the session
	if existingSession.Authenticated == true && s.Authenticated == false {
		delete(sessions, s.ID)
		cbbhttp.Message(w, "session deleted by logging out")
		return
	}

	// if the user just logged in, we update the session with the time
	if existingSession.Authenticated == false && s.Authenticated == true {
		s.LoggedIn = time.Now()
	}

	// finally, we just update the session
	s.LastActivity = time.Now()
	sessions[s.ID] = s
	cbbhttp.Message(w, "existing session updated")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var s session.Session

	cleanupAllOldSessions()
	r.ParseForm()

	err := cbbhttp.GetBody(w, r, &s)
	if err != nil {
		fmt.Printf("getHandler: GetBody failed: %v\n", err)
		return
	}

	if s.ID == "" {
		cbbhttp.Error(w, http.StatusBadRequest)
		return
	}

	s, ok := sessions[s.ID]
	if ok == true {
		delete(sessions, s.ID)
		cbbhttp.Message(w, "session deleted")
	} else {
		cbbhttp.Message(w, "session not found")
	}
}
