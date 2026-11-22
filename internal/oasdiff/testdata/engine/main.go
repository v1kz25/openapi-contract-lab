// Synthetic subprocess for boundary tests, never a semantic replacement for oasdiff.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

func main() {
	mode := os.Getenv("LAB_ENGINE_MODE")
	if slices.Contains(os.Args, "--version") && mode != "version-wait" {
		fmt.Println("oasdiff version 1.39.0")
		return
	}
	switch mode {
	case "stdout-at-limit":
		flood(os.Stdout, 16<<20)
	case "stderr-at-limit":
		fmt.Print(`[]`)
		flood(os.Stderr, 64<<10)
	case "wait", "version-wait":
		config := ""
		if i := slices.Index(os.Args, "--config"); i >= 0 && i+1 < len(os.Args) {
			config = os.Args[i+1]
		}
		data, _ := json.Marshal(struct {
			PID    int
			Config string
		}{os.Getpid(), config})
		if err := os.WriteFile(os.Getenv("LAB_ENGINE_READY"), data, 0600); err != nil {
			os.Exit(9)
		}
		fmt.Print(`[{"id":`)
		for {
			time.Sleep(time.Second)
		}
	case "stdout-limit", "stdout-limit-exit":
		// A valid JSON prefix, then more than the real 16 MiB capture limit.
		fmt.Print(`[{"id":"synthetic","level":3,"text":"`)
		flood(os.Stdout, (16<<20)+1)
		fmt.Print(`"}]`)
	case "stderr-limit", "stderr-limit-exit":
		fmt.Print(`[{"id":"synthetic","level":3}]`)
		flood(os.Stderr, (64<<10)+1)
	case "bad-json":
		fmt.Print(`[{"id":`)
	case "exit":
		fmt.Print(`[{"id":"synthetic","level":3}]`)
		fmt.Fprint(os.Stderr, "synthetic engine failed")
		os.Exit(7)
	default:
		fmt.Print(`[{"id":"synthetic","level":3}]`)
	}
	if mode == "stdout-limit-exit" || mode == "stderr-limit-exit" {
		os.Exit(7)
	}
}

func flood(writer io.Writer, remaining int) {
	chunk := bytes.Repeat([]byte("x"), 8192)
	for remaining > 0 {
		n := min(remaining, len(chunk))
		if _, err := writer.Write(chunk[:n]); err != nil {
			return // Deliberately let the caller choose its exit status.
		}
		remaining -= n
	}
}
