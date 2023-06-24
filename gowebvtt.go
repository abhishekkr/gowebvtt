package gowebvtt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type VTT struct {
	Scenes []Scene
}

type Scene struct {
	StartMilliSec uint64
	EndMilliSec   uint64
	Transcript    []string
}

var (
	MaxTokensOnAFrame = 5
	SplitForMaxTokens = false
)

var (
	RegHeader    = regexp.MustCompile("^\\s*WEBVTT\\s*")
	RegNote      = regexp.MustCompile("^\\s*NOTE\\s")
	RegTimestamp = regexp.MustCompile("^\\s*[0-9\\:\\.]+\\s+\\-+>\\s+[0-9\\:\\.]+\\s*$")
)

func (vtt *VTT) AppendScene(scene Scene) bool {
	if len(scene.Transcript) == 0 || (scene.EndMilliSec-scene.StartMilliSec == 0) {
		return false
	}
	vtt.Scenes = append(vtt.Scenes, scene)
	return true
}

func (scene *Scene) AppendTranscript(line string) bool {
	if len(line) == 0 {
		return false
	}
	scene.Transcript = append(scene.Transcript, line)
	return true
}

func (scene *Scene) ProcessSubtext(line string) {
	if !SplitForMaxTokens {
		scene.AppendTranscript(line)
		return
	}
	words := strings.Split(strings.Trim(line, " "), " ")
	var idx int
	for idx = 0; idx+MaxTokensOnAFrame < len(words); idx += MaxTokensOnAFrame {
		tmpTxt := strings.Join(words[idx:idx+MaxTokensOnAFrame], " ")
		scene.AppendTranscript(tmpTxt)
	}
	if idx < len(words) {
		tmpTxt := strings.Join(words[idx:len(words)], " ")
		scene.AppendTranscript(tmpTxt)
	}
}

func ParseWebVTT(fileScanner *bufio.Scanner) VTT {
	var vtt = VTT{Scenes: []Scene{}}

	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		if RegHeader.MatchString(fileScanner.Text()) {
			break
		}
	}

	var isNote, isSubtext bool
	var scene, muteScene Scene
	var prevEndMilliSec uint64

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "" {
			if isNote {
				isNote = false
			} else if isSubtext {
				isSubtext = false
				vtt.AppendScene(scene)
				prevEndMilliSec = scene.EndMilliSec
				scene = Scene{StartMilliSec: 0, EndMilliSec: 0, Transcript: []string{}}
			}
			continue
		} else if RegNote.MatchString(line) {
			isNote = true
		} else if isSubtext {
			scene.ProcessSubtext(line)
		} else if RegTimestamp.MatchString(line) {
			var errTimeRange error
			scene.StartMilliSec, scene.EndMilliSec, errTimeRange = getTimeRange(line)
			if errTimeRange != nil {
				log.Printf("TimeRange parsing failed: %v", errTimeRange)
			}
			if scene.StartMilliSec-prevEndMilliSec > 50 {
				muteScene = Scene{StartMilliSec: prevEndMilliSec, EndMilliSec: scene.StartMilliSec}
				vtt.AppendScene(muteScene)
			}
			isSubtext = true
		}
	}
	if isSubtext {
		vtt.Scenes = append(vtt.Scenes, scene)
	}
	return vtt
}

func ParseFile(filepath string) VTT {
	readFile, err := os.Open(filepath)
	defer readFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	return ParseWebVTT(fileScanner)
}

func MillsecToVTTTimeString(ms uint64) string {
	seconds := uint64(ms / 1000)
	minutes := uint64(seconds / 60)
	hours := uint64(minutes / 60)
	if hours < 1 {
		return fmt.Sprintf(
			"%02d:%02d.%03d",
			minutes%60,
			seconds%60,
			ms%1000,
		)
	}
	return fmt.Sprintf(
		"%02d:%02d:%02d.%03d",
		hours,
		minutes%60,
		seconds%60,
		ms%1000,
	)
}

func String(vtt VTT) string {
	var result = "WEBVTT\n"
	for _, s := range vtt.Scenes {
		if len(s.Transcript) > 0 {
			result = result + "\n" + MillsecToVTTTimeString(s.StartMilliSec) + " --> " + MillsecToVTTTimeString(s.EndMilliSec) + "\n"
		}
		for _, w := range s.Transcript {
			result = result + w + "\n"
		}
	}
	return result
}

func Println(vtt VTT) {
	for _, s := range vtt.Scenes {
		fmt.Println("\nStarts:", s.StartMilliSec, ",\tEnds:", s.EndMilliSec)
		for _, w := range s.Transcript {
			fmt.Println(w)
		}
	}
}
