// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (a *App) GeneratePathPreRunE(cmd *cobra.Command, args []string) error {
	a.Config.SetLocalFlagsFromFile(cmd)
	if a.Config.GeneratePathSearch && a.Config.GeneratePathWithDescr {
		return errors.New("flags --search and --descr cannot be used together")
	}
	if a.Config.LocalFlags.GeneratePathPathType != "xpath" && a.Config.LocalFlags.GeneratePathPathType != "gnmi" {
		return errors.New("path-type must be one of 'xpath' or 'gnmi'")
	}
	return nil
}

func (a *App) GeneratePathRunE(cmd *cobra.Command, args []string) error {
	return a.PathCmdRun(
		a.Config.GlobalFlags.Dir,
		a.Config.GlobalFlags.File,
		a.Config.GlobalFlags.Exclude,
		pathGenOpts{
			search:        a.Config.LocalFlags.GeneratePathSearch,
			withDescr:     a.Config.LocalFlags.GeneratePathWithDescr,
			withTypes:     a.Config.LocalFlags.GeneratePathWithTypes,
			withPrefix:    a.Config.LocalFlags.GeneratePathWithPrefix,
			pathType:      a.Config.LocalFlags.GeneratePathPathType,
			stateOnly:     a.Config.LocalFlags.GeneratePathState,
			configOnly:    a.Config.LocalFlags.GeneratePathConfig,
			json:          a.Config.LocalFlags.GenerateJSON,
			withNonLeaves: a.Config.LocalFlags.GeneratePathWithNonLeaves,
		},
	)
}

func (a *App) InitGeneratePathFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.Flags().StringVarP(&a.Config.LocalFlags.GeneratePathPathType, "path-type", "", "xpath", "path type xpath or gnmi")
	cmd.Flags().BoolVarP(&a.Config.LocalFlags.GeneratePathWithDescr, "descr", "", false, "print leaf description")
	cmd.Flags().BoolVarP(&a.Config.LocalFlags.GeneratePathWithPrefix, "with-prefix", "", false, "include module/submodule prefix in path elements")
	cmd.Flags().BoolVarP(&a.Config.LocalFlags.GeneratePathWithTypes, "types", "", false, "print leaf type")
	cmd.Flags().BoolVarP(&a.Config.LocalFlags.GeneratePathSearch, "search", "", false, "search through path list")
	cmd.Flags().BoolVarP(&a.Config.LocalFlags.GeneratePathState, "state-only", "", false, "generate paths only for YANG leafs representing state data")
	cmd.Flags().BoolVarP(&a.Config.LocalFlags.GeneratePathConfig, "config-only", "", false, "generate paths only for YANG leafs representing config data")
	cmd.Flags().BoolVarP(&a.Config.LocalFlags.GeneratePathWithNonLeaves, "with-non-leaves", "", false, "also generate paths for non-leaf nodes")
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		a.Config.FileConfig.BindPFlag(fmt.Sprintf("%s-%s", cmd.Name(), flag.Name), flag)
	})
}
