package gtr

import (
	"errors"
	"sync"

	"github.com/xiroxasx/gtr/pkg/locale"
)

var (
	ErrInvalidActiveLocale = errors.New("invalid active locale")
)

type Translator struct {
	translations translations
	activeLocale locale.Locale
	mx           *sync.Mutex
}

type translations struct {
	mappings map[locale.Locale]dictionary
	mx       *sync.Mutex
}

type dictionary map[string]string

func New() *Translator {
	return &Translator{
		translations: translations{
			mappings: make(map[locale.Locale]dictionary, 0),
			mx:       &sync.Mutex{},
		},
		mx: &sync.Mutex{},
	}
}

// Use sets the active locale.
func (t *Translator) Use(loc locale.Locale) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.activeLocale = loc
}

func (t *Translator) Active() locale.Locale {
	return t.activeLocale
}

// Register adds a new dictionary entry for the currently active locale.
// Registering the key twice, will overwrite the old value.
func (t *Translator) Register(key, val string) (err error) {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.activeLocale == (locale.Locale{}) {
		err = ErrInvalidActiveLocale
		return
	}
	t.registerFor(t.activeLocale, key, val)
	return
}

// RegisterFor adds a new dictionary entry for the given locale.
// Registering the key twice, will overwrite the old value.
func (t *Translator) RegisterFor(loc locale.Locale, key, val string) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.registerFor(loc, key, val)
}

// Get gets the translation of the current active locale.
// If no locale is set ot the key has not been found, the returned value is an empty string.
func (t *Translator) Get(k string) string {
	t.translations.mx.Lock()
	defer t.translations.mx.Unlock()

	return t.translations.mappings[t.activeLocale][k]
}

// Load sets the translations for the given locale.
// If replace is true, existing translations will be overwritten.
func (t *Translator) Load(loc locale.Locale, replace bool, dict map[string]string) {
	t.mx.Lock()
	defer t.mx.Unlock()

	d, ok := t.translations.mappings[loc]
	if !ok {
		d = dictionary{}
	}

	if replace {
		for k, v := range dict {
			d[k] = v
		}
	} else {
		for k, v := range dict {
			if !d.hasKey(k) {
				d[k] = v
			}
		}
	}

	t.translations.mappings[loc] = d
}

// ClearAll clears all dictionaries.
func (t *Translator) ClearAll() {
	t.translations.mx.Lock()
	defer t.translations.mx.Unlock()

	t.translations.mappings = make(map[locale.Locale]dictionary, 0)
}

// Clear clears the dictionary for the currently active locale.
func (t *Translator) Clear(loc locale.Locale) (err error) {
	t.translations.mx.Lock()
	defer t.translations.mx.Unlock()

	if t.activeLocale == (locale.Locale{}) {
		err = ErrInvalidActiveLocale
		return
	}

	t.translations.mappings[t.activeLocale] = make(dictionary, 0)
	return
}

// ClearFor clears the dictionary for the given locale.
func (t *Translator) ClearFor(loc locale.Locale) {
	t.translations.mx.Lock()
	defer t.translations.mx.Unlock()

	t.translations.mappings[loc] = make(dictionary, 0)
}

// HasKey checks whether a given key exists for the active translation.
func (t *Translator) HasKey(k string) bool {
	return t.translations.mappings[t.activeLocale].hasKey(k)
}

// HasKey checks whether a given key exists for the given translation.
func (t *Translator) HasKeyIn(loc locale.Locale, k string) bool {
	return t.translations.mappings[loc].hasKey(k)
}

// HasValue checks whether a given value exists for the active translation.
func (t *Translator) HasValue(v string) (ok bool, k string) {
	return t.translations.mappings[t.activeLocale].hasValue(v)
}

// HasValue checks whether a given value exists for the given translation.
func (t *Translator) HasValueIn(loc locale.Locale, v string) (ok bool, k string) {
	return t.translations.mappings[loc].hasValue(v)
}

// registerFor registers a translation in a specific locale.
// Caller must ensure to lock t beforehand!
func (t *Translator) registerFor(loc locale.Locale, k, v string) {
	_, ok := t.translations.mappings[loc]
	if !ok {
		t.translations.mappings[loc] = make(dictionary, 0)
	}
	t.translations.mappings[loc][k] = v
}

func (d dictionary) hasKey(k string) (ok bool) {
	_, ok = d[k]
	return
}

func (d dictionary) hasValue(v string) (ok bool, k string) {
	for key, val := range d {
		if val == v {
			return true, key
		}
	}
	return
}
