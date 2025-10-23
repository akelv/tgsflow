package server

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewBacklog(t *testing.T) {
	b := NewBacklog()
	if b.Version != BacklogVersion {
		t.Errorf("expected version %s, got %s", BacklogVersion, b.Version)
	}
	if len(b.Thoughts) != 0 {
		t.Errorf("expected empty thoughts, got %d", len(b.Thoughts))
	}
}

func TestAddThought(t *testing.T) {
	b := NewBacklog()

	err := b.AddThought("tgs/thoughts/test-123", "Test thought", 10)
	if err != nil {
		t.Fatalf("failed to add thought: %v", err)
	}

	if len(b.Thoughts) != 1 {
		t.Fatalf("expected 1 thought, got %d", len(b.Thoughts))
	}

	th := b.Thoughts[0]
	if th.Dir != "tgs/thoughts/test-123" {
		t.Errorf("expected dir tgs/thoughts/test-123, got %s", th.Dir)
	}
	if th.Title != "Test thought" {
		t.Errorf("expected title 'Test thought', got %s", th.Title)
	}
	if th.Priority != 10 {
		t.Errorf("expected priority 10, got %d", th.Priority)
	}
	if th.Status != StatusQueued {
		t.Errorf("expected status %s, got %s", StatusQueued, th.Status)
	}
}

func TestAddThoughtInvalidDir(t *testing.T) {
	b := NewBacklog()

	err := b.AddThought("invalid/path", "Test", 0)
	if err == nil {
		t.Fatal("expected error for invalid dir, got nil")
	}
}

func TestAddThoughtDuplicate(t *testing.T) {
	b := NewBacklog()

	err := b.AddThought("tgs/thoughts/test-123", "Test", 0)
	if err != nil {
		t.Fatalf("failed to add first thought: %v", err)
	}

	err = b.AddThought("tgs/thoughts/test-123", "Test2", 0)
	if err == nil {
		t.Fatal("expected error for duplicate thought, got nil")
	}
}

func TestRemoveThought(t *testing.T) {
	b := NewBacklog()

	b.AddThought("tgs/thoughts/test-123", "Test", 0)
	if len(b.Thoughts) != 1 {
		t.Fatal("failed to add thought")
	}

	err := b.RemoveThought("tgs/thoughts/test-123")
	if err != nil {
		t.Fatalf("failed to remove thought: %v", err)
	}

	if len(b.Thoughts) != 0 {
		t.Errorf("expected 0 thoughts after remove, got %d", len(b.Thoughts))
	}
}

func TestRemoveThoughtNotFound(t *testing.T) {
	b := NewBacklog()

	err := b.RemoveThought("tgs/thoughts/nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent thought, got nil")
	}
}

func TestClaimNext(t *testing.T) {
	b := NewBacklog()

	// Add multiple thoughts with different priorities
	b.AddThought("tgs/thoughts/low-priority", "Low", 1)
	b.AddThought("tgs/thoughts/high-priority", "High", 10)
	b.AddThought("tgs/thoughts/med-priority", "Med", 5)

	// Claim next should return highest priority
	claimed, err := b.ClaimNext("test-claimer")
	if err != nil {
		t.Fatalf("failed to claim next: %v", err)
	}

	if claimed == nil {
		t.Fatal("expected claimed thought, got nil")
	}

	if claimed.Dir != "tgs/thoughts/high-priority" {
		t.Errorf("expected high-priority thought, got %s", claimed.Dir)
	}

	if claimed.Status != StatusInProgress {
		t.Errorf("expected status %s, got %s", StatusInProgress, claimed.Status)
	}

	if claimed.ClaimedBy == nil || *claimed.ClaimedBy != "test-claimer" {
		t.Error("expected claimedBy to be set")
	}
}

func TestClaimNextEmptyBacklog(t *testing.T) {
	b := NewBacklog()

	claimed, err := b.ClaimNext("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claimed != nil {
		t.Errorf("expected nil for empty backlog, got %+v", claimed)
	}
}

func TestUpdateStatus(t *testing.T) {
	b := NewBacklog()
	b.AddThought("tgs/thoughts/test-123", "Test", 0)

	err := b.UpdateStatus("tgs/thoughts/test-123", StatusInProgress)
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	thought, _ := b.GetThought("tgs/thoughts/test-123")
	if thought.Status != StatusInProgress {
		t.Errorf("expected status %s, got %s", StatusInProgress, thought.Status)
	}
}

