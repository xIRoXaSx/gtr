# gtr
Simplify translations with gtr.  

## Examples
Create a new translator, register a translation and print it:  

```go
package x

import (
	"fmt"

	"github.com/xiroxasx/gtr"
)

func main() {
    tr := gtr.New(gtr.NewLocale("de", "DE"))
    tr.Register("test", "Test")
    fmt.Println("translated:", tr.Get("test"))
}
```

---

Load multiple translations for multiple languages.

```go
package x

import (
	"fmt"

	"github.com/xiroxasx/gtr"
)

func main() {
    de := gtr.NewLocale("de", "DE")
    en := gtr.NewLocale("en", "US")
    es := gtr.NewLocale("es", "ES")
    tr := gtr.New(de)

    // Maps could also be marshalled from files ect...
    tr.Load(de, false, map[string]string{
		"language": "sprache",
        "test":     "test",
	})
    tr.Load(en, false, map[string]string{
		"language": "language",
        "test":     "test",
	})
    tr.Load(es, false, map[string]string{
		"language": "idioma",
        "test":     "prueba",
	})
    fmt.Println("translated:", tr.Get("language"))        // Prints "translated: sprache"
    fmt.Println("translated:", tr.GetFor(es, "language")) // Prints "translated: idioma"
}
```