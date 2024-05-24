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
	MaxReplica     int             `json:"maxReplica" yaml:"maxReplica"`
	MinReplica     int             `json:"minReplica" yaml:"minReplica"`
	Metrics        MetricsSpec     `json:"metrics" yaml:"metrics"`
	Template       PodTemplateSpec `json:"template" yaml:"template"`
	Selector       LabelSelector   `json:"selector" yaml:"selector"`
	AdjustInterval int             `json:"adjustInterval" yaml:"adjustInterval"`
}

type HPAStatus struct {
	Conditions      []HPACondition `json:"conditions" yaml:"conditions"`
	CurrentReplicas int            `json:"currentReplicas" yaml:"currentReplicas"`
	DesiredReplicas int            `json:"desiredReplicas" yaml:"desiredReplicas"`
	LastScaleTime   time.Time      `json:"lastScaleTime" yaml:"lastScaleTime"`
}

type MetricsSpec struct {
	CPUPercentage    float64 `json:"cpuPercentage" yaml:"cpuPercentage"`
	MemoryPercentage float64 `json:"memoryPercentage" yaml:"memoryPercentage"`
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
