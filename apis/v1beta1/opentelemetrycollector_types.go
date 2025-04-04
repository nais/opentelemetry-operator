// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +kubebuilder:skip

package v1beta1

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/open-telemetry/opentelemetry-operator/apis/v1alpha1"
)

// Probe defines the OpenTelemetry's pod probe config. Only Liveness probe is supported currently.
type Probe struct {
	// Number of seconds after the container has started before liveness probes are initiated.
	// Defaults to 0 seconds. Minimum value is 0.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds *int32 `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +optional
	PeriodSeconds *int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	// +optional
	SuccessThreshold *int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	FailureThreshold *int32 `json:"failureThreshold,omitempty"`
	// Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
	// The grace period is the duration in seconds after the processes running in the pod are sent
	// a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
	// value overrides the value provided by the pod spec.
	// Value must be non-negative integer. The value zero indicates stop immediately via
	// the kill signal (no opportunity to shut down).
	// This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
	// Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
}

// Ingress is used to specify how OpenTelemetry Collector is exposed. This
// functionality is only available if one of the valid modes is set.
// Valid modes are: deployment, daemonset and statefulset.
// NOTE: If this feature is activated, all specified receivers are exposed.
// Currently, this has a few limitations. Depending on the ingress controller
// there are problems with TLS and gRPC.
// SEE: https://github.com/open-telemetry/opentelemetry-operator/issues/1306.
// NOTE: As a workaround, port name and appProtocol could be specified directly
// in the CR.
// SEE: OpenTelemetryCollector.spec.ports[index].
type Ingress struct {
	// Type default value is: ""
	// Supported types are: ingress, route
	Type IngressType `json:"type,omitempty"`

	// RuleType defines how Ingress exposes collector receivers.
	// IngressRuleTypePath ("path") exposes each receiver port on a unique path on single domain defined in Hostname.
	// IngressRuleTypeSubdomain ("subdomain") exposes each receiver port on a unique subdomain of Hostname.
	// Default is IngressRuleTypePath ("path").
	RuleType IngressRuleType `json:"ruleType,omitempty"`

	// Hostname by which the ingress proxy can be reached.
	// +optional
	Hostname string `json:"hostname,omitempty"`

	// Annotations to add to ingress.
	// e.g. 'cert-manager.io/cluster-issuer: "letsencrypt"'
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// TLS configuration.
	// +optional
	TLS []networkingv1.IngressTLS `json:"tls,omitempty"`

	// IngressClassName is the name of an IngressClass cluster resource. Ingress
	// controller implementations use this field to know whether they should be
	// serving this Ingress resource.
	// +optional
	IngressClassName *string `json:"ingressClassName,omitempty"`

	// Route is an OpenShift specific section that is only considered when
	// type "route" is used.
	// +optional
	Route OpenShiftRoute `json:"route,omitempty"`
}

// OpenShiftRoute defines openshift route specific settings.
type OpenShiftRoute struct {
	// Termination indicates termination type. By default "edge" is used.
	Termination TLSRouteTerminationType `json:"termination,omitempty"`
}

// OpenTelemetryCollectorSpec defines the desired state of OpenTelemetryCollector.
type OpenTelemetryCollectorSpec struct {
	// OpenTelemetryCommonFields are fields that are on all OpenTelemetry CRD workloads.
	OpenTelemetryCommonFields `json:",inline"`
	// TargetAllocator indicates a value which determines whether to spawn a target allocation resource or not.
	// +optional
	TargetAllocator TargetAllocatorEmbedded `json:"targetAllocator,omitempty"`
	// Mode represents how the collector should be deployed (deployment, daemonset, statefulset or sidecar)
	// +optional
	Mode Mode `json:"mode,omitempty"`
	// UpgradeStrategy represents how the operator will handle upgrades to the CR when a newer version of the operator is deployed
	// +optional
	UpgradeStrategy UpgradeStrategy `json:"upgradeStrategy"`
	// Config is the raw JSON to be used as the collector's configuration. Refer to the OpenTelemetry Collector documentation for details.
	// The empty objects e.g. batch: should be written as batch: {} otherwise they won't work with kustomize or kubectl edit.
	// +required
	// +kubebuilder:pruning:PreserveUnknownFields
	Config Config `json:"config"`
	// Ingress is used to specify how OpenTelemetry Collector is exposed. This
	// functionality is only available if one of the valid modes is set.
	// Valid modes are: deployment, daemonset and statefulset.
	// +optional
	Ingress Ingress `json:"ingress,omitempty"`
	// Liveness config for the OpenTelemetry Collector except the probe handler which is auto generated from the health extension of the collector.
	// It is only effective when healthcheckextension is configured in the OpenTelemetry Collector pipeline.
	// +optional
	LivenessProbe *Probe `json:"livenessProbe,omitempty"`

	// ObservabilitySpec defines how telemetry data gets handled.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Observability"
	Observability v1alpha1.ObservabilitySpec `json:"observability,omitempty"`

	// ConfigMaps is a list of ConfigMaps in the same namespace as the OpenTelemetryCollector
	// object, which shall be mounted into the Collector Pods.
	// Each ConfigMap will be added to the Collector's Deployments as a volume named `configmap-<configmap-name>`.
	ConfigMaps []v1alpha1.ConfigMapsSpec `json:"configmaps,omitempty"`
	// UpdateStrategy represents the strategy the operator will take replacing existing DaemonSet pods with new pods
	// https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/daemon-set-v1/#DaemonSetSpec
	// This is only applicable to Daemonset mode.
	// +optional
	DaemonSetUpdateStrategy appsv1.DaemonSetUpdateStrategy `json:"daemonSetUpdateStrategy,omitempty"`
	// UpdateStrategy represents the strategy the operator will take replacing existing Deployment pods with new pods
	// https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/deployment-v1/#DeploymentSpec
	// This is only applicable to Deployment mode.
	// +optional
	DeploymentUpdateStrategy appsv1.DeploymentStrategy `json:"deploymentUpdateStrategy,omitempty"`
}

