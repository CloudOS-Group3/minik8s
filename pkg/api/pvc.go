package api

type PVC struct {
	APIVersion string     `yaml:"apiVersion" json:"apiVersion"`
	Kind       string     `yaml:"kind" json:"kind"`
	Metadata   ObjectMeta `yaml:"metadata" json:"metadata"`
	Spec       PVCSpec    `yaml:"spec" json:"spec"`
	Status     PVCStatus  `yaml:"status" json:"status"`
}

type PVCSpec struct {
	AccessModes      []string             `yaml:"accessModes" json:"accessModes"`
	VolumeMode       string               `yaml:"volumeMode" json:"volumeMode"`
	VolumeName       string               `yaml:"volumeName" json:"volumeName"`
	StorageClassName string               `yaml:"storageClassName" json:"storageClassName"`
	Selector         LabelSelector        `yaml:"selector" json:"selector"`
	Resources        ResourceRequirements `yaml:"resources" json:"resources"`
}

type PVCStatus struct {
	AccessModes        []string     `yaml:"accessModes" json:"accessModes"`
	AllocatedResources struct{}     `yaml:"allocatedResources" json:"allocatedResources"`
	Capacity           CapacitySpec `yaml:"capacity" json:"capacity"`
	TargetPV           PV           `yaml:"targetPV" json:"targetPV"`
}
