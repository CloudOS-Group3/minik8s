package api

type StorageClass struct {
	APIVersion    string     `yaml:"apiVersion" json:"api_version"`
	Kind          string     `yaml:"kind" json:"kind"`
	Metadata      ObjectMeta `yaml:"metadata" json:"metadata"`
	MountOptions  []string   `yaml:"mountOptions" json:"mount_options"`
	Parameters    struct{}   `yaml:"parameters" json:"parameters"`
	Provisioner   string     `yaml:"provisioner" json:"provisioner"`
	ReclaimPolicy string     `yaml:"reclaimPolicy" json:"reclaim_policy"`
}
