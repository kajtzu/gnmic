// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"testing"

	"github.com/AlekSi/pointer"

	"github.com/openconfig/gnmic/pkg/api/types"
)

type input struct {
	opts   []TargetOption
	config *types.TargetConfig
}

var targetTestSet = map[string]input{
	"address": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			Insecure(true),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400",
			Insecure:   pointer.ToBool(true),
			SkipVerify: pointer.ToBool(false),
			Timeout:    DefaultTargetTimeout,
		},
	},
	"username": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			Username("admin"),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400",
			Username:   pointer.ToString("admin"),
			Insecure:   pointer.ToBool(false),
			SkipVerify: pointer.ToBool(false),
			Timeout:    DefaultTargetTimeout,
		},
	},
	"two_addresses": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			Address("10.0.0.2:57400"),
			Insecure(true),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400,10.0.0.2:57400",
			Insecure:   pointer.ToBool(true),
			SkipVerify: pointer.ToBool(false),
			Timeout:    DefaultTargetTimeout,
		},
	},
	"skip_verify": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			SkipVerify(true),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400",
			Insecure:   pointer.ToBool(false),
			SkipVerify: pointer.ToBool(true),
			Timeout:    DefaultTargetTimeout,
		},
	},
	"tlsca": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			TLSCA("tlsca_path"),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400",
			Insecure:   pointer.ToBool(false),
			SkipVerify: pointer.ToBool(false),
			Timeout:    DefaultTargetTimeout,
			TLSCA:      pointer.ToString("tlsca_path"),
		},
	},
	"tls_key_cert": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			TLSKey("tlskey_path"),
			TLSCert("tlscert_path"),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400",
			Insecure:   pointer.ToBool(false),
			SkipVerify: pointer.ToBool(false),
			Timeout:    DefaultTargetTimeout,
			TLSKey:     pointer.ToString("tlskey_path"),
			TLSCert:    pointer.ToString("tlscert_path"),
		},
	},
	"token": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			Token("token_value"),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400",
			Insecure:   pointer.ToBool(false),
			SkipVerify: pointer.ToBool(false),
			Timeout:    DefaultTargetTimeout,
			Token:      pointer.ToString("token_value"),
		},
	},
	"gzip": {
		opts: []TargetOption{
			Address("10.0.0.1:57400"),
			Gzip(true),
		},
		config: &types.TargetConfig{
			Name:       "10.0.0.1:57400",
			Address:    "10.0.0.1:57400",
			Insecure:   pointer.ToBool(false),
			SkipVerify: pointer.ToBool(false),
			Timeout:    DefaultTargetTimeout,
			Gzip:       pointer.ToBool(true),
		},
	},
}

func TestNewTarget(t *testing.T) {
	for name, item := range targetTestSet {
		t.Run(name, func(t *testing.T) {
			tg, err := NewTarget(item.opts...)
			if err != nil {
				t.Errorf("failed at %q: %v", name, err)
				t.Fail()
			}
			if tg.Config.String() != item.config.String() {
				t.Errorf("failed at %q", name)
				t.Errorf("expected %+v", item.config)
				t.Errorf("     got %+v", tg.Config)
				t.Fail()
			}
		})
	}
}
