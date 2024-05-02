package api

import "time"

type HPA struct {
	APIVersion string     `json:"apiVersion" yaml:"apiVersion"`
	Kind       string     `json:"kind" yaml:"kind"`
	Metadata   ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec       HPASpec    `json:"spec" yaml:"spec"`
	Status     HPAStatus  `json:"status" yaml:"status"`
}

type HPASpec struct {
	Behavior       HPABehavior                 `json:"behavior" yaml:"behavior"`
	MaxReplica     int                         `json:"maxReplica" yaml:"maxReplica"`
	MinReplica     int                         `json:"minReplica" yaml:"minReplica"`
	Metrics        []MetricsSpec                 `json:"metrics" yaml:"metrics"`
	ScaleTargetRef CrossVersionObjectReference `json:"scaleTargetRef" yaml:"scaleTargetRef"`
}

type HPAStatus struct {
	Conditions      []HPACondition `json:"conditions" yaml:"conditions"`
	CurrentReplicas int          `json:"currentReplicas" yaml:"currentReplicas"`
	DesiredReplicas int          `json:"desiredReplicas" yaml:"desiredReplicas"`
	LastScaleTime   time.Time    `json:"lastScaleTime" yaml:"lastScaleTime"`
}

type HPABehavior struct {
	ScaleDown HPAScalingRules `json:"scaleDown" yaml:"scaleDown"`
	ScaleUp   HPAScalingRules `json:"scaleUp" yaml:"scaleUp"`
}

type MetricsSpec struct {
	Type     string               `json:"type" yaml:"type"`
	Resource ResourceMetricSource `json:"resource" yaml:"resource"`
}

type CrossVersionObjectReference struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Name       string `json:"name" yaml:"name"`
}

type HPAScalingRules struct {
	Policies     HPAScalingPolicy `json:"policies" yaml:"policies"`
	SelectPolicy string           `json:"selectPolicy" yaml:"selectPolicy"`
}

type HPAScalingPolicy struct {
	Type          string `json:"type" yaml:"type"`
	Value         int    `json:"value" yaml:"value"`
	PeriodSeconds int    `json:"periodSeconds" yaml:"periodSeconds"`
}

type ResourceMetricSource struct {
	Name   string       `json:"name" yaml:"name"`
	Target MetricTarget `json:"target" yaml:"target"`
}

type MetricTarget struct {
	Type               string `json:"type" yaml:"type"`
	AverageUtilization int    `json:"averageUtilization" yaml:"averageUtilization"`
}

type HPACondition struct {
	LastTransitionTime time.Time `json:"lastTransitionTime" yaml:"lastTransitionTime"`
	Message            string    `json:"message" yaml:"message"`
	Reason             string    `json:"reason" yaml:"reason"`
	Status             string    `json:"status" yaml:"status"`
	Type               string    `json:"type" yaml:"type"`
}
