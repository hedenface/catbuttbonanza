package main

import (
	//"bytes"
	"fmt"
	"net/http"
	//"plugin"
	//"text/template"
	"time"

	"github.com/google/uuid"
)

const (
	templateDir = "html-templates/"
	username = "heden"
	password = "abc"
	maxSessionAge = 60 * 60 * 4
	defaultPort = ":80"
)

var (
	sessions = make(map[string]Session)
)

// TODO: delete session data when a user is deleted
// TODO: mfa for login
type Session struct {
	ID string
	Username string
	Authenticated bool
	LoggedIn time.Time
	ReqRemoteAddr string
	ReqHeaderXForwardedFor string
}

type HTMLPageVars struct {
	Title string
	Head string
	Body string
}


func main() {

	http.HandleFunc("/favicon.png", faviconHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(defaultPort, nil)
}


func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "resources/favicon.png")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := initSession(w, r)

	r.ParseForm()

	if r.Method == "POST" {
		if r.Form["username"][0] == username && r.Form["password"][0] == password {
			session, ok := sessions[sessionID]
			if ok != true {
				// TODO: this shouldn't ever happen
				fmt.Println("this shouldn't happen greegjrehgq")
			} else {
				session.Authenticated = true
				session.Username = r.Form["username"][0]
				session.LoggedIn = time.Now()

				sessions[sessionID] = session
			}
			// TODO: add a redirect to the original page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

	w.Write([]byte(
		htmlTemplatePage(HTMLPageVars{
			Title: "CatButtBonanza Login",
			Head: "",
			Body: htmlTemplateFormLogin(nil),
	})))
}

func handler(w http.ResponseWriter, r *http.Request) {
	sessionID := initSession(w, r)
	redirectNotAuthenticated(w, r, sessionID)
}

func redirectNotAuthenticated(w http.ResponseWriter, r *http.Request, sessionID string) {
	if checkIfAuthenticated(sessionID) != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func checkIfAuthenticated(sessionID string) bool {
	session, ok := sessions[sessionID]
	if ok != true {
		fmt.Println("this shouldn't happen rtteyjeytsgf")
		return false
	}

	return session.Authenticated
}

func startSession(w http.ResponseWriter, r *http.Request) string {
	var sessionID string

	// guarantee a unique sessionID
	for {
		sessionID = uuid.New().String()
		_, ok := sessions[sessionID]

		// if ok isn't true, that means our sessionID wasn't in the map; we have a unique sessionID
		if ok != true {
			break
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: sessionID,
		Path: "/",
		MaxAge: maxSessionAge,
	})

	sessions[sessionID] = Session{
  		ID: sessionID,
  		ReqRemoteAddr: r.RemoteAddr,
  		ReqHeaderXForwardedFor: r.Header.Get("X-Forwarded-For"),
	}

	return sessionID
}

func cleanupAllOldSessions() {
	for key, val := range sessions {
		if time.Since(val.LoggedIn).Seconds() > maxSessionAge {
			delete(sessions, key)
		}
	}
}

func initSession(w http.ResponseWriter, r *http.Request) string {
	sessionCookie, err := r.Cookie("session")

	cleanupAllOldSessions()

	if err != nil || sessionCookie == nil {
		// if we have no session cookie
		// we need to initialize the session and set cookies
		return startSession(w, r)
	} else {
		// if we have a session cookie
		// so now we check if we have corresponding session data

		// if we don't have the session data
		// we should clean up the cookies
		val, ok := sessions[sessionCookie.Value]
		if ok != true {
			deleteCookie(w, "session")
			delete(sessions, sessionCookie.Value)
			return startSession(w, r)
		} else {
			// now lets check that the remoteaddress and x-forwarded-for header matches
			// what was used when the session was set up
			if r.RemoteAddr != val.ReqRemoteAddr || r.Header.Get("X-Forwarded-For") != val.ReqHeaderXForwardedFor {
				deleteCookie(w, "session")
				delete(sessions, sessionCookie.Value)
				return startSession(w, r)
			}
		}
	}

	return sessionCookie.Value
}

func deleteCookie(w http.ResponseWriter, name string) {
 	http.SetCookie(w, &http.Cookie{
 		Name: name,
 		Value: "",
 		Path: "/",
 		MaxAge: -1,
	})
}
