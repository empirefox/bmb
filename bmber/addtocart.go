package bmber

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"go.uber.org/zap"
)

func (b *Bmber) StepAddtocart() error {
	res, err := b.client.Get(fmt.Sprintf(b.config.Addtocart.GetFmt, b.config.Addtocart.Pid, b.config.Addtocart.Count))
	if err != nil {
		b.logger.Error("Addtocart call err", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	// cartid not here
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		b.logger.Error("Addtocart body err", zap.Error(err))
		return err
	}
	if !bytes.Contains(body, []byte(b.config.Addtocart.OkIfContain)) {
		b.logger.Warn("Addtocart wrong body, must check cart!!!", zap.Error(err))
	}

	// open cart page
	body, err := b.GetPageBytes(b.config.Addtocart.PageCart)
	if err != nil {
		b.logger.Error("PageCart failed", zap.Error(err))
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		b.logger.Error("PageCart HTML", zap.Error(err))
		return err
	}

	// get count
	cstr, ok := doc.Find(".count-input").Attr("value")
	if !ok {
		b.logger.Error("count value not found")
		return errors.New("count value not found")
	}
	count, err := strconv.Atoi(cstr)
	if err != nil {
		b.logger.Error("count is not number", zap.Error(err), zap.String("count", cstr))
		return err
	}

	// check count equal
	if uint(count) < b.config.Addtocart.Count {
		b.logger.Error("count value not equal", zap.Uint("should", b.config.Addtocart.Count), zap.Int("got", count))
		return errors.New("count value not equal")
	}

	// get remain
	crc, ok := doc.Find(".remainCount").Attr("value")
	if !ok {
		b.logger.Error("remainCount not found")
		return errors.New("remainCount not found")
	}
	remain, err := strconv.Atoi(crc)
	if err != nil {
		b.logger.Error("remainCount is not number", zap.Error(err), zap.String("remain", crc))
		return err
	}

	if remain == 0 {
		b.logger.Warn("REMAIN not enough!!!!!!")
	}

	var needModCount bool
	if uint(count) > b.config.Addtocart.Count {
		b.logger.Warn("Adjust count because multi add!!!")
		needModCount = true
		count = b.config.Addtocart.Count
	}
	if remain > 0 && count > remain {
		b.logger.Warn("Adjust count because remain not enough!!!", zap.Int("remain", remain), zap.Int("count", count))
		needModCount = true
		b.config.Addtocart.Count = remain
		count = remain
	}

	// get/set cartId
	cartid, ok := doc.Find(".delete").Attr("cartId")
	if !ok {
		b.logger.Error("cartId not found")
		return errors.New("cartId not found")
	}
	b.cartid = cartid

	// mod count
	if needModCount {
		v := make(url.Values)
		v.Set("id", cartid)
		v.Set("count", strconv.Itoa(count))
		v.Set("ty", "set")

		req, err := http.NewRequest("POST", b.config.Addtocart.PostModCount, strings.NewReader(v.Encode()))
		if err != nil {
			b.logger.Error("ModCount POST Request err", zap.Error(err))
			return err
		}

		u, err := url.Parse(b.config.Addtocart.PageCart)
		if err != nil {
			b.logger.Error("PageCart to URL err", zap.Error(err))
			return err
		}
		req.Header.Set("Origin", fmt.Sprintf("%s://%s", u.Scheme, u.Host))

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Referer", b.config.Addtocart.PageCart)
		if _, err = b.client.Do(req); err != nil {
			b.logger.Error("ModCount POST err", zap.Error(err))
			return err
		}
	}

	return nil
}
