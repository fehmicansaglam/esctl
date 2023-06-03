package cmd

import "strings"

const Indentation = "  "

func trim(s string) string {
	return strings.TrimSpace(s)
}

func trimAndIndent(s string) string {
	s = trim(s)
	lines := strings.Split(s, "\n")
	indentedLines := make([]string, len(lines))
	for i, line := range lines {
		trimmed := trim(line)
		indentedLines[i] = Indentation + trimmed
	}
	return strings.Join(indentedLines, "\n")
}
