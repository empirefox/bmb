package bmbconfig

import (
	"github.com/empirefox/bmb/moment"
	"github.com/empirefox/cement/clog"
)

type Bmb struct {
	Dev           bool `env:"DEV"`
	SkipVerifyTLS bool `env:"SKIP_VERIFY_TLS"`
}

type Req struct {
	UA             string `default:"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:49.0) Gecko/20100101 Firefox/49.0"`
	Accept         string `default:"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"`
	AcceptEncoding string `default:"gzip"` // TODO not used now
	AcceptLanguage string `default:"zh-CN,zh;q=0.8"`
	PageHome       string `default:"https://www.bitmain.com/"`
	Page503        string `default:"https://shop.bitmain.com:8080/503.html"`
}

type Login struct {
	Moment   moment.Moment
	Page     string `default:"https://passport.bitmain.com/login?service=https%3A%2F%2Fwww.bitmain.com%2Fuser%2ForderDetails.htm"`
	Post     string `default:"https://passport.bitmain.com/login?service=https%3A%2F%2Fwww.bitmain.com%2Fuser%2ForderDetails.htm"`
	Username string `env:"LOGIN_USER"`
	Password string `env:"LOGIN_PASSWORD"`
}

type AnalyzePid struct {
	PageProducts string `default:"https://shop.bitmain.com/main.htm?lang=zh"`
	PeriodSecond uint
	NotIn        []string
	H2Contain    string
	H3Contain    string
}

type Addtocart struct {
	Moment moment.Moment
	// pid first, count second
	GetFmt            string `default:"https://shop.bitmain.com/user/orderDetails.htm?m=add&pid=%s&count=%d&fitting="`
	Pid               string
	Count             uint
	OkIfContain       string `default:"商品已成功加入购物车"`
	PageProductPrefix string `default:"https://shop.bitmain.com/productDetail.htm?pid="`
	PageCart          string `default:"https://shop.bitmain.com/user/orderDetails.htm"`
	PostModCount      string `default:"https://shop.bitmain.com/user/orderDetails.htm?m=count"`
}

type Shopcart struct {
	Post string `default:"https://shop.bitmain.com/user/orderDetails.htm?m=y"`
}

type Config struct {
	Schema     string `json:"-" yaml:"-" toml:"-"`
	Bmb        Bmb
	Clog       clog.Config
	Req        Req
	Login      Login
	AnalyzePid AnalyzePid
	Addtocart  Addtocart
	Shopcart   Shopcart
}

func (c *Config) GetEnvPtrs() []interface{} {
	return []interface{}{&c.Bmb, &c.Clog}
}
