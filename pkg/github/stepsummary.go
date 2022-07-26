package github

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	StepSummaryEnabledFlag = "github-step-summary"
	StepSummaryPathEnv     = "GITHUB_STEP_SUMMARY"
)

type StepSummary struct {
	flagSet  *flag.FlagSet
	enabled  bool
	filePath string
	file     *os.File
}

func (gss *StepSummary) ReadEnv() error {
	gss.filePath = os.Getenv(StepSummaryPathEnv)
	return nil
}

func (gss *StepSummary) Flags(fs *flag.FlagSet) {
	gss.flagSet = fs // Save this to check if flags actually provided later.
	desc := fmt.Sprintf("write a summary to $%s", StepSummaryPathEnv)
	enabledByDefault := gss.filePath != ""
	fs.BoolVar(&gss.enabled, StepSummaryEnabledFlag, enabledByDefault, desc)
}

func (gss *StepSummary) Open() (io.Writer, error) {
	if !gss.enabled {
		return nil, nil
	}
	if gss.filePath == "" {
		return nil, fmt.Errorf("%s is empty", StepSummaryPathEnv)
	}
	var err error
	if gss.file, err = openAppend(gss.filePath); err != nil {
		return nil, err
	}
	return gss.file, nil
}

func (gss *StepSummary) Close() error {
	if gss.file == nil {
		return nil
	}
	return gss.file.Close()
}

func openAppend(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModePerm)
}
