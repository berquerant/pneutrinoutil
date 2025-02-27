package pathx

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func NewResultElement(basename string, now time.Time, salt int) *ResultElement {
	return &ResultElement{
		Basename: basename,
		Now:      now,
		Salt:     salt,
	}
}

type ResultElement struct {
	Basename string
	Now      time.Time
	Salt     int
}

func (r ResultElement) String() string {
	return fmt.Sprintf(
		"%s__%04d%02d%02d%02d%02d%02d_%d_%d",
		r.Basename,
		r.Now.Year(),
		r.Now.Month(),
		r.Now.Day(),
		r.Now.Hour(),
		r.Now.Minute(),
		r.Now.Second(),
		r.Now.Unix(),
		r.Salt,
	)
}

var (
	resultElementRegexp   = regexp.MustCompile(`^(.+)__([0-9]+)_([0-9]+)_([0-9]+)$`)
	ErrParseResultElement = errors.New("ParseResultElement")
)

const (
	// %Y%m%d%H%M%S
	resultElementTimeFormat = "20060102150405"
)

func ParseResultElement(s string) (*ResultElement, error) {
	matched := resultElementRegexp.FindStringSubmatch(s)
	if len(matched) != 5 {
		return nil, fmt.Errorf("%w: invalid format", ErrParseResultElement)
	}
	matched = matched[1:]

	basename := matched[0]

	now, err := time.Parse(resultElementTimeFormat, matched[1])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid time string", errors.Join(ErrParseResultElement, err))
	}

	if _, err := strconv.ParseInt(matched[2], 10, 64); err != nil {
		return nil, fmt.Errorf("%w: invalid timestamp", errors.Join(ErrParseResultElement, err))
	}

	salt, err := strconv.Atoi(matched[3])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid salt", errors.Join(ErrParseResultElement, err))
	}

	return &ResultElement{
		Basename: basename,
		Now:      now,
		Salt:     salt,
	}, nil
}
