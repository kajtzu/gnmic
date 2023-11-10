// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/protobuf/encoding/prototext"

	"github.com/openconfig/gnmic/pkg/api"
)

func main() {
	// create a target
	tg, err := api.NewTarget(
		api.Name("srl1"),
		api.Address("srl1:57400"),
		api.Username("admin"),
		api.Password("admin"),
		api.SkipVerify(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = tg.CreateGNMIClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer tg.Close()
	// create a gNMI SetRequest
	setReq, err := api.NewSetRequest(
		api.Update(
			api.Path("/interface[name=ethernet-1/1]"),
			api.Value(map[string]interface{}{
				"admin-state": "enable",
			}, "json_ietf")),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prototext.Format(setReq))
	// send the created gNMI SetRequest to the created target
	setResp, err := tg.Set(ctx, setReq)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prototext.Format(setResp))
}
