package bmber

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

var ltRegexp = regexp.MustCompile(`name="lt" value="(.+?)"`)

func (b *Bmber) StepLogin() error {
	// _eventId: submit
	// execution: e1s1
	// lt: xxxxxxxxxxxxxx-passport.bitmain.com
	// password: pppppppp
	// username: uuuuuuuu

	body, err := b.GetPageBytes(b.config.Login.Page)
	if err != nil {
		b.logger.Error("Login page err", zap.Error(err))
		return err
	}

	v := make(url.Values)
	if m := ltRegexp.FindSubmatch(body); len(m) > 0 {
		v.Set("lt", string(m[1]))
	} else {
		b.logger.Error("lt not parsed from login page")
		return errors.New("lt not parsed from login page")
	}
	v.Set("_eventId", "submit")
	v.Set("execution", "e1s1")
	v.Set("username", b.config.Login.Username)
	v.Set("password", b.config.Login.Password)

	req, err := http.NewRequest("POST", b.config.Login.Post, strings.NewReader(v.Encode()))
	if err != nil {
		b.logger.Error("Login POST Request err", zap.Error(err))
		return err
	}

	u, err := url.Parse(b.config.Login.Page)
	if err != nil {
		b.logger.Error("Login.Page to URL err", zap.Error(err))
		return err
	}
	req.Header.Set("Origin", fmt.Sprintf("%s://%s", u.Scheme, u.Host))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", b.config.Login.Page)
	if _, err = b.client.Do(req); err != nil {
		b.logger.Error("Login POST err", zap.Error(err))
		return err
	}
	return nil
}
