package parse

import (
	"fmt"
	"github.com/sheitm/geoform-live/contracts"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	regexPatternElapsedTime = `\d{1,2}:\d{2}:\d{2}`
)

var regexElapsedTime *regexp.Regexp = regexp.MustCompile(regexPatternElapsedTime)

type eventTableParser struct {
}

func (p *eventTableParser) parse(s string) ([]*contracts.RawResult, error) {
	var res []*contracts.RawResult

	lines := getRawResultLines(s)
	for _, line := range lines {
		m := getWords(line)
		r, err := getResult(m)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func getResult(m map[int]string) (*contracts.RawResult, error) {
	rr := &contracts.RawResult{}
	if m[0] == "DSQ" {
		rr.Disqualified = true
	} else {
		rr.Disqualified = false
		p, err := strconv.Atoi(m[0])
		if err != nil {
			return nil, err
		}
		rr.Placement = p
	}
	rr.Athlete = m[1]
	offset := 2
	if len(m) == 6 {
		rr.Club = m[2]
		offset++
	}

	d, err := getDuration(m[offset])
	if err != nil {
		return nil, err
	}
	rr.ElapsedTime = d

	p, err := strconv.ParseFloat(m[offset+2], 64)
	if err != nil {
		return nil, err
	}
	rr.Points = p

	return rr, nil
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

func getWordIndexes(s string) []int {
	var indexes []int
	for i := 0; i < len(s); i++ {
		if string(s[i]) != " " {
			indexes = append(indexes, i)
			i = spool(s, i+1)
			if i < 0 {
				return indexes
			}
		}
	}
	return indexes
}

func spool(s string, start int) int {
	if start >= len(s) {
		return -1
	}
	for i := start; i < len(s) - 1; i++ {
		si := string(s[i])
		sii := string(s[i+1])
		if si == " " && sii == " " {
			return i
		}
	}
	return -1
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