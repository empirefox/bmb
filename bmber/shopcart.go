package bmber

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

func (b *Bmber) StepShopcart() error {
	v := make(url.Values)
	v.Set("shopCarts", b.cartid)

	req, err := http.NewRequest("POST", b.config.Shopcart.Post, strings.NewReader(v.Encode()))
	if err != nil {
		b.logger.Error("Shopcart POST Request err", zap.Error(err))
		return err
	}

	u, err := url.Parse(b.config.Addtocart.PageCart)
	if err != nil {
		b.logger.Error("PageCart to URL err", zap.Error(err))
		return err
	}
	req.Header.Set("Origin", fmt.Sprintf("%s://%s", u.Scheme, u.Host))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", b.config.Login.Page)
	if _, err = b.client.Do(req); err != nil {
		b.logger.Error("Shopcart POST err", zap.Error(err))
		return err
	}
	return nil
}
