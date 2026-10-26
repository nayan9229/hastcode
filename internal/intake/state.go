package intake

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// StateFile is the seen-state's file name in the state directory.
const StateFile = "state.json"

// stateVersion is the state file's format version.
const stateVersion = 1

// Decision is what became of a Lead.
type Decision string

const (
	// Open: listed, not yet decided; it shows in the queue.
	Open Decision = "open"
	// Picked: turned into a draft skeleton (Pick).
	Picked Decision = "picked"
	// Dismissed: hidden for good (Dismiss).
	Dismissed Decision = "dismissed"
	// Expired: left open too long and dismissed by a run.
	Expired Decision = "expired"
	// TooOld: seen by a run but published before its age cutoff
	// (Intake.MaxAge), so never listed: not a Lead the Author decides on.
	// A later run with a wider MaxAge that reaches it lists it as new.
	TooOld Decision = "too-old"
)

// state is the seen-state file: every Lead ever listed, keyed by ID, with
// its decision, and each Source's state between runs. JSON object keys
// marshal sorted, so the file is stable between runs.
type state struct {
	Version int                    `json:"version"`
	Leads   map[string]Lead        `json:"leads"`
	Sources map[string]SourceState `json:"sources"`
}

func loadState(dir string) (*state, error) {
	st := &state{Version: stateVersion, Leads: map[string]Lead{}, Sources: map[string]SourceState{}}
	data, err := os.ReadFile(filepath.Join(dir, StateFile))
	if errors.Is(err, fs.ErrNotExist) {
		return st, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, st); err != nil {
		return nil, fmt.Errorf("%s: %w", filepath.Join(dir, StateFile), err)
	}
	if st.Version != stateVersion {
		return nil, fmt.Errorf("%s: state version %d, want %d", filepath.Join(dir, StateFile), st.Version, stateVersion)
	}
	if st.Leads == nil {
		st.Leads = map[string]Lead{}
	}
	if st.Sources == nil {
		st.Sources = map[string]SourceState{}
	}
	return st, nil
}

// save writes the state atomically (a temp file renamed over the old one),
// so an interrupted run never leaves half a state file.
func (st *state) save(dir string) error {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(filepath.Join(dir, StateFile), append(data, '\n'))
}

func writeFileAtomic(file string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(file), "."+filepath.Base(file)+".*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), file)
}
