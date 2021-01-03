package scrape


import (
	"fmt"
	"github.com/sheitm/ofever/types"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	regexPatternElapsedTime = `\d{1,2}:\d{2}:\d{2}`
	regexPatternELapsedTimeNormal = `\d{1,2}:\d{2}:\d{2} \+  \d{1,2}:\d{2}` // Ex: 1:21:22 +  37:14
	regexPatternELapsedTimeLong = `\d{1,2}:\d{2}:\d{2} \+\d{1}:\d{2}:\d{2}` // Ex: 1:46:31 +1:02:23
	regexPatternPoints = `\d{1,3}\.\d{2}`
	regexPatternPointsComma = `\d{1,3},\d{2}`
	regexPatterMissingControls = `\(\-\d{1,2} poster\)`
	regexPatternMissingControlsDigits = `\d{1,2}`
	regexPatternClubContamination = `\d:\d{2}:\d{2}\s\+`
)

var regexElapsedTime *regexp.Regexp = regexp.MustCompile(regexPatternElapsedTime)
var regexPoints *regexp.Regexp = regexp.MustCompile(regexPatternPoints)
var regexPointsComma *regexp.Regexp = regexp.MustCompile(regexPatternPointsComma)
var regexMissingControls = regexp.MustCompile(regexPatterMissingControls)
var regexMissingControlsDigits = regexp.MustCompile(regexPatternMissingControlsDigits)
var regexClubContamination = regexp.MustCompile(regexPatternClubContamination)

type eventTableParser struct {}

func (p *eventTableParser) parse(columns, s string) ([]*types.Result, error) {
	var res []*types.Result

	lines := getRawResultLines(s)
	for _, line := range lines {
		r, err := getResultFromLine(line)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func getResultFromLine(line string) (*types.Result, error) {
	w := getWords(line)
	r := &types.Result{}


	points, err := getPoints(line)
	if err == nil {
		r.Points = points
	}

	if strings.Contains(line, "DELTATT") {
		r.Disqualified = true
		r.Athlete = w[0]
		r.Club = w[1]
		return r, nil
	}

	placement, err := strconv.Atoi(w[0])
	if err == nil {
		r.Placement = placement
	} else if w[0] == "DSQ" {
		r.Disqualified = true
		mc, err := getMissingControls(line)
		if err == nil {
			r.MissingControls = mc
		}
	}

	r.Athlete = w[1]

	if len(w[2]) > 0 {
		var ru rune
		for _, i := range w[2] {
			ru = i
			break
		}
		if !unicode.IsDigit(ru) {
			r.Club = washClub(w[2])
		}
	}

	dur, err := getDuration(line)
	if err == nil {
		r.ElapsedTime = dur
	}

	return r, nil
}

func getPoints(s string) (float64, error) {
	b := regexPoints.Find([]byte(s))
	if len(b) == 0 {
		b = regexPointsComma.Find([]byte(s))
		if len(b) == 0 {
			return 0, fmt.Errorf("not a valid point (float) format: %s", s)
		}

		ss := strings.ReplaceAll(string(b), ",", ".")
		b = []byte(ss)
	}

	p, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return 0, err
	}
	return p, nil
}

func getMissingControls(s string) (int, error) {
	b := regexMissingControls.Find([]byte(s))
	if len(b) == 0 {
		return 0, fmt.Errorf("not a valid missing controls format: %s", s)
	}
	bb := regexMissingControlsDigits.Find(b)
	i, err := strconv.Atoi(string(bb))
	if err != nil {
		return 0, err
	}
	return i, nil
}

func washClub(s string) string {
	if s == "" {
		return ""
	}

	b := regexClubContamination.Find([]byte(s))
	if len(b) == 0 {
		return s
	}

	l := len(string(b))
	return s[:len(s)-l]
}

func getDuration(s string) (time.Duration, error) {
	b := regexElapsedTime.Find([]byte(s))
	if len(b) == 0 {
		return 0, fmt.Errorf("not a valid duration format: %s", s)
	}
	arr := strings.Split(string(b), ":")
	hours, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(arr[1])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.Atoi(arr[2])
	if err != nil {
		return 0, err
	}

	h := time.Duration(hours) * time.Hour
	m := time.Duration(minutes) * time.Minute
	ss := time.Duration(seconds) * time.Second

	d := h + m + ss
	return d, nil
}

func getWords(s string) map[int]string {
	res := map[int]string{}
	var w  []rune
	inWord := false
	prevWasSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) && !inWord {
			continue
		}
		if unicode.IsSpace(r) && prevWasSpace {
			res[len(res)] = string(w)
			w = []rune{}
			inWord = false
			prevWasSpace = false
			continue
		}
		if unicode.IsSpace(r) {
			prevWasSpace = true
			continue
		}

		if !unicode.IsSpace(r) {
			if prevWasSpace {
				rr, _ := utf8.DecodeLastRuneInString(" ")
				w = append(w, rr)
			}
			inWord = true
			prevWasSpace = false
			w = append(w, r)
			continue
		}
	}

	if len(w) > 0 {
		res[len(res)] = string(w)
	}

	return res
}

func getRawResultLines(s string) []string {
	var lines []string
	rawLines := strings.Split(s, "\n")
	for _, line := range rawLines {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

