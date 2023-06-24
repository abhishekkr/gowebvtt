package gowebvtt

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ParseFile(filepath string) VTT {
	readFile, err := os.Open(filepath)
	defer readFile.Close()
	if err != nil {
		log.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	return ParseWebVTT(fileScanner, VttOptions{})
}

func ParseFileWithOptions(filepath string, opts VttOptions) (VTT, error) {
	readFile, err := os.Open(filepath)
	defer readFile.Close()
	if err != nil {
		return VTT{}, err
	}
	fileScanner := bufio.NewScanner(readFile)

	return ParseWebVTT(fileScanner, opts), nil
}

func ParseWebVTT(fileScanner *bufio.Scanner, opts VttOptions) VTT {
	var vtt = VTT{Scenes: []Scene{}, Options: opts}

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
				scene = newScene(0, 0, vtt.Options.SceneOptions)
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
				muteScene = newScene(prevEndMilliSec, scene.StartMilliSec, vtt.Options.SceneOptions)
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

func String(vtt VTT) string {
	return vtt.String()
}

func Println(vtt VTT) {
	for _, s := range vtt.Scenes {
		fmt.Println("\nStarts:", s.StartMilliSec, ",\tEnds:", s.EndMilliSec)
		for _, w := range s.Transcript {
			fmt.Println(w)
		}
	}
}
