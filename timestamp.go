package gowebvtt

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	TS_HMSMs uint8 = iota
	TS_MSMs
	TS_SMs
	TS_Err
)

var (
	RegTimeHMSMs = regexp.MustCompile("\\s*([0-9]+)\\:([0-9]+)\\:([0-9]+)\\.([0-9]+)\\s*")
	RegTimeMSMs  = regexp.MustCompile("\\s*([0-9]+)\\:([0-9]+)\\.([0-9]+)\\s*")
	RegTimeSMs   = regexp.MustCompile("\\s*([0-9]+)\\.([0-9]+)\\s*")
)

func getTimeRangeTokens(txt string) (uint8, [][]string) {
	if RegTimeHMSMs.MatchString(txt) {
		return TS_HMSMs, RegTimeHMSMs.FindAllStringSubmatch(txt, -1)
	} else if RegTimeMSMs.MatchString(txt) {
		return TS_MSMs, RegTimeMSMs.FindAllStringSubmatch(txt, -1)
	} else if RegTimeSMs.MatchString(txt) {
		return TS_SMs, RegTimeSMs.FindAllStringSubmatch(txt, -1)
	}
	return TS_Err, nil
}

func getTimeRange(txt string) (start, end uint64, err error) {
	tsType, timeSubstr := getTimeRangeTokens(txt)
	if tsType == TS_Err || len(timeSubstr) != 2 {
		err = errors.New(fmt.Sprintf("Failed to get time range for: %v", txt))
	} else if tsType == TS_HMSMs {
		start, err = parseHMSMsToMillisec(timeSubstr[0])
		end, err = parseHMSMsToMillisec(timeSubstr[1])
	} else if tsType == TS_MSMs {
		start, err = parseMSMsToMillisec(timeSubstr[0])
		end, err = parseMSMsToMillisec(timeSubstr[1])
	} else if tsType == TS_SMs {
		start, err = parseSMsToMillisec(timeSubstr[0])
		end, err = parseSMsToMillisec(timeSubstr[1])
	}
	return
}

func parseHMSMsToMillisec(h_m_s_ms []string) (uint64, error) {
	Hr, err := strconv.ParseUint(h_m_s_ms[1], 10, 64)
	if err != nil {
		return uint64(0), err
	}
	m_s_ms, errMSMs := parseMSMsToMillisec(h_m_s_ms[1:len(h_m_s_ms)])
	if errMSMs != nil {
		return uint64(0), errMSMs
	}
	return (Hr * 3600000) + m_s_ms, nil
}

func parseMSMsToMillisec(m_s_ms []string) (uint64, error) {
	Min, err := strconv.ParseUint(m_s_ms[1], 10, 64)
	if err != nil {
		return uint64(0), err
	}
	s_ms, errSMs := parseSMsToMillisec(m_s_ms[1:len(m_s_ms)])
	if errSMs != nil {
		return uint64(0), errSMs
	}
	return (Min * 60000) + s_ms, nil
}

func parseSMsToMillisec(s_ms []string) (uint64, error) {
	Sec, errSec := strconv.ParseUint(s_ms[1], 10, 64)
	if errSec != nil {
		return uint64(0), errSec
	}
	Msec, errMs := strconv.ParseUint(s_ms[2], 10, 64)
	if errMs != nil {
		return uint64(0), errMs
	}
	return (Sec * 1000) + Msec, nil
}
