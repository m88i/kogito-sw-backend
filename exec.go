package main

import (
	"bytes"
	"errors"
	"os/exec"
)

func execCmd(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New(stderr.String())
	}
	return nil
}
