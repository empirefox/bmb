package moment

import (
	"time"
)

var Format = "2006-01-02 15:04:05"

type Moment time.Time

func (m *Moment) MarshalText() (text []byte, err error) {
	return []byte(time.Time(*m).Format(Format)), nil
}

func (m *Moment) UnmarshalText(text []byte) error {
	t, err := time.Parse(Format, string(text))
	if err != nil {
		return err
	}
	*m = Moment(t)
	return nil
}
