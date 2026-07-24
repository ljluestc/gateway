// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package gatewayapi

import egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"

// deprecatedFieldsUsedInEnvoyProxy returns a map of deprecated EnvoyProxy field paths
// to their alternatives.
func deprecatedFieldsUsedInEnvoyProxy(envoyProxy *egv1a1.EnvoyProxy) map[string]string {
	deprecatedFields := make(map[string]string)
	if envoyProxy == nil || envoyProxy.Spec.Telemetry == nil {
		return deprecatedFields
	}

	telemetry := envoyProxy.Spec.Telemetry

	if tracing := telemetry.Tracing; tracing != nil {
		if len(tracing.CustomTags) > 0 {
			deprecatedFields["spec.telemetry.tracing.customTags"] = "spec.telemetry.tracing.tags"
		}
		if tracing.Provider.BackendRef != nil {
			deprecatedFields["spec.telemetry.tracing.provider.backendRef"] = "spec.telemetry.tracing.provider.backendRefs"
		}
		if tracing.Provider.Host != nil {
			deprecatedFields["spec.telemetry.tracing.provider.host"] = "spec.telemetry.tracing.provider.backendRefs"
		}
	}

	if accessLog := telemetry.AccessLog; accessLog != nil {
		for _, setting := range accessLog.Settings {
			for _, sink := range setting.Sinks {
				if sink.ALS != nil && sink.ALS.BackendRef != nil {
					deprecatedFields["spec.telemetry.accessLog.settings[].sinks[].als.backendRef"] = "spec.telemetry.accessLog.settings[].sinks[].als.backendRefs"
				}
				if openTelemetry := sink.OpenTelemetry; openTelemetry != nil {
					if openTelemetry.BackendRef != nil {
						deprecatedFields["spec.telemetry.accessLog.settings[].sinks[].openTelemetry.backendRef"] = "spec.telemetry.accessLog.settings[].sinks[].openTelemetry.backendRefs"
					}
					if openTelemetry.Host != nil {
						deprecatedFields["spec.telemetry.accessLog.settings[].sinks[].openTelemetry.host"] = "spec.telemetry.accessLog.settings[].sinks[].openTelemetry.backendRefs"
					}
					if len(openTelemetry.Resources) > 0 {
						deprecatedFields["spec.telemetry.accessLog.settings[].sinks[].openTelemetry.resources"] = "spec.telemetry.accessLog.settings[].sinks[].openTelemetry.resourceAttributes"
					}
				}
			}
		}
	}

	if metrics := telemetry.Metrics; metrics != nil {
		for _, sink := range metrics.Sinks {
			if openTelemetry := sink.OpenTelemetry; openTelemetry != nil {
				if openTelemetry.BackendRef != nil {
					deprecatedFields["spec.telemetry.metrics.sinks[].openTelemetry.backendRef"] = "spec.telemetry.metrics.sinks[].openTelemetry.backendRefs"
				}
				if openTelemetry.Host != nil {
					deprecatedFields["spec.telemetry.metrics.sinks[].openTelemetry.host"] = "spec.telemetry.metrics.sinks[].openTelemetry.backendRefs"
				}
			}
		}
	}

	return deprecatedFields
}
