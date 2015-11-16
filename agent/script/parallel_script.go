package script

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type ParallelScript struct {
	name       string
	allScripts []Script

	logTag string
	logger boshlog.Logger
}

type scriptResult struct {
	Script Script
	Error  error
}

func NewParallelScript(name string, scripts []Script, logger boshlog.Logger) ParallelScript {
	return ParallelScript{
		name:       name,
		allScripts: scripts,

		logTag: "ParallelScript",
		logger: logger,
	}
}

func (s ParallelScript) Tag() string  { return "" }
func (s ParallelScript) Path() string { return "" }
func (s ParallelScript) Exists() bool { return true }

func (s ParallelScript) RunAsync() (boshsys.Process, boshsys.File, boshsys.File, error) {
	return nil, nil, nil, bosherr.Error("RunAsync not supported for ParallelScript")
}

func (s ParallelScript) Run() error {
	existingScripts := s.findExistingScripts(s.allScripts)

	s.logger.Info(s.logTag, "Will run %d %s scripts in parallel", len(existingScripts), s.name)

	resultsChan := make(chan scriptResult)

	for _, script := range existingScripts {
		script := script
		go func() { resultsChan <- scriptResult{script, script.Run()} }()
	}

	var failedScripts, passedScripts []string

	for i := 0; i < len(existingScripts); i++ {
		select {
		case r := <-resultsChan:
			jobName := r.Script.Tag()

			if r.Error == nil {
				passedScripts = append(passedScripts, jobName)
				s.logger.Info(s.logTag, "'%s' script has successfully executed", r.Script.Path())
			} else {
				failedScripts = append(failedScripts, jobName)
				s.logger.Error(s.logTag, "'%s' script has failed with error: %s", r.Script.Path(), r.Error)
			}
		}
	}

	return s.summarizeErrs(passedScripts, failedScripts)
}

func (s ParallelScript) findExistingScripts(all []Script) []Script {
	var existing []Script

	for _, script := range all {
		if script.Exists() {
			s.logger.Debug(s.logTag, "Found '%s' script in job '%s'", s.name, script.Tag())
			existing = append(existing, script)
		} else {
			s.logger.Debug(s.logTag, "Did not find '%s' script in job '%s'", s.name, script.Tag())
		}
	}

	return existing
}

func (s ParallelScript) summarizeErrs(passedScripts, failedScripts []string) error {
	if len(failedScripts) > 0 {
		errMsg := "Failed Jobs: " + strings.Join(failedScripts, ", ")

		if len(passedScripts) > 0 {
			errMsg += ". Successful Jobs: " + strings.Join(passedScripts, ", ")
		}

		totalRan := len(passedScripts) + len(failedScripts)

		return bosherr.Errorf("%d of %d %s scripts failed. %s.", len(failedScripts), totalRan, s.name, errMsg)
	}

	return nil
}
