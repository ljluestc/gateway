// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package gatewayapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

func TestDeprecatedFieldsUsedInBackendTrafficPolicy(t *testing.T) {
	rateLimitType := egv1a1.GlobalRateLimitType
	headerName := "x-user"

	tests := []struct {
		name     string
		policy   *egv1a1.BackendTrafficPolicy
		expected map[string]string
	}{
		{
			name: "no deprecated fields",
			policy: &egv1a1.BackendTrafficPolicy{
				Spec: egv1a1.BackendTrafficPolicySpec{},
			},
			expected: map[string]string{},
		},
		{
			name: "deprecated rateLimit.type field",
			policy: &egv1a1.BackendTrafficPolicy{
				Spec: egv1a1.BackendTrafficPolicySpec{
					RateLimit: &egv1a1.RateLimitSpec{
						Type: &rateLimitType,
					},
				},
			},
			expected: map[string]string{
				"spec.rateLimit.type": "spec.rateLimit.global and/or spec.rateLimit.local",
			},
		},
		{
			name: "deprecated consistentHash.header field",
			policy: &egv1a1.BackendTrafficPolicy{
				Spec: egv1a1.BackendTrafficPolicySpec{
					ClusterSettings: egv1a1.ClusterSettings{
						LoadBalancer: &egv1a1.LoadBalancer{
							Type: egv1a1.ConsistentHashLoadBalancerType,
							ConsistentHash: &egv1a1.ConsistentHash{
								Type:   egv1a1.HeaderConsistentHashType,
								Header: &egv1a1.Header{Name: headerName},
							},
						},
					},
				},
			},
			expected: map[string]string{
				"spec.loadBalancer.consistentHash.header": "spec.loadBalancer.consistentHash.headers",
			},
		},
		{
			name: "deprecated telemetry tracing customTags field",
			policy: &egv1a1.BackendTrafficPolicy{
				Spec: egv1a1.BackendTrafficPolicySpec{
					Telemetry: &egv1a1.BackendTelemetry{
						Tracing: &egv1a1.Tracing{
							CustomTags: map[string]egv1a1.CustomTag{
								"env": {Type: egv1a1.CustomTagTypeLiteral, Literal: &egv1a1.LiteralCustomTag{Value: "prod"}},
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
			name: "multiple deprecated fields",
			policy: &egv1a1.BackendTrafficPolicy{
				Spec: egv1a1.BackendTrafficPolicySpec{
					RateLimit: &egv1a1.RateLimitSpec{
						Type: &rateLimitType,
					},
					Compression: []*egv1a1.Compression{
						{Type: egv1a1.GzipCompressorType},
					},
				},
			},
			expected: map[string]string{
				"spec.rateLimit.type": "spec.rateLimit.global and/or spec.rateLimit.local",
				"spec.compression":    "spec.compressor",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure targetRefs are valid so TargetRef check is not triggered.
			tt.policy.Spec.TargetRefs = []gwapiv1.LocalPolicyTargetReferenceWithSectionName{
				{
					LocalPolicyTargetReference: gwapiv1.LocalPolicyTargetReference{
						Group: "gateway.networking.k8s.io",
						Kind:  "Gateway",
						Name:  "eg",
					},
				},
			}
			got := deprecatedFieldsUsedInBackendTrafficPolicy(tt.policy)
			assert.Equal(t, tt.expected, got)
		})
	}
}
