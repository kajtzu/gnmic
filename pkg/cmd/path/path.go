// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package path

import (
	"github.com/openconfig/gnmic/pkg/app"
	"github.com/spf13/cobra"
)

// New creates the path command tree.
func New(gApp *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "generate gnmi or xpath style from yang file",
		Annotations: map[string]string{
			"--file": "YANG",
			"--dir":  "DIR",
		},
		PreRunE: gApp.PathPreRunE,
		RunE:    gApp.PathRunE,
		PostRun: func(cmd *cobra.Command, _ []string) {
			cmd.ResetFlags()
			gApp.InitPathFlags(cmd)
		},
		SilenceUsage: true,
	}
	gApp.InitPathFlags(cmd)
	return cmd
}
