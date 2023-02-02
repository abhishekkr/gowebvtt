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

func TestProcessSubtextWithSplit(t *testing.T) {
	scene := Scene{}
	SplitForMaxTokens = true
	scene.ProcessSubtext("one two three four five six seven")
	SplitForMaxTokens = false
	if len(scene.Transcript) < 2 {
		t.Fatalf("Failed for ProcessSubtext, didn't split longer line.")
	}
	for _, sub := range scene.Transcript {
		if len(strings.Split(sub, " ")) > MaxTokensOnAFrame {
			t.Fatalf("Failed for ProcessSubtext, didn't split correct length.")
		}
	}
}

func TestParseWebVTT(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(sample))
	vtt := ParseWebVTT(scanner)
	subs := vtt.Scenes[1] // 0th 2nd 4th index have non-sub scenes
	if subs.StartMilliSec != 1000 || subs.EndMilliSec != 4000 {
		t.Fatalf("Failed for ParseWebVTT with sample content, for time parsing.")
	}
	if subs.Transcript[0] != "This is subtitle at 1sec to 4sec." {
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
