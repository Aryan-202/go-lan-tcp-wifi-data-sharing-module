package transfer

import (
	"encoding/json"
	"fmt"
	"os"
)

// TransferState records progress for resuming interrupted file transfers
type TransferState struct {
	FileID            string       `json:"file_id"`
	FileName          string       `json:"file_name"`
	CompletedSegments map[int]bool `json:"completed_segments"`
}

// NewTransferState creates a new TransferState instance for a given file ID
func NewTransferState(fileID string, fileName string) *TransferState {
	return &TransferState{
		FileID:            fileID,
		FileName:          fileName,
		CompletedSegments: make(map[int]bool),
	}
}

// MarkSegmentCompleted registers a segment index as successfully received and verified
func (ts *TransferState) MarkSegmentCompleted(index int) {
	if ts.CompletedSegments == nil {
		ts.CompletedSegments = make(map[int]bool)
	}
	ts.CompletedSegments[index] = true
}

// IsSegmentCompleted checks if a segment index was already received
func (ts *TransferState) IsSegmentCompleted(index int) bool {
	if ts.CompletedSegments == nil {
		return false
	}
	return ts.CompletedSegments[index]
}

// SaveState persists the progress to a JSON state file
func SaveState(filePath string, state *TransferState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transfer state: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write transfer state file: %w", err)
	}

	return nil
}

// LoadState reads progress from a JSON state file
func LoadState(filePath string) (*TransferState, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read transfer state file: %w", err)
	}

	var state TransferState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer state: %w", err)
	}

	if state.CompletedSegments == nil {
		state.CompletedSegments = make(map[int]bool)
	}

	return &state, nil
}
