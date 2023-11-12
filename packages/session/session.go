package session

import (
    "fmt"
	"time"

	"github.com/hedenface/catbuttbonanza/packages/cbbhttp"
)

// TODO: delete session data when a user is deleted
// TODO: mfa for login
type Session struct {
	ID                     string    `json:"id"`
	Username               string    `json:"username"`
	Authenticated          bool      `json:"authenticated"`
	LoggedIn               time.Time `json:"logged-in"`
	LastActivity           time.Time `json:"last-activity"`
	ReqRemoteAddr          string    `json:"req-remote-addr"`
	ReqHeaderXForwardedFor string    `json:"req-header-x-forwarded-for"`

	HTTPMessage string `json:"message"`
	HTTPError   string `json:"error"`
}

func Get(id string) Session {
	var s Session

	s.ID = id

	err := cbbhttp.APICall("localhost", 8081, "POST", "session/get", s, &s)
    if err != nil {
        // TODO: logging
        fmt.Printf("session package: Get(): err = %v\n", err)
        return Session{}
    }

    return s
}

func Set(s Session) error {
    var r Session
    err := cbbhttp.APICall("localhost", 8081, "POST", "session/set", s, &r)
    if err != nil {
        return err
    }

    if r.HTTPMessage != "" {
        return nil
    }

    return fmt.Errorf("%s", r.HTTPError)
}

func Delete(id string) error {
    var s, r Session
    s.ID = id
    err := cbbhttp.APICall("localhost", 8081, "POST", "session/delete", s, &r)
    if err != nil {
        return err
    }

    if r.HTTPMessage != "" {
        return nil
    }

    return fmt.Errorf("%s", r.HTTPError)
}
