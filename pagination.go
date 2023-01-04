package monzo

import (
	"net/url"
	"strconv"
	"time"
)

type Pagination struct {
	Limit  int
	Since  string
	Before string
}

func (p Pagination) Values(vals ...url.Values) url.Values {
	v := url.Values{}

	if len(vals) > 0 {
		v = vals[0]
	}

	if p.Limit != 0 {
		v.Add("limit", strconv.Itoa(p.Limit))
	}

	if p.Since != "" {
		v.Add("since", p.Since)
	}

	if p.Before != "" {
		v.Add("before", p.Before)
	}

	return v
}

func (p Pagination) SinceTime() time.Time {
	if p.Since == "" {
		return time.Time{}
	}

	t, _ := time.Parse(time.RFC3339Nano, p.Since)
	return t
}

func (p Pagination) BeforeTime() time.Time {
	if p.Before == "" {
		return time.Now()
	}

	t, _ := time.Parse(time.RFC3339Nano, p.Before)
	return t
}
