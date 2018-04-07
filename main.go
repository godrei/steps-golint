package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/godrei/steps-golint/gotool"
)

// Config ...
type Config struct {
	Exclude string `env:"exclude"`
}

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
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Error: %s\n", err)
	}
	stepconf.Print(cfg)

	if !installedInPath("golint") {
		cmd := command.New("go", "get", "-u", "golang.org/x/lint/golint")

		log.Infof("\nInstalling golint")
		log.Donef("$ %s", cmd.PrintableCommandArgs())

		if out, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			failf("Failed to install golint: %s", out)
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		failf("Failed to get working directory: %s", err)
	}

	excludes := strings.Split(cfg.Exclude, ",")

	packages, err := gotool.ListPackages(dir, excludes...)
	if err != nil {
		failf("Failed to list packages: %s", err)
	}

	log.Infof("\nRunning golint...")

	for _, p := range packages {
		cmd := command.NewWithStandardOuts("golint", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("golint failed: %s", err)
		}
	}
}
