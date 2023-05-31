package locale

import (
	"errors"
	"strings"
)

const localeValueSeparator = "_"

var (
	ErrInvalidLocale = errors.New("invalid locale")
)

type Locale struct {
	// Language is the 2 letter abbriviation of the spoken language (e.g. "en").
	Language string
	// Territory is the 2 letter abbriviation of the country (e.g. "US").
	Territory string
}

func New(lang, ter string) Locale {
	return Locale{
		Language:  strings.ToLower(lang),
		Territory: strings.ToUpper(ter),
	}
}

func NewFromString(locValue string) (l Locale, err error) {
	locs := strings.Split(locValue, localeValueSeparator)
	if len(locs) == 1 {
		err = ErrInvalidLocale
		return
	}

	l = Locale{
		Language:  strings.ToLower(locs[0]),
		Territory: strings.ToUpper(locs[1]),
	}
	return
}

func (l *Locale) String() string {
	return l.Language + localeValueSeparator + l.Territory
}
