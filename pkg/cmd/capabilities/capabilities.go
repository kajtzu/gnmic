// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package capabilities

import (
	"github.com/openconfig/gnmic/pkg/app"
	"github.com/spf13/cobra"
)

// capabilitiesCmd represents the capabilities command
func New(gApp *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "capabilities",
		Aliases:      []string{"cap"},
		Short:        "query targets gnmi capabilities",
		PreRunE:      gApp.CapPreRunE,
		RunE:         gApp.CapRunE,
		SilenceUsage: true,
	}
	gApp.InitCapabilitiesFlags(cmd)
	return cmd
}
