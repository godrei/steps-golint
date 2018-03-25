package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	glob "github.com/ryanuber/go-glob"
)

type config struct {
	Include string `env:"include"`
	Exclude string `env:"exclude"`
}

func listFiles(dir, include, exclude string) ([]string, error) {
	var files []string
	return files, filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		if include != "" && !glob.Glob(include, path) {
			return nil
		}

		if exclude != "" && glob.Glob(exclude, path) {
			return nil
		}

		files = append(files, path)

		return nil
	})
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(cfg)

	dir, err := os.Getwd()
	if err != nil {
		log.Errorf("Failed to get working directory: %s", err)
		os.Exit(1)
	}

	files, err := listFiles(dir, cfg.Include, cfg.Exclude)
	if err != nil {
		log.Errorf("Failed to list files: %s", err)
		os.Exit(1)
	}

	for _, file := range files {
		cmd := command.NewWithStandardOuts("golint", file)

		fmt.Println()
		log.Infof("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			log.Errorf("golint failed: %s", err)
			os.Exit(1)
		}
	}
}
