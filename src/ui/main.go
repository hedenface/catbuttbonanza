package main

import (
	//"bytes"
	"fmt"
	"net/http"
	//"plugin"
	//"text/template"
	"time"

	//"github.com/google/uuid"
	"github.com/hedenface/catbuttbonanza/packages/log"
	"github.com/hedenface/catbuttbonanza/packages/session"
)

const (
	templateDir = "html-templates/"
	username = "heden"
	password = "abc"
	defaultPort = ":8080"
)

var (
	logger = log.Setup("ui")
)

func main() {
	http.HandleFunc("/favicon.png", faviconHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", handler)
	logger.Println(http.ListenAndServe(defaultPort, nil))
}


func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "resources/favicon.png")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var s session.Session
	s.ID = initSession(w, r)

	r.ParseForm()

	if r.Method == "POST" {
		if r.Form["username"][0] == username && r.Form["password"][0] == password {

			if s.ID == "" {
				// this really shouldn't happen
				logger.Printf("[loginHandler] No session ID on login Username: %s\n", r.Form["username"][0])
			} else {
				s.Authenticated = true
				s.Username = r.Form["username"][0]
				s.LoggedIn = time.Now()

				_, err := session.Set(s)
				if err != nil {
					logger.Printf("[loginHandler] session.Set() failed: %v\n", err)
				}
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
	if sessionID == "" {
		return false
	}
	s := session.Get(sessionID)
	return s.Authenticated
}

func startSession(w http.ResponseWriter, r *http.Request) string {

	sessionID, err := session.Set(session.Session{
  		ReqRemoteAddr: r.RemoteAddr,
  		ReqHeaderXForwardedFor: r.Header.Get("X-Forwarded-For"),
	})

	if err != nil {
		// this will cause some issues initially, but they should be cleaned up
		// in short order
		logger.Printf("[startSession] session.Set() failed: %v\n", err)
		return ""
	}

	http.SetCookie(w, &http.Cookie{
		Name: "session",
		Value: sessionID,
		Path: "/",
		MaxAge: session.MaxSessionAge,
	})

	return sessionID
}

func initSession(w http.ResponseWriter, r *http.Request) string {
	sessionCookie, err := r.Cookie("session")

	if err != nil || sessionCookie == nil {
		// if we have no session cookie
		// we need to initialize the session and set cookies
		return startSession(w, r)
	}

	// if we have a session cookie
	// so now we check if we have corresponding session data
	s := session.Get(sessionCookie.Value)

	// if we don't have session data, clean up the cookies
	if s.ID == "" {
		deleteCookie(w, "session")
		err := session.Delete(sessionCookie.Value)
		if err != nil {
			logger.Printf("[initSession] session.Delete() failed: %v\n", err)
		}
		return startSession(w, r)
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
