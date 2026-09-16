package policy

import "time"

type Policy struct {
	Name          string
	FirstResponse time.Duration
	Resolve       time.Duration
	WarnRatio     float64
}

func (p Policy) WarnAt(start time.Time, deadline time.Time) time.Time {
	total := deadline.Sub(start)
	warnAfter := time.Duration(float64(total) * p.WarnRatio)
	return start.Add(warnAfter)
}
