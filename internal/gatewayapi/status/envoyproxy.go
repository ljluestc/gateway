// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package status

import (
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

func UpdateEnvoyProxyStatusAccepted(ep *egv1a1.EnvoyProxy, ancestor *gwapiv1.ParentReference,
	reason egv1a1.EnvoyProxyConditionReason, msg string,
) {
	if ep == nil || ancestor == nil {
		return
	}
	status := metav1.ConditionTrue
	if reason != egv1a1.EnvoyProxyReasonAccepted {
		status = metav1.ConditionFalse
	}

	cond := newCondition(string(egv1a1.EnvoyProxyConditionAccepted), status,
		string(reason), msg, ep.Generation)

	for i := range ep.Status.Ancestors {
		item := ep.Status.Ancestors[i]
		if ancestorRefsEqual(&item.AncestorRef, ancestor) {
			ep.Status.Ancestors[i].Conditions = MergeConditions(item.Conditions, cond)
			return
		}
	}

	// ancestor not found, append a new one
	ep.Status.Ancestors = append(ep.Status.Ancestors, egv1a1.EnvoyProxyAncestorStatus{
		AncestorRef: *ancestor,
		Conditions: []metav1.Condition{
			cond,
		},
	})
}

// UpdateEnvoyProxyStatusDeprecatedFields sets a Warning condition on the EnvoyProxy
// status for any deprecated fields that are in use, identified by deprecatedFields
// (a map of deprecated field path → recommended alternative).
// It is a no-op when deprecatedFields is empty or ep/ancestor are nil.
func UpdateEnvoyProxyStatusDeprecatedFields(
	ep *egv1a1.EnvoyProxy,
	ancestor *gwapiv1.ParentReference,
	deprecatedFields map[string]string,
) {
	if ep == nil || ancestor == nil || len(deprecatedFields) == 0 {
		return
	}

	msg := buildEnvoyProxyDeprecationMessage(deprecatedFields)
	cond := newCondition(
		string(egv1a1.EnvoyProxyConditionWarning),
		metav1.ConditionTrue,
		string(egv1a1.EnvoyProxyReasonDeprecatedField),
		msg,
		ep.Generation,
	)

	for i := range ep.Status.Ancestors {
		item := ep.Status.Ancestors[i]
		if ancestorRefsEqual(&item.AncestorRef, ancestor) {
			ep.Status.Ancestors[i].Conditions = MergeConditions(item.Conditions, cond)
			return
		}
	}

	// ancestor not found, append a new one
	ep.Status.Ancestors = append(ep.Status.Ancestors, egv1a1.EnvoyProxyAncestorStatus{
		AncestorRef: *ancestor,
		Conditions:  []metav1.Condition{cond},
	})
}

// buildEnvoyProxyDeprecationMessage builds a human-readable warning message from a map
// of deprecated field paths to their recommended alternatives.
func buildEnvoyProxyDeprecationMessage(deprecatedFields map[string]string) string {
	keys := make([]string, 0, len(deprecatedFields))
	for k := range deprecatedFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, field := range keys {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(field)
		b.WriteString(" is deprecated, use ")
		b.WriteString(deprecatedFields[field])
		b.WriteString(" instead")
	}
	return b.String()
}