// TargetAllocatorEmbedded defines the configuration for the Prometheus target allocator, embedded in the
// OpenTelemetryCollector spec.
type TargetAllocatorEmbedded struct {
	// Replicas is the number of pod instances for the underlying TargetAllocator. This should only be set to a value
	// other than 1 if a strategy that allows for high availability is chosen. Currently, the only allocation strategy
	// that can be run in a high availability mode is consistent-hashing.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// NodeSelector to schedule OpenTelemetry TargetAllocator pods.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Resources to set on the OpenTelemetryTargetAllocator containers.
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// AllocationStrategy determines which strategy the target allocator should use for allocation.
	// The current options are least-weighted, consistent-hashing and per-node. The default is
	// consistent-hashing.
	// +optional
	// +kubebuilder:default:=consistent-hashing
	AllocationStrategy TargetAllocatorAllocationStrategy `json:"allocationStrategy,omitempty"`
	// FilterStrategy determines how to filter targets before allocating them among the collectors.
	// The only current option is relabel-config (drops targets based on prom relabel_config).
	// The default is relabel-config.
	// +optional
	// +kubebuilder:default:=relabel-config
	FilterStrategy TargetAllocatorFilterStrategy `json:"filterStrategy,omitempty"`
	// ServiceAccount indicates the name of an existing service account to use with this instance. When set,
	// the operator will not automatically create a ServiceAccount for the TargetAllocator.
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`
	// Image indicates the container image to use for the OpenTelemetry TargetAllocator.
	// +optional
	Image string `json:"image,omitempty"`
	// Enabled indicates whether to use a target allocation mechanism for Prometheus targets or not.
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// If specified, indicates the pod's scheduling constraints
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// PrometheusCR defines the configuration for the retrieval of PrometheusOperator CRDs ( servicemonitor.monitoring.coreos.com/v1 and podmonitor.monitoring.coreos.com/v1 )  retrieval.
	// All CR instances which the ServiceAccount has access to will be retrieved. This includes other namespaces.
	// +optional
	PrometheusCR TargetAllocatorPrometheusCR `json:"prometheusCR,omitempty"`
	// SecurityContext configures the container security context for
	// the targetallocator.
	// +optional
	SecurityContext *v1.SecurityContext `json:"securityContext,omitempty"`
	// PodSecurityContext configures the pod security context for the
	// targetallocator.
	// +optional
	PodSecurityContext *v1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	// TopologySpreadConstraints embedded kubernetes pod configuration option,
	// controls how pods are spread across your cluster among failure-domains
	// such as regions, zones, nodes, and other user-defined topology domains
	// https://kubernetes.io/docs/concepts/workloads/pods/pod-topology-spread-constraints/
	// +optional
	TopologySpreadConstraints []v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// Toleration embedded kubernetes pod configuration option,
	// controls how pods can be scheduled with matching taints
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// ENV vars to set on the OpenTelemetry TargetAllocator's Pods. These can then in certain cases be
	// consumed in the config file for the TargetAllocator.
	// +optional
	Env []v1.EnvVar `json:"env,omitempty"`
	// ObservabilitySpec defines how telemetry data gets handled.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Observability"
	Observability v1alpha1.ObservabilitySpec `json:"observability,omitempty"`
	// PodDisruptionBudget specifies the pod disruption budget configuration to use
	// for the target allocator workload.
	//
	// +optional
	PodDisruptionBudget *PodDisruptionBudgetSpec `json:"podDisruptionBudget,omitempty"`
}

// OpenTelemetryCollectorStatus defines the observed state of OpenTelemetryCollector.
type OpenTelemetryCollectorStatus struct {
	// Scale is the OpenTelemetryCollector's scale subresource status.
	// +optional
	Scale v1alpha1.ScaleSubresourceStatus `json:"scale,omitempty"`

	// Version of the managed OpenTelemetry Collector (operand)
	// +optional
	Version string `json:"version,omitempty"`

	// Image indicates the container image to use for the OpenTelemetry Collector.
	// +optional
	Image string `json:"image,omitempty"`

	// Messages about actions performed by the operator on this resource.
	// +optional
	// +listType=atomic
	// Deprecated: use Kubernetes events instead.
	Messages []string `json:"messages,omitempty"`

	// Replicas is currently not being set and might be removed in the next version.
	// +optional
	// Deprecated: use "OpenTelemetryCollector.Status.Scale.Replicas" instead.
	Replicas int32 `json:"replicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OpenTelemetryCollector is the Schema for the opentelemetrycollectors API.
type OpenTelemetryCollector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenTelemetryCollectorSpec   `json:"spec,omitempty"`
	Status OpenTelemetryCollectorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenTelemetryCollectorList contains a list of OpenTelemetryCollector.
type OpenTelemetryCollectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenTelemetryCollector `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenTelemetryCollector{}, &OpenTelemetryCollectorList{})
}
