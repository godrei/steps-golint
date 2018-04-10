package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
)

func installedInPath(name string) bool {
	cmd := exec.Command("which", name)
	outBytes, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(outBytes)) != ""
}

func failf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	packages := os.Getenv("packages")

	log.Infof("Configs:")
	log.Printf("- exclude: %s", packages)

	if packages == "" {
		failf("Required input not defined: packages")
	}

	if !installedInPath("golint") {
		cmd := command.New("go", "get", "-u", "golang.org/x/lint/golint")

		log.Infof("\nInstalling golint")
		log.Donef("$ %s", cmd.PrintableCommandArgs())

		if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			failf("Failed to install golint: %s", out)
		}
	}

	log.Infof("\nRunning golint...")

	for _, p := range strings.Split(packages, "\n") {
		cmd := command.NewWithStandardOuts("golint", "-set_exit_status", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("golint failed: %s", err)
		}
	}
}
