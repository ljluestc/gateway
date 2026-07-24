// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

func TestUpdateEnvoyProxyStatusDeprecatedFields(t *testing.T) {
	ancestor := &gwapiv1.ParentReference{
		Name: gwapiv1.ObjectName("envoy-gateway"),
	}

	tests := []struct {
		name             string
		deprecatedFields map[string]string
		expectConditions int
		expectType       string
		expectReason     string
		expectMessage    string
	}{
		{
			name:             "no deprecated fields - no condition set",
			deprecatedFields: map[string]string{},
			expectConditions: 0,
		},
		{
			name: "single deprecated field",
			deprecatedFields: map[string]string{
				"spec.telemetry.tracing.customTags": "spec.telemetry.tracing.tags",
			},
			expectConditions: 1,
			expectType:       string(egv1a1.EnvoyProxyConditionWarning),
			expectReason:     string(egv1a1.EnvoyProxyReasonDeprecatedField),
			expectMessage:    "spec.telemetry.tracing.customTags is deprecated, use spec.telemetry.tracing.tags instead",
		},
		{
			name: "multiple deprecated fields - deterministic order",
			deprecatedFields: map[string]string{
				"spec.telemetry.tracing.provider.host":   "spec.telemetry.tracing.provider.backendRefs",
				"spec.telemetry.accessLog.sinks[].host": "spec.telemetry.accessLog.sinks[].backendRefs",
			},
			expectConditions: 1,
			expectType:       string(egv1a1.EnvoyProxyConditionWarning),
			expectReason:     string(egv1a1.EnvoyProxyReasonDeprecatedField),
			expectMessage: "spec.telemetry.accessLog.sinks[].host is deprecated, use spec.telemetry.accessLog.sinks[].backendRefs instead; " +
				"spec.telemetry.tracing.provider.host is deprecated, use spec.telemetry.tracing.provider.backendRefs instead",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &egv1a1.EnvoyProxy{
				ObjectMeta: metav1.ObjectMeta{
					Generation: 5,
				},
			}

			UpdateEnvoyProxyStatusDeprecatedFields(ep, ancestor, tt.deprecatedFields)

			if tt.expectConditions == 0 {
				assert.Empty(t, ep.Status.Ancestors, "no ancestor status should be set for empty deprecated fields")
				return
			}

			if assert.Len(t, ep.Status.Ancestors, 1) {
				cond := ep.Status.Ancestors[0].Conditions
				if assert.Len(t, cond, 1) {
					assert.Equal(t, tt.expectType, cond[0].Type)
					assert.Equal(t, metav1.ConditionTrue, cond[0].Status)
					assert.Equal(t, tt.expectReason, cond[0].Reason)
					assert.Equal(t, int64(5), cond[0].ObservedGeneration)
					assert.Equal(t, tt.expectMessage, cond[0].Message)
				}
			}
		})
	}
}

func TestUpdateEnvoyProxyStatusDeprecatedFieldsMergesIntoExistingAncestor(t *testing.T) {
	ancestor := &gwapiv1.ParentReference{
		Name: gwapiv1.ObjectName("envoy-gateway"),
	}
	ep := &egv1a1.EnvoyProxy{
		ObjectMeta: metav1.ObjectMeta{Generation: 2},
	}

	// First set the Accepted condition via the normal path.
	UpdateEnvoyProxyStatusAccepted(ep, ancestor, egv1a1.EnvoyProxyReasonAccepted, "EnvoyProxy has been accepted.")

	// Then add a deprecated-field Warning condition.
	UpdateEnvoyProxyStatusDeprecatedFields(ep, ancestor, map[string]string{
		"spec.telemetry.tracing.customTags": "spec.telemetry.tracing.tags",
	})

	// Both conditions should live under the single existing ancestor entry.
	if assert.Len(t, ep.Status.Ancestors, 1) {
		conditions := ep.Status.Ancestors[0].Conditions
		assert.Len(t, conditions, 2, "Accepted + Warning conditions expected under one ancestor")

		types := make(map[string]bool)
		for _, c := range conditions {
			types[c.Type] = true
		}
		assert.True(t, types[string(egv1a1.EnvoyProxyConditionAccepted)], "Accepted condition must be present")
		assert.True(t, types[string(egv1a1.EnvoyProxyConditionWarning)], "Warning condition must be present")
	}
}
