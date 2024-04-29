// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package prometheus_output

import (
	"context"
	"time"

	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/openconfig/gnmic/pkg/formatters"
	"github.com/openconfig/gnmic/pkg/outputs"
)

func (p *prometheusOutput) collectFromCache(ch chan<- prometheus.Metric) {
	notifications, err := p.gnmiCache.ReadAll()
	if err != nil {
		p.logger.Printf("failed to read from cache: %v", err)
		return
	}
	numNotifications := len(notifications)
	prometheusNumberOfCachedMetrics.Set(float64(numNotifications))

	p.targetsMeta.DeleteExpired()
	events := make([]*formatters.EventMsg, 0, numNotifications)
	for subName, notifs := range notifications {
		// build events without processors
		for _, notif := range notifs {
			targetName := notif.GetPrefix().GetTarget()
			var meta outputs.Meta
			if item := p.targetsMeta.Get(subName + "/" + targetName); item != nil {
				meta = item.Value()
			}
			ievents, err := formatters.ResponseToEventMsgs(
				subName,
				&gnmi.SubscribeResponse{
					Response: &gnmi.SubscribeResponse_Update{Update: notif},
				},
				meta)
			if err != nil {
				p.logger.Printf("failed to convert gNMI notifications to events: %v", err)
				return
			}
			events = append(events, ievents...)
		}
	}

	if p.cfg.CacheConfig.Debug {
		p.logger.Printf("got %d events from cache pre processors", len(events))
	}
	for _, proc := range p.evps {
		events = proc.Apply(events...)
	}
	if p.cfg.CacheConfig.Debug {
		p.logger.Printf("got %d events from cache post processors", len(events))
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.cfg.Timeout)
	defer cancel()
	now := time.Now()
	for _, ev := range events {
		for _, pm := range p.mb.MetricsFromEvent(ev, now) {
			select {
			case <-ctx.Done():
				p.logger.Printf("collection context terminated: %v", ctx.Err())
				return
			case ch <- pm:
			}
		}
	}
}
