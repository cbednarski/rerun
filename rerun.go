package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
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

	for {
		total++
		cmd := exec.Command(command, os.Args[2:]...)
		output, runtimeerr := cmd.CombinedOutput()

		// This is a byte array; we need to convert it to a byte slice
		hashbin := sha256.Sum256(output)
		outputHash := hex.EncodeToString(hashbin[:])
		filename := fmt.Sprintf("rerun-failure-%s.log", outputHash)

		if runtimeerr == nil {
			passes++
			log.Printf("%s SUCCEEDED (total: %d, passing: %d, errors: %d)", terminal, total, passes, errors)
		} else {
			errors++
			// log.Printf(runtimeerr.Error())
			log.Printf("%s FAILED (total: %d, passing: %d, errors: %d)", terminal, total, passes, errors)
			if !fileExists(filename) {
				ioerr := ioutil.WriteFile(filename, output, 0600)
				if ioerr != nil {
					log.Printf("Unable to write error log: %s", ioerr.Error())
				}
			}
		}
	}
}
