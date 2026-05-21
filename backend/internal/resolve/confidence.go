package resolve

const (
	StatusHigh            = "high"
	StatusNoPath          = "no_path"
	StatusUnresolved      = "unresolved"
	StatusAmbiguous       = "ambiguous"
	StatusDuplicatePrefix = "duplicate_prefix"
	StatusRoleInvalid     = "role_invalid"
	StatusInvalidForMap   = "invalid_for_map"
	StatusMissingCoords   = "missing_coordinates"
	StatusMissingRF       = "missing_rf_evidence"
	StatusDistanceGate    = "rejected_distance_gate"
)

type Candidate struct {
	NodeID    string
	PublicKey string
	Name      string
	Role      string
	IATA      string
	Latitude  *float64
	Longitude *float64
}

type ResolvedHop struct {
	Prefix     string      `json:"prefix"`
	Confidence string      `json:"confidence"`
	Candidate  Candidate   `json:"candidate"`
	Candidates []Candidate `json:"candidates,omitempty"`
}

type Result struct {
	Status string        `json:"status"`
	Reason string        `json:"reason,omitempty"`
	Hops   []ResolvedHop `json:"hops,omitempty"`
}

func (r Result) IsHigh() bool {
	return r.Status == StatusHigh
}
