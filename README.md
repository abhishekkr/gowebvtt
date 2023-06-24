## gowebvtt

> a Go package to parse **WebVTT subtitles** (Web Video Text Tracks Format) content 

* This package currently supports main features i.e. Time Marker, Subtitles & same-line Notes.

* To parse use `gowebvtt.ParseFile(..)`

* Example usage with Options

```
package main

import (
        "github.com/abhishekkr/gowebvtt"
)

func main() {
        vttOpts := gowebvtt.VttOptions{Enabled: true, MaxLinesPerScene: 2}
        vtt, err := gowebvtt.ParseFileWithOptions("sample.vtt", vttOpts)
        if errVtt != nil {
                panic(errVtt)
        }
        gowebvtt.Println(vtt)
}
```

---

### ToDo

> [source: VTT Mozilla doc](https://developer.mozilla.org/en-US/docs/Web/API/WebVTT_API)

* add support for multi-line NOTE

* add support for Style & Cue blocks

---
