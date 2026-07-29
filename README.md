# GoShare: High-Performance Cross-Platform Peer-to-Peer File Transfer Engine

GoShare is a decentralized, cross-platform peer-to-peer (P2P) file sharing application built in Go (Golang). It enables high-speed, secure, and direct file transfers across local area networks (LAN/Wi-Fi) without relying on external servers, cloud storage, or internet connections.

The system features automatic local service discovery, end-to-end cryptographic key negotiation, authenticated payload encryption, parallel stream transmission, and interrupted state resumption.

---

## Key Features

- **Zero-Configuration mDNS Peer Discovery**: Automatically discovers active GoShare instances on local subnets using Multicast DNS (`_goshare._tcp.local`) with manual IP fallback for firewall-restricted environments.
- **End-to-End Cryptography**: Secures all data exchanges via ephemeral X25519 Elliptic Curve Diffie-Hellman (ECDH) key agreement, HKDF-SHA256 key expansion, and AES-256-GCM authenticated encryption.
- **High-Throughput Parallel Transfer**: Partitions files of any size into 4MB data segments, streaming encrypted segments concurrently across TCP worker sockets to saturate LAN bandwidth.
- **Data Integrity Verification**: Calculates and validates SHA-256 checksums for both the complete file and individual data segments prior to disk writing.
- **Interrupted Transfer Resumption**: Tracks segment receipt states to resume interrupted downloads from the last successfully written segment.
- **Cross-Platform Compatibility**: Native binary builds supporting Linux, macOS, Windows, Android (ARM64 PIE), and iOS.

---

## Architecture Overview

GoShare is built with a modular package architecture following standard Go project conventions:

```text
.
├── cmd/                # Command-Line Interface (Cobra CLI definitions)
│   ├── root.go         # Root command & argument sanitizer
│   ├── discover.go     # mDNS discovery & manual IP probing command
│   ├── receive.go      # TCP listener & file receiver command
│   └── send.go         # Parallel TCP file sender command
├── internal/           # Private application & engine packages
│   ├── discovery/      # mDNS announcer, resolver, and peer management
│   ├── security/       # X25519, HKDF-SHA256, AES-256-GCM, & passphrase pairing
│   └── transfer/       # File segmenter, manifest, protocol framing, & state persistence
├── test/               # Comprehensive unit and integration test suites
│   ├── commands/       # CLI command execution tests
│   ├── discovery/       # Discovery package unit tests
│   ├── security/       # Security & cryptography unit tests
│   └── transfer/      # File transfer engine & state tests
├── docs/               # Software Requirements Specification (SRS)
├── main.go             # Main application entry point
├── go.mod              # Go module definitions
└── go.sum              # Dependency checksums
```

---

## Security Architecture

GoShare implements zero-trust end-to-end security designed to prevent network eavesdropping, tampering, and Man-in-the-Middle (MITM) attacks.

### Cryptographic Workflow

1. **Ephemeral Key Pair Generation**: Upon initiating a connection, each peer generates a fresh X25519 private/public key pair using `crypto/ecdh`.
2. **Public Key Exchange**: Peers exchange 32-byte raw public key bytes over the established TCP socket.
3. **Shared Secret Computation**: Each peer computes the raw 32-byte Diffie-Hellman shared secret via `privKey.ECDH(peerPubKey)`.
4. **HKDF Key Expansion**: The raw secret is passed through HKDF-SHA256 (`golang.org/x/crypto/hkdf`) with a domain-separated salt and info label to derive a uniform 256-bit (32-byte) AES key.
5. **AES-256-GCM AEAD Encryption**: Every 4MB data segment is encrypted using AES-256-GCM (`crypto/cipher`). A 12-byte cryptographically secure random nonce (`crypto/rand`) is prepended to each ciphertext payload, along with segment-specific Additional Authenticated Data (AAD).
6. **Passphrase Verification**: Peer connections support verification using a 4-word pairing passphrase generated via a cryptographically secure random word selection algorithm.

---

## Installation & Build

### Prerequisites

- Go 1.25 or higher
- Git

### Build from Source

```bash
# Clone repository
git clone https://github.com/Aryan-202/go-lan-tcp-wifi-data-sharing-module.git
cd go-lan-tcp-wifi-data-sharing-module

# Build native binary
go build -o goshare main.go
```

### Cross-Compilation for Mobile (Android ARM64)

To build a Position-Independent Executable (PIE) binary compatible with Android / Termux:

```bash
GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build -o goshare_android main.go
```

---

## Usage

### 1. Peer Discovery

Discover active GoShare instances on your local subnet:

```bash
# Discover peers with default 5-second timeout
./goshare discover

# Discover with custom 2-second timeout
./goshare discover --timeout 2s

# Manual IP probe fallback for firewalled networks
./goshare discover --ip 192.168.1.50
```

### 2. File Receiver

Start listening for incoming file transfers:

```bash
# Listen on default port 8829 and save to ~/Downloads
./goshare receive --dir ~/Downloads

# Listen on custom port
./goshare receive --port 9000 --dir ./custom_folder
```

### 3. File Sender

Encrypt and transmit a file to a peer:

```bash
# Send file to target peer IP
./goshare send /path/to/file.mp4 --ip 192.168.1.50

# Send using custom port and passphrase verification
./goshare send /path/to/file.mp4 --ip 192.168.1.50 --port 9000 --passphrase "orbit-falcon-amber-crest"
```

---

## Verification & Testing

GoShare includes unit and integration test suites covering security algorithms, mDNS resolution, file segmentation, and CLI execution.

```bash
# Run all unit and integration test suites
go test -v ./test/...

# Measure code coverage across packages
go test -v -coverpkg=./... ./test/...
```

---

## License

This project is licensed under the MIT License.

---

## Author

- **Aryan Vishwakarma**
- **GitHub**: [Aryan-202](https://github.com/Aryan-202)
- **Email**: [aryanvishwakarma275@gmail.com](mailto:aryanvishwakarma275@gmail.com)
