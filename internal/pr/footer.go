package pr

import "strings"

const (
	StartMarker = "<!-- gts:stack-start -->"
	EndMarker   = "<!-- gts:stack-end -->"
)

// ParseBlock returns the content between markers (excluding markers, including its trailing newline).
// ok=false if markers are absent or malformed.
func ParseBlock(body string) (string, bool) {
	si := strings.Index(body, StartMarker)
	if si < 0 {
		return "", false
	}
	ei := strings.Index(body, EndMarker)
	if ei < 0 || ei <= si {
		return "", false
	}
	// Content begins after StartMarker + its trailing newline.
	content := body[si+len(StartMarker) : ei]
	content = strings.TrimPrefix(content, "\n")
	return content, true
}

// ReplaceOrAppend writes block between markers. If markers are absent, append with a blank line.
// `block` should NOT include the markers themselves and should NOT include trailing newlines.
func ReplaceOrAppend(body, block string) string {
	wrapped := StartMarker + "\n" + block + "\n" + EndMarker

	si := strings.Index(body, StartMarker)
	ei := strings.Index(body, EndMarker)
	if si >= 0 && ei > si {
		// Replace from start marker through end marker (inclusive).
		end := ei + len(EndMarker)
		return body[:si] + wrapped + body[end:]
	}

	// Append. Ensure exactly one blank line before block.
	trimmed := strings.TrimRight(body, "\n")
	if trimmed == "" {
		return wrapped + "\n"
	}
	return trimmed + "\n\n" + wrapped + "\n"
}
