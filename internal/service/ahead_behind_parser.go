package service

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseAheadBehind(out string) (behind int, ahead int, err error) {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != 2 {
		return 0, 0, fmt.Errorf("invalid ahead/behind output %q", out)
	}

	behind, err = strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid behind count %q: %w", fields[0], err)
	}

	ahead, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid ahead count %q: %w", fields[1], err)
	}

	return behind, ahead, nil
}
