package gtr

import (
	"testing"

	r "github.com/stretchr/testify/require"
	"github.com/xiroxasx/gtr/pkg/locale"
)

func TestTranslator(t *testing.T) {
	tr := New()
	r.Exactly(t, locale.Locale{}, tr.activeLocale)
	_, err := locale.NewFromString("deDE")
	r.ErrorIs(t, locale.ErrInvalidLocale, err)
	r.Exactly(t, locale.Locale{}, tr.Active())
	r.ErrorIs(t, ErrInvalidActiveLocale, tr.Register("test", "Test"))

	trKey := "test"
	trKey2 := "test2"
	trVal := "Test"
	trVal2 := "Test2"
	loc := "de_DE"
	l, err := locale.NewFromString(loc)
	r.NoError(t, err)
	r.Empty(t, tr.Get("NotExisting"))
	r.Empty(t, tr.Get(""))
	r.Error(t, tr.Register(trKey, trVal))

	tr.Use(l)
	r.NoError(t, tr.Register(trKey, trVal))
	tr.registerFor(l, trKey2, trVal2)
	r.Exactly(t, trVal, tr.Get(trKey))
	r.Exactly(t, trVal2, tr.Get(trKey2))
}
