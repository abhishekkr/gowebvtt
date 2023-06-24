package gowebvtt

import (
	"bufio"
	"strings"
	"testing"
)

var (
	sample = `WEBVTT

00:01.000 --> 00:04.000
This is subtitle at 1sec to 4sec.

00:05.000 --> 00:09.000
— These are split subtitles from marker of 5sec.
— Upto 9sec, split equally.

00:10.000 --> 00:14.000
The text could be as longs as you want or as small as you want, file starts with WEBVTT marker. Can hold comments starting with NOTE. The lines above with --> are time markers for when subtitles need to be displayed.
`
)

func TestParseWebVTT(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(sample))
	vtt := ParseWebVTT(scanner, VttOptions{})
	if len(vtt.Scenes) != 3 {
		t.Fatalf("Failed for ParseWebVTT with sample content, for Scene parsing.")
	}
	subs1 := vtt.Scenes[0]
	if subs1.StartMilliSec != 1000 || subs1.EndMilliSec != 4000 {
		t.Fatalf("Failed for ParseWebVTT with sample content, for time parsing.")
	}
	if subs1.Transcript[0] != "This is subtitle at 1sec to 4sec." {
		t.Fatalf("Failed for ParseWebVTT with sample content, for transcript.")
	}
	subs2 := vtt.Scenes[1]
	if subs2.StartMilliSec != 5000 || subs2.EndMilliSec != 9000 {
		t.Fatalf("Failed for ParseWebVTT with sample content, for time parsing.")
	}
	if len(subs2.Transcript) != 2 {
		t.Fatalf("Failed for ParseWebVTT with sample content, for transcript.")
	}
	if String(vtt) != sample {
		t.Fatalf("Failed for ParseWebVTT with sample content.")
	}
}

func TestParseFile(t *testing.T) {
	vtt := ParseFile("sample.vtt")
	if String(vtt) != sample {
		t.Fatalf("Failed for ParseFile with sample.vtt content.")
	}
}
