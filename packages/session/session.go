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

const (
    MaxSessionAge = 60 * 60 * 4
)

func Create(s Session) (string, error) {
    var r Session
    err := cbbhttp.APICall("localhost", 8081, "PUT", "session/create", s, &r)
    if err != nil {
        return "", err
    }

    if r.HTTPMessage == "" {
        return r.ID, nil
    }

    return "", fmt.Errorf("%s", r.HTTPError)
}

func Read(id string) (Session, error) {
	var s Session

	err := cbbhttp.APICall("localhost", 8081, "GET", fmt.Sprintf("session/read/%s", id), "", &s)
    if err != nil {
        return Session{}, err
    }

    return s, nil
}

func Update(s Session) (Session, error) {
    var r Session
    err := cbbhttp.APICall("localhost", 8081, "PATCH", "session/update", s, &r)
    if err != nil {
        return Session{}, err
    }

    if r.HTTPMessage == "" && r.HTTPError == "" {
        return r, nil
    }

    return Session{}, fmt.Errorf("%s", r.HTTPError)
}

func Delete(id string) error {
    var s, r Session
    s.ID = id
    err := cbbhttp.APICall("localhost", 8081, "DELETE", fmt.Sprintf("session/delete/%s", s.ID), "", &r)
    if err != nil {
        return err
    }

    if r.HTTPMessage != "" {
        return nil
    }

    return fmt.Errorf("%s", r.HTTPError)
}
