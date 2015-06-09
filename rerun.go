package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func removeTimingInfo(content []byte) []byte {
	// Match things like 02/15
	date := regexp.MustCompile("\\d+(/\\d+)+")
	// Match things like 9:52
	time := regexp.MustCompile("\\d+(:\\d+)+")
	// Match things like 2.3s
	duration := regexp.MustCompile("\\d+(\\.\\d+)[μsmhd]")

	content = date.ReplaceAll(content, nil)
	content = time.ReplaceAll(content, nil)
	content = duration.ReplaceAll(content, nil)

	return content
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: rerun [command]")
		os.Exit(1)
	}

	total := 0
	errors := 0
	passes := 0

	command := os.Args[1]
	args := strings.Join(os.Args[2:], " ")
	terminal := command + " " + args
	seen := map[string]int{}

	for {
		total++
		cmd := exec.Command(command, os.Args[2:]...)
		output, runtimeerr := cmd.CombinedOutput()

		// This is a byte array; we need to convert it to a byte slice
		hashbin := md5.Sum(removeTimingInfo(output))
		outputHash := hex.EncodeToString(hashbin[:])
		filename := fmt.Sprintf("rerun-failure-%s.log", outputHash)

		if runtimeerr == nil {
			passes++
			log.Printf("%s SUCCEEDED (total: %d, passing: %d, errors: %d)", terminal, total, passes, errors)
		} else {
			errors++

			value, ok := seen[outputHash]
			if ok {
				seen[outputHash] = value + 1
			} else {
				seen[outputHash]++
			}

			log.Printf("%s FAILED (total: %d, passing: %d, errors: %d)", terminal, total, passes, errors)
			log.Printf("Seen error %s %d times", outputHash, seen[outputHash])
			if !fileExists(filename) {
				ioerr := ioutil.WriteFile(filename, output, 0600)
				if ioerr == nil {
					log.Printf("Wrote error log to %s", filename)
				} else {
					log.Printf("Unable to write error log: %s", ioerr.Error())
				}
			}
		}
	}
}
