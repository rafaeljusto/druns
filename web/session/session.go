package session

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/password"
	"github.com/rafaeljusto/druns/core/session"
	"github.com/rafaeljusto/druns/core/user"
	"github.com/rafaeljusto/druns/web/config"
)

func NewSession(sqler db.SQLer, email string, ipAddress net.IP) (*http.Cookie, error) {
	user, err := user.NewService().FindByEmail(sqler, email)
	if err != nil {
		return nil, err
	}

	s, err := session.NewService().Create(sqler, user, ipAddress)
	if err != nil {
		return nil, err
	}

	secret, err := password.Decrypt(config.DrunsConfig.Session.Secret)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:   "session",
		Value:  s.Fingerprint(secret),
		Path:   "/",
		Secure: true,
	}, nil
}

func LoadAndCheckSession(sqler db.SQLer, cookie *http.Cookie,
	ipAddress net.IP) (session.Session, error) {

	var s session.Session
	var err error

	sessionId, err := session.SessionFingerprintId(cookie.Value)
	if err != nil {
		return s, err
	}

	s, err = session.NewService().FindById(sqler, sessionId)
	if err != nil {
		return s, err
	}

	if !s.IPAddress.Equal(ipAddress) {
		return s, errors.New(fmt.Errorf("IP address '%s' does not match with session IP '%s'",
			ipAddress, s.IPAddress))
	}

	secret, err := password.Decrypt(config.DrunsConfig.Session.Secret)
	if err != nil {
		return s, err
	}

	if !s.CheckFingerprint(cookie.Value, secret) {
		return s, errors.New(fmt.Errorf("Fingerprint does not match"))
	}

	if time.Now().Sub(s.LastAccessAt) > config.DrunsConfig.Session.Timeout.Duration {
		return s, errors.New(fmt.Errorf("Session expired"))
	}

	s.LastAccessAt = time.Now()
	if err = session.NewService().Save(sqler, &s); err != nil {
		return s, err
	}

	return s, nil
}
