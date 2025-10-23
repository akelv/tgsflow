// Package server provides work orchestration for TGS Server Mode.
package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// BacklogVersion is the current backlog file format version.
const BacklogVersion = "1"

// Thought status values.
const (
	StatusQueued     = "queued"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// Backlog represents the work queue for approved thoughts.
type Backlog struct {
	Version  string         `json:"version"`
	Thoughts []ThoughtEntry `json:"thoughts"`
	mu       sync.Mutex     // Protects concurrent access
}

// ThoughtEntry represents a single thought in the backlog.
type ThoughtEntry struct {
	Dir         string     `json:"dir"`
	Title       string     `json:"title"`
	Priority    int        `json:"priority"`
	Status      string     `json:"status"`
	AddedAt     time.Time  `json:"added_at"`
	ClaimedBy   *string    `json:"claimed_by"`
	ClaimedAt   *time.Time `json:"claimed_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// NewBacklog creates a new empty backlog.
func NewBacklog() *Backlog {
	return &Backlog{
		Version:  BacklogVersion,
		Thoughts: []ThoughtEntry{},
	}
}

// Load reads a backlog from a JSON file.
func Load(path string) (*Backlog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty backlog if file doesn't exist
			return NewBacklog(), nil
		}
		return nil, fmt.Errorf("read backlog: %w", err)
	}

	var b Backlog
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("unmarshal backlog: %w", err)
	}

	// Validate version
	if b.Version != BacklogVersion {
		return nil, fmt.Errorf("unsupported backlog version: %s (expected %s)", b.Version, BacklogVersion)
	}

	return &b, nil
}

// Save writes the backlog to a JSON file.
func (b *Backlog) Save(path string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create backlog dir: %w", err)
	}

	// Marshal with indentation for readability
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal backlog: %w", err)
	}

	// Write atomically via temp file + rename
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write backlog temp: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename backlog: %w", err)
	}

	return nil
}

// AddThought adds a new thought to the backlog.
func (b *Backlog) AddThought(dir, title string, priority int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Validate directory path
	if !strings.HasPrefix(dir, "tgs/thoughts/") {
		return fmt.Errorf("invalid thought dir: must start with tgs/thoughts/")
	}

	// Check for duplicate
	for _, t := range b.Thoughts {
		if t.Dir == dir {
			return fmt.Errorf("thought already in backlog: %s", dir)
		}
	}

	// Add new entry
	now := time.Now()
	b.Thoughts = append(b.Thoughts, ThoughtEntry{
		Dir:      dir,
		Title:    title,
		Priority: priority,
		Status:   StatusQueued,
		AddedAt:  now,
	})

	return nil
}

// RemoveThought removes a thought from the backlog.
func (b *Backlog) RemoveThought(dir string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, t := range b.Thoughts {
		if t.Dir == dir {
			// Remove by swapping with last element and truncating
			b.Thoughts[i] = b.Thoughts[len(b.Thoughts)-1]
			b.Thoughts = b.Thoughts[:len(b.Thoughts)-1]
			return nil
		}
	}

	return fmt.Errorf("thought not found: %s", dir)
}

// ClaimNext atomically claims the next available thought.
// Returns the claimed thought entry or nil if backlog is empty.
func (b *Backlog) ClaimNext(claimedBy string) (*ThoughtEntry, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Find highest priority queued thought
	var best *ThoughtEntry
	var bestIdx int

	for i := range b.Thoughts {
		t := &b.Thoughts[i]
		if t.Status != StatusQueued {
			continue
		}
		if best == nil || t.Priority > best.Priority {
			best = t
			bestIdx = i
		}
	}

	if best == nil {
		return nil, nil // No work available
	}

	// Claim it
	now := time.Now()
	b.Thoughts[bestIdx].Status = StatusInProgress
	b.Thoughts[bestIdx].ClaimedBy = &claimedBy
	b.Thoughts[bestIdx].ClaimedAt = &now
	b.Thoughts[bestIdx].StartedAt = &now

	// Return a copy
	claimed := b.Thoughts[bestIdx]
	return &claimed, nil
}

// UpdateStatus updates the status of a thought.
func (b *Backlog) UpdateStatus(dir, status string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Validate status
	validStatuses := map[string]bool{
		StatusQueued:     true,
		StatusInProgress: true,
		StatusCompleted:  true,
		StatusFailed:     true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	// Find thought
	for i := range b.Thoughts {
		if b.Thoughts[i].Dir == dir {
			oldStatus := b.Thoughts[i].Status

			// Validate transition
			if err := validateStatusTransition(oldStatus, status); err != nil {
				return err
			}

			// Update status
			b.Thoughts[i].Status = status

			// Update timestamp if completing
			if status == StatusCompleted || status == StatusFailed {
				now := time.Now()
				b.Thoughts[i].CompletedAt = &now
			}

			return nil
		}
	}

	return fmt.Errorf("thought not found: %s", dir)
}

// ListThoughts returns a copy of all thoughts, sorted by priority (desc) then added time (asc).
func (b *Backlog) ListThoughts() []ThoughtEntry {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Make a copy
	thoughts := make([]ThoughtEntry, len(b.Thoughts))
	copy(thoughts, b.Thoughts)

	// Sort by priority (desc), then by added time (asc)
	sort.Slice(thoughts, func(i, j int) bool {
		if thoughts[i].Priority != thoughts[j].Priority {
			return thoughts[i].Priority > thoughts[j].Priority
		}
		return thoughts[i].AddedAt.Before(thoughts[j].AddedAt)
	})

	return thoughts
}

// GetThought returns a thought by directory path.
func (b *Backlog) GetThought(dir string) (*ThoughtEntry, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := range b.Thoughts {
		if b.Thoughts[i].Dir == dir {
			// Return a copy
			t := b.Thoughts[i]
			return &t, nil
		}
	}

	return nil, fmt.Errorf("thought not found: %s", dir)
}

// Validate checks backlog consistency.
func (b *Backlog) Validate() []string {
	b.mu.Lock()
	defer b.mu.Unlock()

	var errors []string

	// Check for duplicate dirs
	seen := make(map[string]bool)
	for _, t := range b.Thoughts {
		if seen[t.Dir] {
			errors = append(errors, fmt.Sprintf("duplicate thought dir: %s", t.Dir))
		}
		seen[t.Dir] = true

		// Validate dir format
		if !strings.HasPrefix(t.Dir, "tgs/thoughts/") {
			errors = append(errors, fmt.Sprintf("invalid thought dir format: %s", t.Dir))
		}

		// Validate status
		validStatuses := map[string]bool{
			StatusQueued:     true,
			StatusInProgress: true,
			StatusCompleted:  true,
			StatusFailed:     true,
		}
		if !validStatuses[t.Status] {
			errors = append(errors, fmt.Sprintf("invalid status for %s: %s", t.Dir, t.Status))
		}

		// Check claimed thoughts have claimedBy
		if t.Status == StatusInProgress && t.ClaimedBy == nil {
			errors = append(errors, fmt.Sprintf("in_progress thought missing claimedBy: %s", t.Dir))
		}
	}

	return errors
}

// validateStatusTransition checks if a status transition is valid.
func validateStatusTransition(oldStatus, newStatus string) error {
	// Valid transitions:
	// queued -> in_progress, failed
	// in_progress -> completed, failed, queued (retry)
	// completed -> (none)
	// failed -> queued (retry)

	if oldStatus == newStatus {
		return nil // No-op
	}

	switch oldStatus {
	case StatusQueued:
		if newStatus == StatusInProgress || newStatus == StatusFailed {
			return nil
		}
	case StatusInProgress:
		if newStatus == StatusCompleted || newStatus == StatusFailed || newStatus == StatusQueued {
			return nil
		}
	case StatusFailed:
		if newStatus == StatusQueued {
			return nil
		}
	case StatusCompleted:
		return fmt.Errorf("cannot transition from completed to %s", newStatus)
	}

	return fmt.Errorf("invalid status transition: %s -> %s", oldStatus, newStatus)
}
