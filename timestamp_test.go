package gowebvtt

import (
	"testing"
)

func TestGetTimeRangeTokens(t *testing.T) {
	if tsType, _ := getTimeRangeTokens(""); tsType != TS_Err {
		t.Fatal("Failed to getTimeRangeTokens for empty line.")
	}

	hmsm, hmsmMatchs := getTimeRangeTokens("00:00:00.100")
	msm, msmMatchs := getTimeRangeTokens("00:00.100")
	sm, smMatchs := getTimeRangeTokens("00.100")
	if hmsm != TS_HMSMs || len(hmsmMatchs[0]) != 5 {
		t.Fatalf("Failed to getTimeRangeTokens for HMSMs.")
	}
	if msm != TS_MSMs || len(msmMatchs[0]) != 4 {
		t.Fatalf("Failed to getTimeRangeTokens for MSMs.")
	}
	if sm != TS_SMs || len(smMatchs[0]) != 3 {
		t.Fatalf("Failed to getTimeRangeTokens for SMs.")
	}
}

func TestGetTimeRange(t *testing.T) {
	var data = []string{
		"",
		"0:00:00.000 1:01:01:10",
	}

	for _, ts := range data {
		if _, _, err := getTimeRange(ts); err == nil {
			t.Fatalf("Failed to err on getTimeRange for: %s", ts)
		}
	}

	if start, end, err := getTimeRange("00:00:00.100 ---> 01:10:01.100"); start != 100 || end != 4201100 || err != nil {
		t.Fatalf("Failed to getTimeRange for HMSMs. start:%d end:%d err:%v", start, end, err)
	}
	if start, end, err := getTimeRange("00:10.100 --> 01:10.500"); start != 10100 || end != 70500 || err != nil {
		t.Fatalf("Failed to getTimeRange for MSMs. start:%d end:%d err:%v", start, end, err)
	}
	if start, end, err := getTimeRange("00.500 -> 10.100"); start != 500 || end != 10100 || err != nil {
		t.Fatalf("Failed to getTimeRange for SMs. start:%d end:%d err:%v", start, end, err)
	}
}

func TestParseHMSMsToMillisec(t *testing.T) {
	timeX := []string{"00:00:01.100", "00", "00", "01", "100"}
	timeXMs := uint64(1100)
	if ms, err := parseHMSMsToMillisec(timeX); ms != timeXMs || err != nil {
		t.Fatalf("Failed to parseHMSMsToMillisec for %v with %d == %d.", timeX, ms, timeXMs)
	}

	timeY := []string{"01:10:01.100", "01", "10", "01", "100"}
	timeYMs := uint64(4201100)
	if ms, err := parseHMSMsToMillisec(timeY); ms != timeYMs || err != nil {
		t.Fatalf("Failed to parseHMSMsToMillisec for %v with %d == %d.", timeY, ms, timeYMs)
	}
}

func TestParseMSMsToMillisec(t *testing.T) {
	timeX := []string{"00:01.100", "00", "01", "100"}
	timeXMs := uint64(1100)
	if ms, err := parseMSMsToMillisec(timeX); ms != timeXMs || err != nil {
		t.Fatalf("Failed to parseMSMsToMillisec for %v with %d == %d.", timeX, ms, timeXMs)
	}

	timeY := []string{"10:01.100", "10", "01", "100"}
	timeYMs := uint64(601100)
	if ms, err := parseMSMsToMillisec(timeY); ms != timeYMs || err != nil {
		t.Fatalf("Failed to parseMSMsToMillisec for %v with %d == %d.", timeY, ms, timeYMs)
	}
}

func TestParseSMsToMillisec(t *testing.T) {
	timeX := []string{"01.100", "01", "100"}
	timeXMs := uint64(1100)
	if ms, err := parseSMsToMillisec(timeX); ms != timeXMs || err != nil {
		t.Fatalf("Failed to parseSMsToMillisec for %v with %d == %d.", timeX, ms, timeXMs)
	}

	timeY := []string{"10.100", "10", "100"}
	timeYMs := uint64(10100)
	if ms, err := parseSMsToMillisec(timeY); ms != timeYMs || err != nil {
		t.Fatalf("Failed to parseSMsToMillisec for %v with %d == %d.", timeY, ms, timeYMs)
	}
}
