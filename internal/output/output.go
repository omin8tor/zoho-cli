package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func JSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func JSONRaw(data []byte) error {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		_, err2 := os.Stdout.Write(data)
		return err2
	}
	buf.WriteByte('\n')
	_, err := os.Stdout.Write(buf.Bytes())
	return err
}

func Err(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
