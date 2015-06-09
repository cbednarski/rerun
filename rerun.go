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

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: rerun [command]")
		os.Exit(1)
	}

	seen := map[string]int{}

	cmd := exec.Command(strings.Join(os.Args[1:], " "))
	for {
		output, runtimeerr := cmd.CombinedOutput()

		// This is a byte array; we need to convert it to a byte slice
		hashbin := sha256.Sum256(output)
		outputHash := hex.EncodeToString(hashbin[:])
		filename := fmt.Sprintf("rerun-failure-%s.log", outputHash)

		if runtimeerr != nil {
			value, ok := seen[outputHash]
			if ok {
				seen[outputHash] = value + 1
			} else {

			}
			seen[outputHash]++
			ioerr := ioutil.WriteFile(filename, output, 0600)
			if ioerr != nil {
				log.Printf("Unable to write error log: %s", ioerr.Error())
			}
		}
	}
}
