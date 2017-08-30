package bmber

import (
	"bytes"

	"go.uber.org/zap"

	"github.com/PuerkitoBio/goquery"
)

func (b *Bmber) StepPidInfo() error {
	defer func() {
		if err := recover(); err != nil {
			b.logger.Error("StepPidInfo recover", zap.Any(err))
		}
	}()

	body, err := b.GetPageBytes(b.config.Addtocart.PageProductPrefix + b.config.Addtocart.Pid)
	if err != nil {
		b.logger.Error("PID info failed", zap.Error(err))
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		b.logger.Error("PID HTML", zap.Error(err))
		return err
	}

	b.logger.Info(
		"PID INFO",
		zap.String("Price", doc.Find(".u-p").Text()),
		zap.String("MinCount", doc.Find("#productBuyMinCount").Text()),
		zap.String("MaxCount", doc.Find("#productBuyMaxCount").Text()),
		zap.String("PerMaxCount", doc.Find("#productBuyPerMaxCount").Text()),
		zap.String("RemainCount", doc.Find("#productRemainCount").Text()),
	)
}
