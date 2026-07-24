// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

// This file contains code derived from Contour,
// https://github.com/projectcontour/contour
// from the source file
// https://github.com/projectcontour/contour/blob/main/internal/status/gatewayclass.go
// and is provided here subject to the following:
// Copyright Project Contour Authors
// SPDX-License-Identifier: Apache-2.0

package status

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	ReasonOlderGatewayClassExists gwapiv1.GatewayClassConditionReason = "OlderGatewayClassExists"

	// GatewayClassConditionWarning is a condition type added to GatewayClass status when the
	// referenced EnvoyProxy uses deprecated fields that need to be migrated.
	GatewayClassConditionWarning gwapiv1.GatewayClassConditionType = "Warning"

	// GatewayClassReasonDeprecatedField is the reason used with GatewayClassConditionWarning
	// when deprecated fields are detected in the referenced EnvoyProxy.
	GatewayClassReasonDeprecatedField gwapiv1.GatewayClassConditionReason = "DeprecatedField"

	MsgOlderGatewayClassExists   = "Invalid GatewayClass: another older GatewayClass with the same Spec.Controller exists"
	MsgValidGatewayClass         = "Valid GatewayClass"
	MsgGatewayClassInvalidParams = "Invalid parametersRef"
)

// SetGatewayClassAccepted inserts or updates the Accepted condition
// for the provided GatewayClass.
func SetGatewayClassAccepted(gc *gwapiv1.GatewayClass, accepted bool, reason, msg string) {
	gc.Status.Conditions = MergeConditions(gc.Status.Conditions, computeGatewayClassAcceptedCondition(gc, accepted, reason, msg))
	// Disable SupportedFeatures until the field moves from experimental to stable to avoid
	// status failures due to changes in the datatype. This can occur because we cannot control
	// how a CRD is installed in the cluster
	// gc.Status.SupportedFeatures = GatewaySupportedFeatures
}

// SetGatewayClassDeprecatedFields inserts or updates a Warning condition on the GatewayClass
// when the referenced EnvoyProxy uses deprecated fields.
// It is a no-op when gc is nil or deprecatedFields is empty.
func SetGatewayClassDeprecatedFields(gc *gwapiv1.GatewayClass, deprecatedFields map[string]string) {
	if gc == nil || len(deprecatedFields) == 0 {
		return
	}
	msg := buildEnvoyProxyDeprecationMessage(deprecatedFields)
	cond := metav1.Condition{
		Type:               string(GatewayClassConditionWarning),
		Status:             metav1.ConditionTrue,
		Reason:             string(GatewayClassReasonDeprecatedField),
		Message:            truncateConditionMessage(msg),
		ObservedGeneration: gc.Generation,
	}
	gc.Status.Conditions = MergeConditions(gc.Status.Conditions, cond)
}

// computeGatewayClassAcceptedCondition computes the GatewayClass Accepted status condition.
func computeGatewayClassAcceptedCondition(gatewayClass *gwapiv1.GatewayClass,
	accepted bool,
	reason, msg string,
) metav1.Condition {
	switch accepted {
	case true:
		return metav1.Condition{
			Type:               string(gwapiv1.GatewayClassConditionStatusAccepted),
			Status:             metav1.ConditionTrue,
			Reason:             reason,
			Message:            msg,
			ObservedGeneration: gatewayClass.Generation,
		}
	default:
		return metav1.Condition{
			Type:               string(gwapiv1.GatewayClassConditionStatusAccepted),
			Status:             metav1.ConditionFalse,
			Reason:             reason,
			Message:            msg,
			ObservedGeneration: gatewayClass.Generation,
		}
	}
}
