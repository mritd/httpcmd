package main

import (
	"bytes"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func try(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	w := logrus.StandardLogger().Writer()
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func run(command string, args ...string) (string, string, error) {
	cmd := exec.Command(command, args...)
	var sout bytes.Buffer
	var serr bytes.Buffer
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	//cmd.CombinedOutput()
	err := cmd.Run()
	if err != nil {
		return sout.String(), serr.String(), err
	} else {
		return sout.String(), serr.String(), nil
	}
}
