package dfuse

import (
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	u := t.Unix()
	return []byte(strconv.Itoa(int(u))), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	secs, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*t = Timestamp{time.Unix(secs, 0)}
	return nil
}
