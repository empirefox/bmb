package moment

import "testing"

func TestMoment(t *testing.T) {
	var m Moment
	err := m.UnmarshalText([]byte("2017-08-28 11:30:56"))
	if err != nil {
		t.Fatal(err)
	}

	b, err := m.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "2017-08-28 11:30:56" {
		t.Errorf("monent should be 2017-08-28 11:30:56, but got %s", b)
	}
}
