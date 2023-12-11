// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"

	"github.com/openconfig/gnmic/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// newVersionUpgradeCmd creates the version upgrade command tree.
func newVersionUpgradeCmd(gApp *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Aliases: []string{"up"},
		Short:   "upgrade gnmic to latest available version",
		PreRun: func(cmd *cobra.Command, _ []string) {
			gApp.Config.SetLocalFlagsFromFile(cmd)
		},
		RunE: gApp.VersionUpgradeRun,
	}
	initVersionUpgradeFlags(cmd, gApp)
	return cmd
}

func initVersionUpgradeFlags(cmd *cobra.Command, gApp *app.App) {
	cmd.Flags().Bool("use-pkg", false, "upgrade using package")
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		gApp.Config.FileConfig.BindPFlag(fmt.Sprintf("%s-%s", cmd.Name(), flag.Name), flag)
	})
}
