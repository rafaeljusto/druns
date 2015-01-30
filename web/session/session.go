package session

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/core/password"
	"github.com/rafaeljusto/druns/web/config"
)

func NewSession(sqler dao.SQLer, email string, ipAddress net.IP) (*http.Cookie, error) {
	userDAO := dao.NewUser(sqler, nil, "")
	user, err := userDAO.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	session := model.NewSession(user, ipAddress)
	sessionDAO := dao.NewSession(sqler)

	if err := sessionDAO.Save(&session); err != nil {
		return nil, err
	}

	secret, err := password.Decrypt(config.DrunsConfig.Session.Secret)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:   "session",
		Value:  session.Fingerprint(secret),
		Path:   "/",
		Secure: true,
	}, nil
}

func LoadAndCheckSession(sqler dao.SQLer, cookie *http.Cookie,
	ipAddress net.IP) (model.Session, error) {

	var session model.Session
	var err error

	sessionId, err := model.SessionFingerprintId(cookie.Value)
	if err != nil {
		return session, err
	}

	sessionDAO := dao.NewSession(sqler)
	session, err = sessionDAO.FindById(sessionId)
	if err != nil {
		return session, err
	}

	if !session.IPAddress.Equal(ipAddress) {
		return session, core.NewError(fmt.Errorf("IP address '%s' does not match with session IP '%s'",
			ipAddress, session.IPAddress))
	}

	secret, err := password.Decrypt(config.DrunsConfig.Session.Secret)
	if err != nil {
		return session, err
	}

	if !session.CheckFingerprint(cookie.Value, secret) {
		return session, core.NewError(fmt.Errorf("Fingerprint does not match"))
	}

	if time.Now().Sub(session.LastAccessAt) > config.DrunsConfig.Session.Timeout.Duration {
		return session, core.NewError(fmt.Errorf("Session expired"))
	}

	session.LastAccessAt = time.Now()
	if err := sessionDAO.Save(&session); err != nil {
		return session, err
	}

	return session, nil
}
