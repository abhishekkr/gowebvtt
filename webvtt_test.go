package gowebvtt

import (
	"strings"
	"testing"
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

func TestSplitAndAppendSceneTwoLines(t *testing.T) {
	vttOpts := VttOptions{Enabled: true, MaxLinesPerScene: 2}
	vtt := VTT{Options: vttOpts}
	sceneA := newScene(0, 5000, vttOpts.SceneOptions)
	sceneA.Transcript = []string{
		"this is",
		"a test",
		"of multiple",
		"lines in",
		"a scene",
	}
	errSplitA := vtt.splitAndAppendScene(sceneA)
	if errSplitA != nil {
		t.Fatalf("Failed to splitAndAppendScene for scene: %v", sceneA)
	}
	if len(vtt.Scenes) != 3 {
		t.Fatalf("Failed to splitAndAppendScene in correct count for scene: %v", sceneA)
	}

	sceneA0 := vtt.Scenes[0]
	if sceneA0.StartMilliSec != 0 || sceneA0.EndMilliSec != 1666 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.0: %v", sceneA0)
	}
	if sceneA0.Transcript[0] != "this is" || sceneA0.Transcript[1] != "a test" {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.0: %v", sceneA0)
	}

	sceneA1 := vtt.Scenes[1]
	if sceneA1.StartMilliSec != 1666 || sceneA1.EndMilliSec != 3332 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.1: %v", sceneA0)
	}
	if sceneA1.Transcript[0] != "of multiple" || sceneA1.Transcript[1] != "lines in" {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.1: %v", sceneA0)
	}

	sceneA2 := vtt.Scenes[2]
	if sceneA2.StartMilliSec != 3332 || sceneA2.EndMilliSec != 5000 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.2: %v", sceneA0)
	}
	if sceneA2.Transcript[0] != "a scene" {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.2: %v", sceneA0)
	}
}

func TestSplitAndAppendSceneThreeLines(t *testing.T) {
	vttOpts := VttOptions{Enabled: true, MaxLinesPerScene: 3}
	vtt := VTT{Options: vttOpts}
	sceneA := newScene(0, 5000, vttOpts.SceneOptions)
	sceneA.Transcript = []string{
		"this is",
		"a test",
		"of multiple",
		"lines in",
		"a scene",
	}
	errSplitA := vtt.splitAndAppendScene(sceneA)
	if errSplitA != nil {
		t.Fatalf("Failed to splitAndAppendScene for scene: %v", sceneA)
	}
	if len(vtt.Scenes) != 2 {
		t.Fatalf("Failed to splitAndAppendScene in correct count for scene: %v", sceneA)
	}

	sceneA0 := vtt.Scenes[0]
	if sceneA0.StartMilliSec != 0 || sceneA0.EndMilliSec != 2500 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.0: %v", sceneA0)
	}
	if sceneA0.Transcript[0] != "this is" || sceneA0.Transcript[1] != "a test" || sceneA0.Transcript[2] != "of multiple" {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.0: %v", sceneA0)
	}

	sceneA1 := vtt.Scenes[1]
	if sceneA1.StartMilliSec != 2500 || sceneA1.EndMilliSec != 5000 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.1: %v", sceneA1)
	}
	if sceneA1.Transcript[0] != "lines in" || sceneA1.Transcript[1] != "a scene" {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.1: %v", sceneA1)
	}
}

func TestSplitAndAppendSceneFourLines(t *testing.T) {
	vttOpts := VttOptions{Enabled: true, MaxLinesPerScene: 4}
	vtt := VTT{Options: vttOpts}
	sceneA := newScene(0, 5000, vttOpts.SceneOptions)
	sceneA.Transcript = []string{
		"this is",
		"a test",
		"of multiple",
		"lines in",
		"a scene",
	}
	errSplitA := vtt.splitAndAppendScene(sceneA)
	if errSplitA != nil {
		t.Fatalf("Failed to splitAndAppendScene for scene: %v", sceneA)
	}
	if len(vtt.Scenes) != 2 {
		t.Fatalf("Failed to splitAndAppendScene in correct count for scene: %v", sceneA)
	}

	sceneA0 := vtt.Scenes[0]
	if sceneA0.StartMilliSec != 0 || sceneA0.EndMilliSec != 2500 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.0: %v", sceneA0)
	}
	if sceneA0.Transcript[0] != "this is" || sceneA0.Transcript[1] != "a test" || sceneA0.Transcript[2] != "of multiple" || sceneA0.Transcript[3] != "lines in" {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.0: %v", sceneA0)
	}

	sceneA1 := vtt.Scenes[1]
	if sceneA1.StartMilliSec != 2500 || sceneA1.EndMilliSec != 5000 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.1: %v", sceneA1)
	}
	if sceneA1.Transcript[0] != "a scene" {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.1: %v", sceneA1)
	}
}

func TestSplitAndAppendSceneFiveLines(t *testing.T) {
	vttOpts := VttOptions{Enabled: true, MaxLinesPerScene: 5}
	vtt := VTT{Options: vttOpts}
	sceneA := newScene(0, 5000, vttOpts.SceneOptions)
	sceneA.Transcript = []string{
		"this is",
		"a test",
		"of multiple",
		"lines in",
		"a scene",
	}
	errSplitA := vtt.splitAndAppendScene(sceneA)
	if errSplitA != nil {
		t.Fatalf("Failed to splitAndAppendScene for scene: %v", sceneA)
	}
	if len(vtt.Scenes) != 1 {
		t.Fatalf("Failed to splitAndAppendScene in correct count for scene: %v", sceneA)
	}

	sceneA0 := vtt.Scenes[0]
	if sceneA0.StartMilliSec != 0 || sceneA0.EndMilliSec != 5000 {
		t.Fatalf("Failed to splitAndAppendScene get correct time slot for scene.0: %v", sceneA0)
	}
	if len(sceneA0.Transcript) != 5 {
		t.Fatalf("Failed to splitAndAppendScene get correct transcript for scene.0: %v", sceneA0)
	}
}
