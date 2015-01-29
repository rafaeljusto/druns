package session

import (
	"net"
	"net/http"
	"time"

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

func CheckSession(sqler dao.SQLer, cookie *http.Cookie, ipAddress net.IP) (bool, error) {
	sessionId, err := model.SessionFingerprintId(cookie.Value)
	if err != nil {
		return false, err
	}

	sessionDAO := dao.NewSession(sqler)
	session, err := sessionDAO.FindById(sessionId)
	if err != nil {
		return false, err
	}

	if !session.IPAddress.Equal(ipAddress) {
		return false, nil
	}

	secret, err := password.Decrypt(config.DrunsConfig.Session.Secret)
	if err != nil {
		return false, err
	}

	if !session.CheckFingerprint(cookie.Value, secret) {
		return false, nil
	}

	if time.Now().Sub(session.LastAccessAt) > config.DrunsConfig.Session.Timeout.Duration {
		return false, nil
	}

	session.LastAccessAt = time.Now()
	if err := sessionDAO.Save(&session); err != nil {
		return false, err
	}

	return true, nil
}
