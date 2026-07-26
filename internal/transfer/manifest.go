package transfer

// DefaultSegmentSize is 4 MB as specified in SRS Section 3.2
const DefaultSegmentSize = 4 * 1024 * 1024

// FileManifest contains metadata about a file being prepared or sent
type FileManifest struct {
	FileID        string   `json:"file_id"`
	FileName      string   `json:"file_name"`
	FileSize      int64    `json:"file_size"`
	SegmentSize   int64    `json:"segment_size"`
	TotalSegments int      `json:"total_segments"`
	FileHash      string   `json:"file_hash"`
	SegmentHashes []string `json:"segment_hashes"`
}

// SegmentPacket represents a network packet carrying an encrypted segment payload
type SegmentPacket struct {
	FileID       string `json:"file_id"`
	SegmentIndex int    `json:"segment_index"`
	DataLength   int    `json:"data_length"`
	Payload      []byte `json:"payload"` // AES-256-GCM encrypted segment data
}