package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// really sad about this but i've tried for ages to do it properly :(
func opensslCommand(input string, key string, iv string) string {
	args := []string{"enc", "-aes-128-cbc", "-base64"}
	args = append(args, "-iv", fmt.Sprintf("%X", iv))
	args = append(args, "-K", fmt.Sprintf("%s", key))

	cmd := exec.Command("openssl", args...)
	cmd.Stdin = strings.NewReader(input)
	result, err := cmd.CombinedOutput()
	if err != nil {
		if e, ok := err.(*exec.Error); ok && e.Err == exec.ErrNotFound {
			return err.Error()
		}
		fmt.Println("cmd error:", err)
		log.Fatalf("result %q", result)
	}
	if n := len(result) - 1; result[n] == '\n' {
		result = result[:n]
	}
	return string(result)
}
