package github

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	StepSumarryEnabledFlag = "github-step-summary"
	StepSummaryPathEnv     = "GITHUB_STEP_SUMMARY"
)

type StepSummary struct {
	flagSet  *flag.FlagSet
	enabled  bool
	filePath string
}

func (gss *StepSummary) ReadEnv() error {
	gss.filePath = os.Getenv(StepSumarryEnabledFlag)
	return nil
}

func (gss *StepSummary) Flags(fs *flag.FlagSet) {
	gss.flagSet = fs // Save this to check if flags actually provided later.
	desc := fmt.Sprintf("write a summary to $%s", StepSummaryPathEnv)
	enabledByDefault := gss.filePath != ""
	fs.BoolVar(&gss.enabled, StepSumarryEnabledFlag, enabledByDefault, desc)
}

func (gss *StepSummary) Write(s string, echo io.Writer) error {
	if !gss.enabled {
		return nil
	}
	var (
		w          io.Writer
		closeError error
	)
	if gss.filePath == "" {
		log.Printf("warning: %s is empty, so not writing to that file.", StepSummaryPathEnv)
		w = echo
	} else {
		file, err := os.OpenFile(gss.filePath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer func() { closeError = file.Close() }()
		log.Printf("Writing step summary to %s", gss.filePath)
		w = io.MultiWriter(file, echo)
	}
	_, err := w.Write([]byte(s))
	if err != nil {
		return err
	}

	return closeError
}
