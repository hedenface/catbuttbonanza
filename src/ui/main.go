package main

import (
	//"bytes"
	//"fmt"
	"net/http"
	//"plugin"
	//"text/template"
	"time"

	//"github.com/google/uuid"
	"github.com/hedenface/catbuttbonanza/packages/log"
	"github.com/hedenface/catbuttbonanza/packages/session"
)

const (
	templateDir = "_html-templates/"
	username = "heden"
	password = "abc"
	defaultPort = ":8080"
)

func main() {
	log.PushStack("main")
	defer log.PopStack()

	log.Setup("ui", log.DebugLevel)

	http.HandleFunc("/favicon.png", faviconHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/", handler)
	log.Errorln(http.ListenAndServe(defaultPort, nil))
}


func faviconHandler(w http.ResponseWriter, r *http.Request) {
	log.PushStack("main")
	defer log.PopStack()

	http.ServeFile(w, r, "resources/favicon.png")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.PushStack("loginHandler")
	defer log.PopStack()

	var s session.Session
	s.ID = initSession(w, r)

	if checkIfAuthenticated(s.ID) == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	r.ParseForm()

	if r.Method == "POST" {
		if r.Form["username"][0] == username && r.Form["password"][0] == password {

			if s.ID == "" {
				// this really shouldn't happen
				log.Warn("No session ID on login Username: %s\n", r.Form["username"][0])
			} else {
				s.Authenticated = true
				s.Username = r.Form["username"][0]
				s.LoggedIn = time.Now()

				_, err := session.Update(s)
				if err != nil {
					log.Error("session.Update() failed: %v\n", err)
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	log.PushStack("logoutHandler")
	defer log.PopStack()

	sessionID := initSession(w, r)
	redirectNotAuthenticated(w, r, sessionID)

	err := session.Delete(sessionID)
	if err != nil {
		log.Error("session.Update() failed: %v\n", err)
	}

	deleteCookie(w, "session")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.PushStack("handler")
	defer log.PopStack()

	sessionID := initSession(w, r)
	redirectNotAuthenticated(w, r, sessionID)
}

func redirectNotAuthenticated(w http.ResponseWriter, r *http.Request, sessionID string) {
	log.PushStack("redirectNotAuthenticated")
	defer log.PopStack()

	if checkIfAuthenticated(sessionID) != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func checkIfAuthenticated(sessionID string) bool {
	log.PushStack("checkIfAuthenticated")
	defer log.PopStack()

	if sessionID == "" {
		return false
	}
	s, err := session.Read(sessionID)
	if err != nil {
		log.Error("error reading sessionID (%s)\n", sessionID)
		return false
	}
	return s.Authenticated
}

func startSession(w http.ResponseWriter, r *http.Request) string {
	log.PushStack("startSession")
	defer log.PopStack()


	sessionID, err := session.Create(session.Session{
  		ReqRemoteAddr: r.RemoteAddr,
  		ReqHeaderXForwardedFor: r.Header.Get("X-Forwarded-For"),
	})

	if err != nil {
		// this will cause some issues initially, but they should be cleaned up
		// in short order
		log.Error("session.Create() failed: %v\n", err)
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
	log.PushStack("initSession")
	defer log.PopStack()

	sessionCookie, err := r.Cookie("session")

	if err != nil || sessionCookie == nil {
		// if we have no session cookie
		// we need to initialize the session and set cookies
		return startSession(w, r)
	}

	// if we have a session cookie
	// so now we check if we have corresponding session data
	s, err := session.Read(sessionCookie.Value)
	if err != nil {
		log.Warn("error reading sessionCooke.Value (%s)\n", sessionCookie.Value)
		return startSession(w, r)
	}

	// if we don't have session data, clean up the cookies
	if s.ID == "" {
		deleteCookie(w, "session")
		err := session.Delete(sessionCookie.Value)
		if err != nil {
			log.Error("session.Delete() failed: %v\n", err)
		}
		return startSession(w, r)
	}

	return sessionCookie.Value
}

func deleteCookie(w http.ResponseWriter, name string) {
	log.PushStack("deleteCookie")
	defer log.PopStack()

 	http.SetCookie(w, &http.Cookie{
 		Name: name,
 		Value: "",
 		Path: "/",
 		MaxAge: -1,
	})
}
