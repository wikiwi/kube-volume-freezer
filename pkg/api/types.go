package api

// Health contains current health status.
type Health struct {
	Status string `json:"status"`
}

// Error contains an API top-level error response.
type Error struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Issues  IssueList `json:"issues,omitempty"`
}

// Issue can be added to an Error.
type Issue struct {
	Reason       string `json:"reason"`
	Message      string `json:"message"`
	Domain       string `json:"domain,omitempty"`
	Location     string `json:"location,omitempty"`
	LocationType string `json:"locationType,omitempty"`
}

// IssueList is a list of issues.
type IssueList []*Issue

// VolumeList contains list of pod volumes.
type VolumeList struct {
	Items  []string `json:"items"`
	PodUID string   `json:"podUID"`
}

// Volume contains the information of a pod volume.
type Volume struct {
	Name   string `json:"name"`
	PodUID string `json:"podUID"`
}

// FreezeThawRequest is sent by the consumer to freeze/thaw a volume.
type FreezeThawRequest struct {
	Action string `json:"action" enum:"freeze|thaw" required:"true"`
}
