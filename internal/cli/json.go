package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf16"
)

// writeJSON emits ASCII-only JSON so legacy shells cannot misdecode UTF-8
// before ConvertFrom-Json runs. JSON escapes preserve the original Unicode.
func writeJSON(writer io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	output := make([]byte, 0, len(data)+1)
	// In valid serialized JSON, non-ASCII runes occur only inside strings.
	for _, r := range string(data) {
		switch {
		case r <= 0x7f:
			output = append(output, byte(r))
		case r <= 0xffff:
			output = fmt.Appendf(output, `\u%04x`, r)
		default:
			high, low := utf16.EncodeRune(r)
			output = fmt.Appendf(output, `\u%04x\u%04x`, high, low)
		}
	}
	output = append(output, '\n')
	n, err := writer.Write(output)
	if err == nil && n != len(output) {
		return io.ErrShortWrite
	}
	return err
}
