package gtr_test

import (
	"fmt"
	"testing"

	r "github.com/stretchr/testify/require"
	"github.com/xiroxasx/gtr"
)

func TestTranslator(t *testing.T) {
	tr := gtr.New()
	r.Exactly(t, gtr.Locale{}, tr.Active())
	_, err := gtr.NewFromString("deDE")
	r.ErrorIs(t, gtr.ErrInvalidLocale, err)
	r.Exactly(t, gtr.Locale{}, tr.Active())
	r.ErrorIs(t, gtr.ErrInvalidActiveLocale, tr.Register("test", "Test"))

	trLen := 3
	createDummyEntries := func(tr *gtr.Translator, baseKey, baseValue string) {
		locs := []gtr.Locale{
			gtr.NewLocale("de", "DE"),
			gtr.NewLocale("en", "US"),
		}
		for _, l := range locs {
			for i := 0; i < trLen; i++ {
				tr.RegisterFor(l, fmt.Sprintf("%s%d", baseKey, i), fmt.Sprintf("%s%d", baseValue, i))
			}
		}
	}

	const trKey = "test"
	const trVal = "Test"
	const loc = "de_DE"
	deDE, err := gtr.NewFromString(loc)
	enUS := gtr.NewLocale("en", "US")
	lastKey := fmt.Sprintf("%s%d", trKey, trLen-1)
	lastValue := fmt.Sprintf("%s%d", trVal, trLen-1)
	r.NoError(t, err)
	r.Empty(t, tr.Get("NotExisting"))
	r.Empty(t, tr.Get(""))
	r.Error(t, tr.Register(trKey, trVal))
	r.ErrorIs(t, gtr.ErrInvalidActiveLocale, tr.Clear())

	tr.Use(deDE)
	r.Exactly(t, loc, tr.Active().String())
	r.NoError(t, tr.Register(trKey, trVal))
	createDummyEntries(tr, trKey, trVal)

	r.Exactly(t, trVal, tr.Get(trKey))
	r.Exactly(t, trLen+1, tr.Len())
	r.Exactly(t, trLen, tr.LenFor(enUS))
	r.Empty(t, tr.Get(fmt.Sprintf("%s%d", trKey, trLen)))
	r.Exactly(t, fmt.Sprintf("%s%d", trVal, trLen-1), tr.Get(lastKey))

	r.NoError(t, tr.Clear())
	r.Exactly(t, 0, tr.Len())
	tr.ClearFor(enUS)
	r.Exactly(t, 0, tr.LenFor(enUS))

	createDummyEntries(tr, trKey, trVal)
	r.False(t, tr.HasValue("NotExisting"))
	r.False(t, tr.HasValueFor(enUS, "NotExisting"))
	r.True(t, tr.HasKey(lastKey))
	r.True(t, tr.HasValue(lastValue))
	r.True(t, tr.HasKeyFor(enUS, lastKey))
	r.True(t, tr.HasValueFor(enUS, lastValue))

	const newVal = "new"
	dict := map[string]string{
		trKey:   trVal,
		lastKey: lastValue,
	}
	tr.ClearAll()
	r.Exactly(t, 0, tr.LenFor(deDE))
	r.Exactly(t, 0, tr.LenFor(enUS))
	tr.Load(enUS, false, dict)
	r.Exactly(t, len(dict), tr.LenFor(enUS))
	tr.Load(enUS, true, map[string]string{trKey: newVal})
	r.Exactly(t, newVal, tr.GetFor(enUS, trKey))
	r.Exactly(t, dict[lastKey], tr.GetFor(enUS, lastKey))
	tr.Use(enUS)
	r.Exactly(t, enUS, tr.Active())
	r.Exactly(t, dict[lastKey], tr.Get(lastKey))
	r.Exactly(t, len(dict), tr.LenFor(enUS))
	r.Exactly(t, 0, tr.LenFor(deDE))
	tr.Use(deDE)
	r.Exactly(t, deDE, tr.Active())
}
