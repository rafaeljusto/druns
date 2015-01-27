package session

import (
	"net"
	"net/http"
	"time"

	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/model"
)

func NewSession(sqler db.SQLer, user model.User, ipAddress net.IP) (*http.Cookie, error) {
	session := model.NewSession(user, ipAddress)
	sessionDAO := dao.NewSession(sqler)

	if err := sessionDAO.Save(session); err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:   "session",
		Value:  session.Fingerprint(),
		Path:   "/",
		Secure: true,
	}, nil
}

func CheckSession(sqler db.SQLer, cookie *http.Cookie) (bool, error) {
	sessionId := model.SessionFingerprintId(cookie.Value)
	sessionDAO := dao.NewSession(sqler)

	session, err := sessionDAO.FindById(sessionId)
	if err != nil {
		return false, err
	}

	secret := "" // TODO!

	if !session.CheckFingerprint(cookie.Value, secret) {
		return false, nil
	}

	// TODO: Check session timeout!

	session.LastAccessAt = time.Now()
	if err := sessionDAO.Save(session); err != nil {
		return false, err
	}

	return true, nil
}
