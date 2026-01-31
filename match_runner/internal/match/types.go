package match

import "encoding/json"

// Details represents the information about a match to be run
type Details struct {
	MatchID      string `json:"match_id"`
	ServerImage  string `json:"server_image"`
	Client1Image string `json:"client1_image"`
	Client2Image string `json:"client2_image"`
}

// Result represents the outcome of a match
type Result struct {
	MatchID string `json:"match_id"`
	Winner  string `json:"winner"` // Name of the client image that won
	Log     string `json:"log"`
}

// ToJSON encodes a Details struct to a JSON byte slice
func (d *Details) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

// FromJSON decodes a Details struct from a JSON byte slice
func (d *Details) FromJSON(data []byte) error {
	return json.Unmarshal(data, d)
}

// ToJSON encodes a Result struct to a JSON byte slice
func (r *Result) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
