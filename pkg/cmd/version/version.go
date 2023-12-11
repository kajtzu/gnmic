// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"github.com/openconfig/gnmic/pkg/app"
	"github.com/spf13/cobra"
)

// New creates the version command tree.
func New(gApp *app.App) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "show gnmic version",
		PreRun: func(cmd *cobra.Command, _ []string) {
			gApp.Config.SetLocalFlagsFromFile(cmd)
		},
		Run: gApp.VersionRun,
	}
	versionCmd.AddCommand(newVersionUpgradeCmd(gApp))

	return versionCmd
}
