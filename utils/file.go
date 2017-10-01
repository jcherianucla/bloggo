package utils

import (
	"bufio"
	"strings"
)

const (
	DELIM = "x---x"
)

// Go through markdown file and grab the meta data for post
func ReadMDMeta(scan *bufio.Scanner) (map[string]string, error) {
	meta := make(map[string]string)
	started := false
	for scan.Scan() {
		line := strings.Trim(scan.Text(), " ")
		if line == DELIM {
			if started {
				return meta, nil
			} else {
				started = true
			}
		} else if started {
			sections := strings.SplitN(line, ":", 2)
			if len(sections) < 2 {
				return nil, MDMissingSections
			}
			meta[sections[0]] = strings.Trim(sections[1], " ")
		}
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return nil, MDEmpty
}
