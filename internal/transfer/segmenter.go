package transfer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

// GenerateManifest inspects a local file and generates its FileManifest metadata
func GenerateManifest(filePath string, segmentSize int64) (*FileManifest, error) {
	if segmentSize <= 0 {
		segmentSize = DefaultSegmentSize
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file for manifest: %w", err)
	}

	fileSize := info.Size()
	totalSegments := int(math.Ceil(float64(fileSize) / float64(segmentSize)))
	if totalSegments == 0 {
		totalSegments = 1
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for manifest: %w", err)
	}
	defer file.Close()

	// Calculate full file SHA-256
	fileHasher := sha256.New()
	if _, err := io.Copy(fileHasher, file); err != nil {
		return nil, fmt.Errorf("failed to compute file SHA-256: %w", err)
	}
	fullFileHash := hex.EncodeToString(fileHasher.Sum(nil))

	// Calculate per-segment SHA-256 hashes
	segmentHashes := make([]string, totalSegments)
	for i := 0; i < totalSegments; i++ {
		segmentData, hash, err := ReadSegment(file, i, segmentSize)
		if err != nil && len(segmentData) == 0 {
			return nil, fmt.Errorf("failed to read segment %d for manifest: %w", i, err)
		}
		segmentHashes[i] = hash
	}

	fileID := fmt.Sprintf("%s-%d-%s", filepath.Base(filePath), fileSize, fullFileHash[:8])

	return &FileManifest{
		FileID:        fileID,
		FileName:      filepath.Base(filePath),
		FileSize:      fileSize,
		SegmentSize:   segmentSize,
		TotalSegments: totalSegments,
		FileHash:      fullFileHash,
		SegmentHashes: segmentHashes,
	}, nil
}

// ReadSegment reads a specific segment from an open file thread-safely
func ReadSegment(file *os.File, segmentIndex int, segmentSize int64) ([]byte, string, error) {
	offset := int64(segmentIndex) * segmentSize
	buffer := make([]byte, segmentSize)

	n, err := file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		return nil, "", fmt.Errorf("failed to read segment at index %d: %w", segmentIndex, err)
	}

	segmentData := buffer[:n]
	hash := sha256.Sum256(segmentData)
	return segmentData, hex.EncodeToString(hash[:]), nil
}
