package api

type PV struct {
	APIVersion string     `yaml:"apiVersion" json:"apiVersion"`
	Kind       string     `yaml:"kind" json:"kind"`
	Metadata   ObjectMeta `yaml:"metadata" json:"metadata"`
	Spec       PVSpec     `yaml:"spec" json:"spec"`
	Status     PVStatus   `yaml:"status" json:"status"`
}

type PVSpec struct {
	AccessModes []string        `yaml:"accessModes" json:"accessModes"`
	Capacity    CapacitySpec    `yaml:"capacity" json:"capacity"`
	NFS         NFSVolumeSource `yaml:"nfs" json:"nfs"`
}

type CapacitySpec struct {
	Storage string `yaml:"storage" json:"storage"`
}

type NFSVolumeSource struct {
	Path     string `json:"path" yaml:"path"`
	ReadOnly bool   `json:"readOnly" yaml:"readOnly"`
	Server   string `json:"server" yaml:"server"`
}

type PVStatus struct {
	Message string `yaml:"message" json:"message"`
	Phase   string `yaml:"phase" json:"phase"`
	Reason  string `yaml:"reason" json:"reason"`
}
