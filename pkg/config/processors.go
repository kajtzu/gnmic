// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"

	"github.com/openconfig/gnmic/pkg/formatters"
)

func (c *Config) GetEventProcessors() (map[string]map[string]interface{}, error) {
	eps := c.FileConfig.GetStringMap("processors")
	for name, epc := range eps {
		switch epc := epc.(type) {
		case map[string]interface{}:
			c.logger.Printf("validating processor %q config", name)
			err := c.validateProcessorConfig(epc)
			if err != nil {
				return nil, err
			}
			c.Processors[name] = epc
		case nil:
			return nil, fmt.Errorf("empty processor %q config", name)
		default:
			c.logger.Printf("malformed processors config, %+v", epc)
			return nil, fmt.Errorf("malformed processors config, got %T", epc)
		}
	}
	for n, es := range c.Processors {
		for nn, p := range es {
			es[nn] = convert(p)
		}
		c.Processors[n] = es
	}
	for n := range c.Processors {
		expandMapEnv(c.Processors[n], "expression", "condition")
	}
	if c.Debug {
		c.logger.Printf("processors: %+v", c.Processors)
	}
	return c.Processors, nil
}

func (c *Config) validateProcessorConfig(pcfg map[string]interface{}) error {
	for epType := range pcfg {
		if !strInlist(epType, formatters.EventProcessorTypes) {
			return fmt.Errorf("unknown processors type: %s", epType)
		}
	}
	return nil
}

func strInlist(s string, ls []string) bool {
	for _, ss := range ls {
		if ss == s {
			return true
		}
	}
	return false
}
