package types

// UpdateStatus represents the update status for a single container image
type ImageUpdateStatus struct {
	ContainerID        string `json:"containerId,omitempty" yaml:"containerId" csv:"container_id"`
	ContainerName      string `json:"containerName,omitempty" yaml:"containerName" csv:"container_name"`
	Image              Image  `json:"image" yaml:"image" csv:"image"`
	OriginalTag        string `json:"originalTag,omitempty" yaml:"originalTag" csv:"original_tag"`
	LatestAvailableTag string `json:"latestAvailableTag,omitempty" yaml:"latestAvailable_tag" csv:"latest_available_tag"`
	UpdateAvailable    bool   `json:"updateAvailable,omitempty" yaml:"updateAvailable" csv:"update_available"`
	StatusMessage      string `json:"statusMessage,omitempty" yaml:"statusMessage" csv:"status_message"`
	Error              string `json:"error,omitempty" yaml:"error" csv:"error"`
}
