// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package gatewayapi

import (
	"testing"

	"github.com/stretchr/testify/assert"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

func TestDeprecatedFieldsUsedInEnvoyProxy(t *testing.T) {
	hostStr := "otel-collector"

	tests := []struct {
		name        string
		envoyProxy  *egv1a1.EnvoyProxy
		expectedLen int
		expected    map[string]string
	}{
		{
			name:        "nil envoyProxy",
			envoyProxy:  nil,
			expectedLen: 0,
		},
		{
			name: "no telemetry",
			envoyProxy: &egv1a1.EnvoyProxy{
				Spec: egv1a1.EnvoyProxySpec{},
			},
			expectedLen: 0,
		},
		{
			name: "deprecated tracing customTags",
			envoyProxy: &egv1a1.EnvoyProxy{
				Spec: egv1a1.EnvoyProxySpec{
					Telemetry: &egv1a1.ProxyTelemetry{
						Tracing: &egv1a1.ProxyTracing{
							Provider: egv1a1.TracingProvider{
								Type: egv1a1.TracingProviderTypeOpenTelemetry,
							},
							Tracing: egv1a1.Tracing{
								CustomTags: map[string]egv1a1.CustomTag{
									"env": {Type: egv1a1.CustomTagTypeLiteral, Literal: &egv1a1.LiteralCustomTag{Value: "prod"}},
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"spec.telemetry.tracing.customTags": "spec.telemetry.tracing.tags",
			},
		},
		{
			name: "deprecated tracing provider host",
			envoyProxy: &egv1a1.EnvoyProxy{
				Spec: egv1a1.EnvoyProxySpec{
					Telemetry: &egv1a1.ProxyTelemetry{
						Tracing: &egv1a1.ProxyTracing{
							Provider: egv1a1.TracingProvider{
								Type: egv1a1.TracingProviderTypeOpenTelemetry,
								Host: &hostStr,
							},
						},
					},
				},
			},
			expected: map[string]string{
				"spec.telemetry.tracing.provider.host": "spec.telemetry.tracing.provider.backendRefs",
			},
		},
		{
			name: "deprecated access log openTelemetry host",
			envoyProxy: &egv1a1.EnvoyProxy{
				Spec: egv1a1.EnvoyProxySpec{
					Telemetry: &egv1a1.ProxyTelemetry{
						AccessLog: &egv1a1.ProxyAccessLog{
							Settings: []egv1a1.ProxyAccessLogSetting{
								{
									Sinks: []egv1a1.ProxyAccessLogSink{
										{
											Type:          egv1a1.ProxyAccessLogSinkTypeOpenTelemetry,
											OpenTelemetry: &egv1a1.OpenTelemetryEnvoyProxyAccessLog{Host: &hostStr},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"spec.telemetry.accessLog.settings[].sinks[].openTelemetry.host": "spec.telemetry.accessLog.settings[].sinks[].openTelemetry.backendRefs",
			},
		},
		{
			name: "deprecated access log openTelemetry resources",
			envoyProxy: &egv1a1.EnvoyProxy{
				Spec: egv1a1.EnvoyProxySpec{
					Telemetry: &egv1a1.ProxyTelemetry{
						AccessLog: &egv1a1.ProxyAccessLog{
							Settings: []egv1a1.ProxyAccessLogSetting{
								{
									Sinks: []egv1a1.ProxyAccessLogSink{
										{
											Type: egv1a1.ProxyAccessLogSinkTypeOpenTelemetry,
											OpenTelemetry: &egv1a1.OpenTelemetryEnvoyProxyAccessLog{
												Resources: map[string]string{"k8s.cluster.name": "my-cluster"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"spec.telemetry.accessLog.settings[].sinks[].openTelemetry.resources": "spec.telemetry.accessLog.settings[].sinks[].openTelemetry.resourceAttributes",
			},
		},
		{
			name: "deprecated metrics openTelemetry host",
			envoyProxy: &egv1a1.EnvoyProxy{
				Spec: egv1a1.EnvoyProxySpec{
					Telemetry: &egv1a1.ProxyTelemetry{
						Metrics: &egv1a1.ProxyMetrics{
							Sinks: []egv1a1.ProxyMetricSink{
								{
									Type: egv1a1.MetricSinkTypeOpenTelemetry,
									OpenTelemetry: &egv1a1.ProxyOpenTelemetrySink{
										Host: &hostStr,
									},
								},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"spec.telemetry.metrics.sinks[].openTelemetry.host": "spec.telemetry.metrics.sinks[].openTelemetry.backendRefs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deprecatedFieldsUsedInEnvoyProxy(tt.envoyProxy)
			if tt.expected == nil {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
