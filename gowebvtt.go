package gowebvtt

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

const (
	TS_HMSMs uint8 = iota
	TS_MSMs
	TS_SMs
)

var (
	MaxTokensOnAFrame = 5
	SplitForMaxTokens = false
)

var (
	RegHeader    = regexp.MustCompile("^\\s*WEBVTT\\s*")
	RegNote      = regexp.MustCompile("^\\s*NOTE\\s")
	RegTimestamp = regexp.MustCompile("^\\s*[0-9\\:\\.]+\\s*[\\-]+>\\s*[0-9\\:\\.]+\\s*$")
	RegTimeHMSMs = regexp.MustCompile("\\s*([0-9]+)\\:([0-9]+)\\:([0-9]+)\\.([0-9]+)\\s*")
	RegTimeMSMs  = regexp.MustCompile("\\s*([0-9]+)\\:([0-9]+)\\.([0-9]+)\\s*")
	RegTimeSMs   = regexp.MustCompile("\\s*([0-9]+)\\.([0-9]+)\\s*")
)

func getTimeRangeTokens(txt string) (uint8, [][]string) {
	if RegTimeHMSMs.MatchString(txt) {
		return TS_HMSMs, RegTimeHMSMs.FindAllStringSubmatch(txt, -1)
	} else if RegTimeMSMs.MatchString(txt) {
		return TS_MSMs, RegTimeMSMs.FindAllStringSubmatch(txt, -1)
	}
	return TS_SMs, RegTimeSMs.FindAllStringSubmatch(txt, -1)
}

func getTimeRange(txt string) (start, end uint64) {
	tsType, timeSubstr := getTimeRangeTokens(txt)
	if tsType == TS_HMSMs {
		startHr, _ := strconv.ParseUint(timeSubstr[0][1], 10, 64)
		startMin, _ := strconv.ParseUint(timeSubstr[0][2], 10, 64)
		startSec, _ := strconv.ParseUint(timeSubstr[0][3], 10, 64)
		startMsec, _ := strconv.ParseUint(timeSubstr[0][4], 10, 64)
		start = (((startHr * 3600) + (startMin * 60) + startSec) * 1000) + startMsec

		endHr, _ := strconv.ParseUint(timeSubstr[1][1], 10, 64)
		endMin, _ := strconv.ParseUint(timeSubstr[1][2], 10, 64)
		endSec, _ := strconv.ParseUint(timeSubstr[1][3], 10, 64)
		endMsec, _ := strconv.ParseUint(timeSubstr[1][4], 10, 64)
		end = (((endHr * 3600) + (endMin * 60) + endSec) * 1000) + endMsec
	} else if tsType == TS_MSMs {
		startMin, _ := strconv.ParseUint(timeSubstr[0][1], 10, 64)
		startSec, _ := strconv.ParseUint(timeSubstr[0][2], 10, 64)
		startMsec, _ := strconv.ParseUint(timeSubstr[0][3], 10, 64)
		start = (((startMin * 60) + startSec) * 1000) + startMsec

		endMin, _ := strconv.ParseUint(timeSubstr[1][1], 10, 64)
		endSec, _ := strconv.ParseUint(timeSubstr[1][2], 10, 64)
		endMsec, _ := strconv.ParseUint(timeSubstr[1][3], 10, 64)
		end = (((endMin * 60) + endSec) * 1000) + endMsec
	} else if tsType == TS_SMs {
		startSec, _ := strconv.ParseUint(timeSubstr[0][1], 10, 64)
		startMsec, _ := strconv.ParseUint(timeSubstr[0][2], 10, 64)
		start = (startSec * 1000) + startMsec

		endSec, _ := strconv.ParseUint(timeSubstr[1][1], 10, 64)
		endMsec, _ := strconv.ParseUint(timeSubstr[1][2], 10, 64)
		end = (endSec * 1000) + endMsec
	}
	return
}

func (scene *Scene) ProcessSubtext(line string) {
	if !SplitForMaxTokens {
		scene.Transcript = append(scene.Transcript, line)
		return
	}
	words := strings.Split(strings.Trim(line, " "), " ")
	var idx int
	for idx = 0; idx+MaxTokensOnAFrame < len(words); idx += MaxTokensOnAFrame {
		tmpTxt := strings.Join(words[idx:idx+MaxTokensOnAFrame], " ")
		scene.Transcript = append(scene.Transcript, tmpTxt)
	}
	if idx < len(words) {
		tmpTxt := strings.Join(words[idx:len(words)], " ")
		scene.Transcript = append(scene.Transcript, tmpTxt)
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
				vtt.Scenes = append(vtt.Scenes, scene)
				prevEndMilliSec = scene.EndMilliSec
				scene = Scene{StartMilliSec: 0, EndMilliSec: 0, Transcript: []string{}}
			}
			continue
		} else if RegNote.MatchString(line) {
			isNote = true
		} else if isSubtext {
			scene.ProcessSubtext(line)
		} else if RegTimestamp.MatchString(line) {
			scene.StartMilliSec, scene.EndMilliSec = getTimeRange(line)
			if scene.StartMilliSec-prevEndMilliSec > 50 {
				muteScene = Scene{StartMilliSec: prevEndMilliSec, EndMilliSec: scene.StartMilliSec}
				vtt.Scenes = append(vtt.Scenes, muteScene)
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
	totalSeconds := uint64(ms / 1000)
	totalMinutes := uint64(totalSeconds / 60)
	hours := uint64(totalMinutes / 60)
	minutes := totalMinutes - (hours * 60)
	seconds := totalSeconds - (totalMinutes * 60)
	fraction := ms - (totalSeconds * 1000)
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, fraction)
	}
	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, fraction)
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
