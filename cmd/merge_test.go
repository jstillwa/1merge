package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"1merge/internal/items"
	"1merge/internal/models"
	"1merge/internal/op"
)

type stubOpClient struct {
	calls      []string
	editErr    error
	archiveErr error
}

func (s *stubOpClient) RunOpCmd(args ...string) ([]byte, error) {
	s.calls = append(s.calls, strings.Join(args, " "))
	// Check if this is an edit command
	if len(args) >= 2 && args[0] == "item" && args[1] == "edit" && s.editErr != nil {
		return nil, s.editErr
	}
	// Check if this is an archive (delete) command
	if len(args) >= 2 && args[0] == "item" && args[1] == "delete" && s.archiveErr != nil {
		return nil, s.archiveErr
	}
	return nil, nil
}

func TestApplyMergeAndReport_LogsSuccess(t *testing.T) {
	stub := &stubOpClient{}
	items.SetOpClient(stub)
	t.Cleanup(func() { items.SetOpClient(op.DefaultClient) })

	winner := models.Item{ID: "winner"}
	losers := []models.Item{{ID: "loser1"}, {ID: "loser2"}}

	var out bytes.Buffer
	if err := applyMergeAndReport(&out, winner, losers, false); err != nil {
		t.Fatalf("applyMergeAndReport returned error: %v", err)
	}

	if !strings.Contains(out.String(), "Successfully merged 2 items into winner") {
		t.Fatalf("expected success message to be printed, got %q", out.String())
	}

	if len(stub.calls) != 3 {
		t.Fatalf("expected three op calls (edit + two archives), got %d: %v", len(stub.calls), stub.calls)
	}
}

func TestApplyMergeAndReport_ForwardsError(t *testing.T) {
	stub := &stubOpClient{archiveErr: errors.New("archive failed")}
	items.SetOpClient(stub)
	t.Cleanup(func() { items.SetOpClient(op.DefaultClient) })

	winner := models.Item{ID: "winner"}
	losers := []models.Item{{ID: "loser1"}}

	var out bytes.Buffer
	err := applyMergeAndReport(&out, winner, losers, false)
	if err == nil {
		t.Fatal("expected error from applyMergeAndReport, got nil")
	}

	if out.Len() != 0 {
		t.Fatalf("expected no success output when an error occurs, got %q", out.String())
	}

	if len(stub.calls) != 2 {
		t.Fatalf("expected edit and one archive attempts, got %d: %v", len(stub.calls), stub.calls)
	}
}
