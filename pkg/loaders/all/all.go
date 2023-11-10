// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package all

import (
	_ "github.com/openconfig/gnmic/pkg/loaders/consul_loader"
	_ "github.com/openconfig/gnmic/pkg/loaders/docker_loader"
	_ "github.com/openconfig/gnmic/pkg/loaders/file_loader"
	_ "github.com/openconfig/gnmic/pkg/loaders/http_loader"
)
