package models

import "time"

type Peer struct {
	ID       string    `json:"id"`
	IP       string    `json:"ip"`
	Port     int       `json:"port"`
	LastSeen time.Time `json:"last_seen"`
	IsBusy   bool      `json:"is_busy"`
}

type FileInfo struct {
	Name       string
	Size       int64
	FullHash   []byte
	ChunkSize  int
	ChunkCount int
}

type ChunkInfo struct {
	Index   int
	Status  string
	Hash    []byte
	Offset  int64
	Length  int
	Retries int
}

type TransferState struct {
	SessionID   string
	LastUpdated time.Time
	Chunks      []ChunkInfo
}
