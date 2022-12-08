package dpkg

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func parseStatusFile(r io.Reader, entryFunc func(lineNo int, key string, value string) error,
	commitFunc func(lineNo int) error) error {
	scanner := bufio.NewScanner(r)
	key := ""
	value := ""
	entryLineNo := 0
	currentLineNo := 0
	for scanner.Scan() {
		line := scanner.Text()
		currentLineNo++
		if line == "" {
			if key != "" {
				err := entryFunc(entryLineNo, key, value)
				if err != nil {
					return err
				}
				key = ""
				value = ""
			}
			err := commitFunc(currentLineNo)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(line, " ") {
			if key == "" {
				return fmt.Errorf("invalid continuation at line %d", currentLineNo)
			}
			value += "\n"
			v := line[1:]
			if v != "." {
				value += v
			}
		} else {
			if key != "" {
				err := entryFunc(entryLineNo, key, value)
				if err != nil {
					return err
				}
				key = ""
				value = ""
			}
			p := strings.Index(line, ":")
			key = line[0:p]
			entryLineNo = currentLineNo
			if p+2 < len(line) {
				value = line[p+2:]
			}
		}
	}
	if key != "" {
		err := entryFunc(entryLineNo, key, value)
		if err != nil {
			return err
		}
		key = ""
		value = ""
	}
	return nil
}
