package bmber

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrPidNotFound = errors.New("pid not found")
)

type product struct {
	pid string
	h2  string
	h3  string
}

func (b *Bmber) StepAnalyzePid() error {
	for b.config.Addtocart.Pid == "" {
		err := b.analyzePid()
		if err == ErrPidNotFound {
			time.Sleep(time.Second * time.Duration(b.config.AnalyzePid.PeriodSecond))
			continue
		}
		if err != nil {
			return err
		}
		break
	}
	return nil
}

func (b *Bmber) analyzePid() error {
	body, err := b.GetPageBytes(b.config.AnalyzePid.PageProducts)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	var ps []product
	doc.Find(".homepage-prtList ul li").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find(".goodsImg a").Attr("href")
		if !ok {
			b.logger.Error("Cannot find '.goodsImg a'", zap.String("href", href))
			return
		}

		hrefs := strings.Split(href, "=")
		if len(hrefs) != 2 {
			b.logger.Error("Extract pid from href failed", zap.String("href", href))
			return
		}

		pid := hrefs[1]
		for _, not := range b.config.AnalyzePid.NotIn {
			if not == pid {
				b.logger.Error("Pid found but in blacklist", zap.String("pid", pid))
				return
			}
		}

		h2 := s.Find(".goodsDescrip h2").Text()
		h3 := s.Find(".goodsDescrip h3").Text()
		if strings.Contains(h2, b.config.AnalyzePid.H2Contain) && strings.Contains(h3, b.config.AnalyzePid.H3Contain) {
			ps = append(ps, product{pid, h2, h3})
		}
	})

	b.logger.Info("PID FOUND", zap.Int(len(ps)))
	for _, p := range ps {
		b.logger.Info("PID FOUND ====> ", zap.String("pid", p.pid), zap.String("h2", p.h2), zap.String("h3", p.h3))
	}

	switch len(ps) {
	case 0:
		return ErrPidNotFound
	case 1:
		p := ps[0]
		b.logger.Info("USING PID ====> ", zap.String("pid", p.pid), zap.String("h2", p.h2), zap.String("h3", p.h3))
		b.config.Addtocart.Pid = p.pid
		return nil
	}

	b.logger.Error("More than 1 pid found, please choose 1 to set")
	return errors.New("pid wrong")
}
