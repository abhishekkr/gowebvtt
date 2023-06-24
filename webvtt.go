package gowebvtt

import (
	"errors"
	"log"
	"math"
	"regexp"
	"strings"
)

type VTT struct {
	Scenes  []Scene
	Options VttOptions
}

type Scene struct {
	StartMilliSec uint64
	EndMilliSec   uint64
	Transcript    []string
	Options       SceneOptions
}

type VttOptions struct {
	Enabled          bool
	MaxLinesPerScene int
	SceneOptions     SceneOptions
}

type SceneOptions struct {
	Enabled         bool
	MaxCharsPerLine int
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

func newScene(startMs, endMs uint64, opts SceneOptions) Scene {
	return Scene{
		StartMilliSec: startMs,
		EndMilliSec:   endMs,
		Transcript:    []string{},
		Options:       opts,
	}
}

func (vtt *VTT) AppendScene(scene Scene) bool {
	if len(scene.Transcript) == 0 || (scene.EndMilliSec-scene.StartMilliSec == 0) {
		return false
	}
	if vtt.Options.Enabled && len(scene.Transcript) > vtt.Options.MaxLinesPerScene {
		if errSplit := vtt.splitAndAppendScene(scene); errSplit == nil {
			return true
		} else {
			log.Println(errSplit)
		}
	}
	vtt.Scenes = append(vtt.Scenes, scene)
	return true
}

func (vtt *VTT) splitAndAppendScene(scene Scene) error {
	transcriptLen := len(scene.Transcript)
	startMs := scene.StartMilliSec
	currEndMs := startMs
	endMs := scene.EndMilliSec
	sceneCount := math.Ceil(float64(transcriptLen) / float64(vtt.Options.MaxLinesPerScene))
	if sceneCount <= 1.0 {
		vtt.Scenes = append(vtt.Scenes, scene)
		return nil
	}
	splitDurationMs := uint64((endMs - startMs) / uint64(sceneCount))
	for idx := 0; idx < transcriptLen; {
		idxUpto := idx + vtt.Options.MaxLinesPerScene
		startMs = currEndMs
		currEndMs += splitDurationMs
		if sceneCount == 1 {
			currEndMs = endMs
			idxUpto = transcriptLen
		}
		sceneX := newScene(startMs, currEndMs, scene.Options)
		sceneX.Transcript = scene.Transcript[idx:idxUpto]
		vtt.Scenes = append(vtt.Scenes, sceneX)
		idx += vtt.Options.MaxLinesPerScene
		sceneCount--
	}
	if currEndMs != endMs {
		return errors.New("splitAndAppendScene failed with last scene EndMilliSec not being same as input scene's EndMilliSec or sceneCount not being completed")
	}
	return nil
}

func (vtt *VTT) String() string {
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
