package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hedenface/catbuttbonanza/packages/cbbhttp"
)

const (
	maxSessionAge = 60 * 60 * 4
	defaultPort = ":8081"
)

var (
	sessions = make(map[string]Session)
	lock sync.Mutex
)

// TODO: delete session data when a user is deleted
// TODO: mfa for login
type Session struct {
	ID string
	Username string
	Authenticated bool
	LoggedIn time.Time
	LastActivity time.Time
	ReqRemoteAddr string
	ReqHeaderXForwardedFor string
}

func main() {
	http.HandleFunc("/session/set", setHandler)
	http.HandleFunc("/session/get", getHandler)
	http.HandleFunc("/session/delete", deleteHandler)
	fmt.Println(http.ListenAndServe(defaultPort, nil))
}

func cleanupAllOldSessions() {
	fmt.Printf("Sessions: %+v\n", sessions)
	for key, val := range sessions {
		if val.Authenticated == true && time.Since(val.LoggedIn).Seconds() > maxSessionAge {
			fmt.Println("deleting : %d\n", time.Since(val.LoggedIn).Seconds())
			delete(sessions, key)
		}
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var session Session

	cleanupAllOldSessions()
	r.ParseForm()

	if r.Method != "POST" {
		cbbhttp.Error(w, http.StatusMethodNotAllowed)
		return
	}

	err := cbbhttp.GetBody(w, r, &session)
	if err != nil {
		fmt.Printf("getHandler: GetBody failed: %v\n", err)
		return
	}

	if session.ID == "" {
		cbbhttp.Error(w, http.StatusBadRequest)
		return
	}

	session, ok := sessions[session.ID]
	if ok == true {
		cbbhttp.ReturnObject(w, session)
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

	var session Session

	cleanupAllOldSessions()
	r.ParseForm()

	fmt.Printf("Sessions: %+v\n", sessions)

	if r.Method != "POST" {
		cbbhttp.Error(w, http.StatusMethodNotAllowed)
		return
	}

	err := cbbhttp.GetBody(w, r, &session)
	if err != nil {
		fmt.Printf("setHandler: GetBody failed: %v\n", err)
		return
	}

	if session.ID == "" {
		cbbhttp.Error(w, http.StatusBadRequest)
		return
	}

	existingSession, ok := sessions[session.ID]

	// this means its a new session
	if ok != true {
		session.LastActivity = time.Now()

		sessions[session.ID] = session
		cbbhttp.Message(w, "session added")
		return
	}

	// something fishy is happening here
	if existingSession.ReqRemoteAddr != session.ReqRemoteAddr || existingSession.ReqHeaderXForwardedFor != session.ReqHeaderXForwardedFor {
		delete(sessions, session.ID)
		cbbhttp.Message(w, "session deleted due to suspicious activity")
		return
	}

	// if the user logged out, we just clear the session
	if existingSession.Authenticated == true && session.Authenticated == false {
		delete(sessions, session.ID)
		cbbhttp.Message(w, "session deleted by logging out")
		return
	}

	// if the user just logged in, we update the session with the time
	if existingSession.Authenticated == false && session.Authenticated == true {
		session.LoggedIn = time.Now()
	}

	// finally, we just update the session
	session.LastActivity = time.Now()
	sessions[session.ID] = session
	cbbhttp.Message(w, "existing session updated")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	lock.Lock()
	defer lock.Unlock()

	var session Session

	cleanupAllOldSessions()
	r.ParseForm()

	err := cbbhttp.GetBody(w, r, &session)
	if err != nil {
		fmt.Printf("getHandler: GetBody failed: %v\n", err)
		return
	}

	if session.ID == "" {
		cbbhttp.Error(w, http.StatusBadRequest)
		return
	}

	session, ok := sessions[session.ID]
	if ok == true {
		delete(sessions, session.ID)
		cbbhttp.Message(w, "session deleted")
	} else {
		cbbhttp.Message(w, "session not found")
	}
}