func TestUpdateStatusInvalidTransition(t *testing.T) {
	b := NewBacklog()
	b.AddThought("tgs/thoughts/test-123", "Test", 0)

	// Mark as completed
	b.UpdateStatus("tgs/thoughts/test-123", StatusInProgress)
	b.UpdateStatus("tgs/thoughts/test-123", StatusCompleted)

	// Try to transition from completed to something else (should fail)
	err := b.UpdateStatus("tgs/thoughts/test-123", StatusQueued)
	if err == nil {
		t.Fatal("expected error for invalid transition, got nil")
	}
}

func TestListThoughts(t *testing.T) {
	b := NewBacklog()

	// Add in random order
	b.AddThought("tgs/thoughts/low", "Low", 1)
	time.Sleep(10 * time.Millisecond)
	b.AddThought("tgs/thoughts/high", "High", 10)
	time.Sleep(10 * time.Millisecond)
	b.AddThought("tgs/thoughts/high2", "High2", 10)

	thoughts := b.ListThoughts()

	// Should be sorted by priority (desc), then by added time (asc)
	if len(thoughts) != 3 {
		t.Fatalf("expected 3 thoughts, got %d", len(thoughts))
	}

	// First two should be priority 10
	if thoughts[0].Priority != 10 || thoughts[1].Priority != 10 {
		t.Error("expected first two to have priority 10")
	}

	// Among priority 10, earlier added time should come first
	if thoughts[0].Dir != "tgs/thoughts/high" {
		t.Errorf("expected high to be first, got %s", thoughts[0].Dir)
	}

	// Last should be priority 1
	if thoughts[2].Priority != 1 {
		t.Errorf("expected last to have priority 1, got %d", thoughts[2].Priority)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	backlogPath := filepath.Join(tmpDir, "backlog.json")

	// Create and save
	b := NewBacklog()
	b.AddThought("tgs/thoughts/test-123", "Test thought", 5)

	err := b.Save(backlogPath)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(backlogPath); err != nil {
		t.Fatalf("backlog file not created: %v", err)
	}

	// Load
	loaded, err := Load(backlogPath)
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.Version != BacklogVersion {
		t.Errorf("version mismatch after load")
	}

	if len(loaded.Thoughts) != 1 {
		t.Fatalf("expected 1 thought after load, got %d", len(loaded.Thoughts))
	}

	th := loaded.Thoughts[0]
	if th.Dir != "tgs/thoughts/test-123" {
		t.Errorf("dir mismatch after load: %s", th.Dir)
	}
	if th.Title != "Test thought" {
		t.Errorf("title mismatch after load: %s", th.Title)
	}
	if th.Priority != 5 {
		t.Errorf("priority mismatch after load: %d", th.Priority)
	}
}

func TestLoadNonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	backlogPath := filepath.Join(tmpDir, "nonexistent.json")

	// Should return empty backlog, not error
	b, err := Load(backlogPath)
	if err != nil {
		t.Fatalf("unexpected error loading nonexistent file: %v", err)
	}

	if len(b.Thoughts) != 0 {
		t.Errorf("expected empty backlog, got %d thoughts", len(b.Thoughts))
	}
}

func TestValidate(t *testing.T) {
	b := NewBacklog()
	b.AddThought("tgs/thoughts/test-123", "Test", 0)

	errors := b.Validate()
	if len(errors) != 0 {
		t.Errorf("expected no validation errors, got: %v", errors)
	}
}

func TestValidateDuplicates(t *testing.T) {
	b := NewBacklog()

	// Manually create duplicate (bypassing AddThought validation)
	b.Thoughts = []ThoughtEntry{
		{Dir: "tgs/thoughts/test-123", Title: "Test", Priority: 0, Status: StatusQueued, AddedAt: time.Now()},
		{Dir: "tgs/thoughts/test-123", Title: "Test2", Priority: 0, Status: StatusQueued, AddedAt: time.Now()},
	}

	errors := b.Validate()
	if len(errors) == 0 {
		t.Error("expected validation errors for duplicates, got none")
	}
}

func TestValidateInvalidStatus(t *testing.T) {
	b := NewBacklog()

	// Manually create thought with invalid status
	b.Thoughts = []ThoughtEntry{
		{Dir: "tgs/thoughts/test-123", Title: "Test", Priority: 0, Status: "invalid_status", AddedAt: time.Now()},
	}

	errors := b.Validate()
	if len(errors) == 0 {
		t.Error("expected validation errors for invalid status, got none")
	}
}
