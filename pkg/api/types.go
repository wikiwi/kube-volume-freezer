package api

// Health contains current health status.
type Health struct {
	Status string `json:"status"`
}

// Error contains an API top-level error response.
type Error struct {
	// Code corresponds to HTTP error codes.
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

// IssueList is a list of Issues.
type IssueList []*Issue

// VolumeList contains list of Pod Volumes.
type VolumeList struct {
	Items  []string `json:"items"`
	PodUID string   `json:"podUID"`
}

// Volume contains the information of a Pod Volume.
type Volume struct {
	Name   string `json:"name"`
	PodUID string `json:"podUID"`
}

// FreezeThawRequest describes a request for freezing or thawing a Volume.
type FreezeThawRequest struct {
	Action string `json:"action" enum:"freeze|thaw" required:"true"`
}
