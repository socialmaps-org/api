package multiline

import "strings"

func Dedent(s string) string {
	lines := strings.Split(s, "\n")

	var result []string
	for i, line := range lines {
		if i == 0 && line == "" {
			continue
		}
		result = append(result, strings.TrimLeft(line, "\t"))
	}

	return strings.Join(result, "\n")
}
