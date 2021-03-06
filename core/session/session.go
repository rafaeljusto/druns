package session

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/user"
)

var (
	SessionFingerprint = regexp.MustCompile("^[0-9]+-[0-9A-Z]+$")
)

type Session struct {
	Id           int
	User         user.User
	IPAddress    net.IP
	CreatedAt    time.Time
	LastAccessAt time.Time
	revision     uint64
}

func NewSession(u user.User, ip net.IP) Session {
	return Session{
		User:         u,
		IPAddress:    ip,
		CreatedAt:    time.Now(),
		LastAccessAt: time.Now(),
	}
}

func (s *Session) Fingerprint(secret string) string {
	id := strconv.Itoa(s.Id)
	data := s.IPAddress.String() + secret
	mac := hmac.New(sha1.New, []byte(id))
	mac.Write([]byte(data))
	hash := strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))
	return fmt.Sprintf("%s-%s", id, hash)
}

func SessionFingerprintId(fingerprint string) (int, error) {
	if !SessionFingerprint.MatchString(fingerprint) {
		return 0, errors.New(fmt.Errorf("Session fingerprint '%s' has an invalid format",
			fingerprint))
	}

	idStr := strings.Split(fingerprint, "-")[0]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return id, errors.New(err)
	}

	return id, nil
}

func (s *Session) CheckFingerprint(fingerprint, secret string) bool {
	return s.Fingerprint(secret) == fingerprint
}
