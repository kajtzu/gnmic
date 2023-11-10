// © 2022 Nokia.
//
// This code is a Contribution to the gNMIc project (“Work”) made under the Google Software Grant and Corporate Contributor License Agreement (“CLA”) and governed by the Apache License 2.0.
// No other rights or licenses in or to any of Nokia’s intellectual property are granted for any other purpose.
// This code is provided on an “as is” basis without any warranties of any kind.
//
// SPDX-License-Identifier: Apache-2.0

package nats_output

import "github.com/prometheus/client_golang/prometheus"

var NatsNumberOfSentMsgs = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "gnmic",
	Subsystem: "nats_output",
	Name:      "number_of_nats_msgs_sent_success_total",
	Help:      "Number of msgs successfully sent by gnmic nats output",
}, []string{"publisher_id", "subject"})

var NatsNumberOfSentBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "gnmic",
	Subsystem: "nats_output",
	Name:      "number_of_written_nats_bytes_total",
	Help:      "Number of bytes written by gnmic nats output",
}, []string{"publisher_id", "subject"})

var NatsNumberOfFailSendMsgs = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "gnmic",
	Subsystem: "nats_output",
	Name:      "number_of_nats_msgs_sent_fail_total",
	Help:      "Number of failed msgs sent by gnmic nats output",
}, []string{"publisher_id", "reason"})

var NatsSendDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "gnmic",
	Subsystem: "nats_output",
	Name:      "msg_send_duration_ns",
	Help:      "gnmic nats output send duration in ns",
}, []string{"publisher_id"})

func initMetrics() {
	NatsNumberOfSentMsgs.WithLabelValues("", "").Add(0)
	NatsNumberOfSentBytes.WithLabelValues("", "").Add(0)
	NatsNumberOfFailSendMsgs.WithLabelValues("", "").Add(0)
	NatsSendDuration.WithLabelValues("").Set(0)
}

func registerMetrics(reg *prometheus.Registry) error {
	initMetrics()
	var err error
	if err = reg.Register(NatsNumberOfSentMsgs); err != nil {
		return err
	}
	if err = reg.Register(NatsNumberOfSentBytes); err != nil {
		return err
	}
	if err = reg.Register(NatsNumberOfFailSendMsgs); err != nil {
		return err
	}
	if err = reg.Register(NatsSendDuration); err != nil {
		return err
	}
	return nil
}
