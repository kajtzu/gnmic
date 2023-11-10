// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package script_action

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/openconfig/gnmic/pkg/actions"
)

const (
	loggingPrefix = "[script_action] "
	actionType    = "script"
	defaultShell  = "/bin/bash"
)

func init() {
	actions.Register(actionType, func() actions.Action {
		return &scriptAction{
			logger: log.New(io.Discard, "", 0),
		}
	})
}

type scriptAction struct {
	Name    string `mapstructure:"name,omitempty"`
	Shell   string `mapstructure:"shell,omitempty"`
	Command string `mapstructure:"command,omitempty"`
	File    string `mapstructure:"file,omitempty"`
	Debug   bool   `mapstructure:"debug,omitempty"`

	logger *log.Logger
}

func (s *scriptAction) Init(cfg map[string]interface{}, opts ...actions.Option) error {
	err := actions.DecodeConfig(cfg, s)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		opt(s)
	}
	if s.Name == "" {
		return fmt.Errorf("action type %q missing name field", actionType)
	}
	err = s.setDefaults()
	if err != nil {
		return err
	}
	s.logger.Printf("action name %q of type %q initialized: %v", s.Name, actionType, s)
	return nil
}

func (s *scriptAction) Run(_ context.Context, aCtx *actions.Context) (interface{}, error) {
	if s.Command == "" && s.File == "" {
		return nil, nil
	}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	var cmd *exec.Cmd
	if s.Command != "" {
		cmds := strings.Split(s.Command, "\n")
		args := append([]string{"-c"}, strings.Join(cmds, "; "))
		cmd = exec.Command(s.Shell, args...)
	}
	if s.File != "" {
		cmd = exec.Command(s.File)
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = os.Environ()
	for k, v := range aCtx.Env {
		k = strings.ReplaceAll(k, "-", "_")
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range aCtx.Vars {
		k = strings.ReplaceAll(k, "-", "_")
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}
	if stderr.String() != "" {
		return map[string]string{"stderr": stderr.String()}, nil
	}
	return map[string]string{"stdout": stdout.String()}, nil
}

func (s *scriptAction) NName() string { return s.Name }

func (s *scriptAction) setDefaults() error {
	if s.Shell == "" {
		s.Shell = defaultShell
	}
	return nil
}
