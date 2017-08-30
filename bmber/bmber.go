package bmber

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/empirefox/bmb/bmbconfig"
	"github.com/empirefox/cement/clog"
	"github.com/empirefox/go-cloudflare-scraper"

	"golang.org/x/net/publicsuffix"
)

type Bmber struct {
	config *bmbconfig.Config
	client *http.Client
	logger clog.Logger

	cartid string
	jar    *cookiejar.Jar
}

func NewBmber(config *bmbconfig.Config) (*Bmber, error) {
	logger, err := clog.NewLogger(config.Clog)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "locale",
		Value:  "zh-CN",
		Path:   "/",
		Domain: ".bitmain.com",
	}
	cookies = append(cookies, cookie)

	u, _ := url.Parse(config.Req.PageHome)
	jar.SetCookies(u, cookies)
	u, _ = url.Parse(config.Login.Page)
	jar.SetCookies(u, cookies)

	header := make(http.Header)
	header.Set("User-Agent", config.Req.UA)
	header.Set("Accept", config.Req.Accept)
	//	header.Set("Accept-Encoding", config.Req.AcceptEncoding)
	header.Set("Accept-Language", config.Req.AcceptLanguage)
	headerTp := &HeaderTransport{
		Upstream: http.DefaultTransport,
		Header:   header,
	}

	tp, err := scraper.NewTransport(headerTp, config.Req.UA, jar)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: tp,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.String() == config.Req.Page503 {
				return ErrFake503
			}
			return nil
		},
	}

	b := &Bmber{
		config: config,
		client: client,
		logger: logger,

		jar: jar,
	}

	return b, nil
}
