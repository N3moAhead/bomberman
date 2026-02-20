package match

import (
	"encoding/json"
	"time"

	"github.com/N3moAhead/bombahead/match_runner/internal/history"
)

// Details represents the information about a match to be run
type Details struct {
	MatchID      string `json:"match_id"`
	ServerImage  string `json:"server_image"`
	Client1Image string `json:"client1_image"`
	Client2Image string `json:"client2_image"`
}

// Result represents the outcome of a match
type Result struct {
	MatchID       string               `json:"match_id"`
	Winner        string               `json:"winner"` // Name of the client image that won
	Client1GameID string               `json:"client1GameId"`
	Client2GameID string               `json:"client2GameId"`
	Log           *history.GameHistory `json:"log"`
}

// Failure represents a permanently failed match handling attempt.
type Failure struct {
	MatchID    string          `json:"match_id"`
	Reason     string          `json:"reason"`
	Error      string          `json:"error"`
	RetryCount int             `json:"retry_count"`
	FailedAt   time.Time       `json:"failed_at"`
	Payload    json.RawMessage `json:"payload"`
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

// ToJSON encodes a Failure struct to a JSON byte slice
func (f *Failure) ToJSON() ([]byte, error) {
	return json.Marshal(f)
}
