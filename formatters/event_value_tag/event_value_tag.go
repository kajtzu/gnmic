// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package event_value_tag

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/openconfig/gnmic/formatters"
	"github.com/openconfig/gnmic/types"
	"github.com/openconfig/gnmic/utils"
)

const (
	processorType = "event-value-tag"
	loggingPrefix = "[" + processorType + "] "
)

type valueTag struct {
	TagName   string `mapstructure:"tag-name,omitempty" json:"tag-name,omitempty"`
	ValueName string `mapstructure:"value-name,omitempty" json:"value-name,omitempty"`
	Consume   bool   `mapstructure:"consume,omitempty" json:"consume,omitempty"`
	Debug     bool   `mapstructure:"debug,omitempty" json:"debug,omitempty"`
	logger    *log.Logger
}

func init() {
	formatters.Register(processorType, func() formatters.EventProcessor {
		return &valueTag{logger: log.New(io.Discard, "", 0)}
	})
}

func (vt *valueTag) Init(cfg interface{}, opts ...formatters.Option) error {
	err := formatters.DecodeConfig(cfg, vt)
	if err != nil {
		return err
	}
	for _, opt := range opts {
		opt(vt)
	}

	if vt.logger.Writer() != io.Discard {
		b, err := json.Marshal(vt)
		if err != nil {
			vt.logger.Printf("initialized processor '%s': %+v", processorType, vt)
			return nil
		}
		vt.logger.Printf("initialized processor '%s': %s", processorType, string(b))
	}
	return nil
}

type foo struct {
	tags  map[string]string
	value interface{}
}

func (vt *valueTag) Apply(evs ...*formatters.EventMsg) []*formatters.EventMsg {
	if vt.TagName == "" {
		vt.TagName = vt.ValueName
	}
	// Look for events with ValueName
	toApply := make([]foo, 0)
	for _, ev := range evs {
		for k, v := range ev.Values {
			if vt.ValueName == k {
				toApply = append(toApply, foo{ev.Tags, v})
				if vt.Consume {
					delete(ev.Values, k)
				}
			}
		}
	}
	for _, bar := range toApply {
		for _, ev := range evs {
			if checkKeys(bar.tags, ev.Tags) {
				if _, ok := ev.Values[vt.ValueName]; !ok {
					ev.Tags[vt.TagName] = fmt.Sprint(bar.value)
				}
			}
		}
	}
	return evs
}

func (vt *valueTag) WithLogger(l *log.Logger) {
	if vt.Debug && l != nil {
		vt.logger = log.New(l.Writer(), loggingPrefix, l.Flags())
	} else if vt.Debug {
		vt.logger = log.New(os.Stderr, loggingPrefix, utils.DefaultLoggingFlags)
	}
}

func (vt *valueTag) WithTargets(tcs map[string]*types.TargetConfig) {}

func (vt *valueTag) WithActions(act map[string]map[string]interface{}) {}

func checkKeys(a map[string]string, b map[string]string) bool {
	for k, v := range a {
		if vv, ok := b[k]; ok {
			if v != vv {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (vt *valueTag) WithProcessors(procs map[string]map[string]any) {}
