package utils

import "strings"

const Indentation = "  "

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func TrimAndIndent(s string) string {
	s = Trim(s)
	lines := strings.Split(s, "\n")
	indentedLines := make([]string, len(lines))
	for i, line := range lines {
		trimmed := Trim(line)
		indentedLines[i] = Indentation + trimmed
	}
	return strings.Join(indentedLines, "\n")
}
