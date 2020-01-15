package dfuse

import "time"

type progress struct {
	reqId    string
	interval time.Duration
	nextTime time.Time
}

func newProgress(reqId string, interval time.Duration) progress {
	return progress{
		reqId:    reqId,
		interval: interval,
		nextTime: time.Now().Add(interval),
	}
}

func (p *progress) refreshTime() {
	p.nextTime = p.nextTime.Add(p.interval)
}

func (p *progress) isTimeout(nowTime time.Time) bool {
	return nowTime.After(p.nextTime)
}
