// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package app

import "fmt"

func (a *App) targetLockKey(s string) string {
	if a.Config.Clustering == nil {
		return s
	}
	if s == "" {
		return s
	}
	return fmt.Sprintf("gnmic/%s/targets/%s", a.Config.Clustering.ClusterName, s)
}
