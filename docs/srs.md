![Realistic image of two command line interfaces showing file transfer
progress on nearby laptops connected over a local
network.](media/image1.jpg){width="8.0in" height="2.2222222222222223in"}

1\. INTRODUCTION

# 1.1 Purpose

This Software Requirements Specification (SRS) documents the complete
requirements for \"P2Pxfer,\" a Command-Line Interface (CLI)-first,
peer-to-peer (P2P) local file transfer system. P2Pxfer is designed to
enable secure, fast, and robust direct file exchange between devices on
the same Local Area Network (LAN) without relying on intermediate
servers, complex router configuration, or physical cables.

# 1.2 Scope

P2Pxfer\'s primary scope is to facilitate high-performance, resilient,
and secure transfers of large files, specifically targeting the range of
20GB to 100GB+. It is designed for environments where traditional
internet connectivity is unavailable, slow, or undesirable for security
and speed reasons. The system is inherently decentralized and requires
**no external internet connection, pre-configured routers, or physical
cables** beyond the existing network infrastructure (e.g., Wi-Fi or
local Ethernet).

# 1.3 Definitions and Acronyms

  -----------------------------------------------------------------------
  Term                                Definition
  ----------------------------------- -----------------------------------
  P2P                                 Peer-to-Peer: Direct communication
                                      between two computing systems.

  UDP                                 User Datagram Protocol:
                                      Connectionless network protocol,
                                      used for discovery.

  TCP                                 Transmission Control Protocol:
                                      Connection-oriented network
                                      protocol, used for reliable data
                                      transfer.

  mDNS                                Multicast Domain Name System:
                                      Protocol for local service
                                      discovery.

  SHA-256                             Secure Hash Algorithm 256-bit:
                                      Cryptographic hash function used
                                      for integrity checking.

  AES-GCM                             Advanced Encryption Standard -
                                      Galois/Counter Mode: Authenticated
                                      encryption scheme used for
                                      security.

  CLI                                 Command-Line Interface: Text-based
                                      user interface.

  LAN                                 Local Area Network: Computer
                                      network that interconnects devices
                                      within a limited area.

  MTU                                 Maximum Transmission Unit: The size
                                      of the largest packet that can be
                                      transmitted.

  SRS                                 Software Requirements
                                      Specification: Document detailing
                                      requirements.
  -----------------------------------------------------------------------

# 1.4 References

  -----------------------------------------------------------------------
  Reference                           Description
  ----------------------------------- -----------------------------------
  IEEE 830                            IEEE Recommended Practice for
                                      Software Requirements
                                      Specifications.

  RFC 768                             User Datagram Protocol (UDP).

  RFC 793                             Transmission Control Protocol
                                      (TCP).

  RFC 6762                            Multicast DNS (mDNS).

  Syncthing                           Open-source file synchronization
                                      application (used as a conceptual
                                      reference for decentralized
                                      design).

  croc                                Open-source tool for peer-to-peer
                                      file transfer (used as a conceptual
                                      reference for CLI implementation).

  Magic Wormhole                      Secure, simple file transfer
                                      application (used as a conceptual
                                      reference for secure
                                      authentication/pairing).
  -----------------------------------------------------------------------

2\. OVERALL DESCRIPTION

# 2.1 Product Perspective

P2Pxfer is a **standalone, decentralized** file transfer utility. It
does not require any central server, cloud services, or external
accounts, operating entirely within the local subnet. It is packaged as
a single binary executable for simplicity and portability.

# 2.2 Product Functions

P2Pxfer must provide the following core functions:

- **Auto-Discovery:** Automatically locate other P2Pxfer peers on the
  local subnet using UDP broadcasting.

- **Direct P2P TCP Transfer:** Establish a reliable, direct TCP
  connection between peers for data exchange.

- **Chunked Transfer:** Split files into smaller, manageable chunks for
  efficient parallel transfer, integrity checking, and resume
  capability.

- **Pause/Resume:** Allow transfers to be interrupted and resumed later,
  maintaining transfer state persistence.

- **Integrity Check:** Verify the correctness of data using
  cryptographic hashing per chunk and for the full file.

- **Parallel Transfer:** Utilize multiple concurrent connections (worker
  pools) to maximize local network throughput.

- **Encryption:** Secure all data transferred over the network using
  strong, authenticated encryption.

- **Compression:** Apply compression (where effective) to reduce network
  payload and improve speed.

# 2.3 User Classes and Characteristics

  -----------------------------------------------------------------------
  User Class                          Characteristics
  ----------------------------------- -----------------------------------
  Developers                          Requires high control, scripting
                                      capability, needs advanced flags
                                      (port, buffer size).

  Power Users                         Values speed and security, uses the
                                      CLI, appreciates features like
                                      pause/resume and compression.

  General Users                       Requires zero-configuration setup,
                                      simple send/receive commands, clear
                                      progress display.
  -----------------------------------------------------------------------

# 2.4 Operating Environment

  -----------------------------------------------------------------------
  Component                           Requirement
  ----------------------------------- -----------------------------------
  OS Support                          Windows 10+, Linux Kernel 5+, macOS
                                      12+ (Future support for
                                      Android/iOS).

  Runtime                             Go 1.21+ runtime environment
                                      (compiled binary only).

  Permissions                         Must run without root/administrator
                                      privileges.

  Memory                              Memory usage must remain below 2x
                                      the configured chunk size for the
                                      transfer buffer.

  Network                             Requires local network (LAN, Wi-Fi,
                                      Ethernet) access. UDP and TCP ports
                                      must be allowed.
  -----------------------------------------------------------------------

# 2.5 Design and Implementation Constraints

1.  **Language:** Must be implemented entirely in Go (Golang) version
    1.21 or later.

2.  **Distribution:** Must be a single, statically linked binary.

3.  **Permissions:** Must not require root or administrator privileges
    for primary operation.

4.  **Assumptions:** The operating subnet allows UDP broadcast/multicast
    and TCP connections between hosts. The target filesystem supports
    files larger than 4GB.

3\. SYSTEM FEATURES

# F1 Discovery (Priority: Critical, Phase: 1)

The system must locate other P2Pxfer peers on the local subnet.

- **F1.1 Broadcast:** Use UDP broadcast/multicast on port 9998 to
  announce presence every 5 seconds.

- **F1.2 mDNS:** Future feature to leverage mDNS/Zeroconf for peer
  discovery.

- **F1.3 Manual Peer:** Support specifying a peer\'s IP address and port
  manually via CLI flag.

# F2 Connection (Priority: Critical, Phase: 1)

The system must establish and maintain a reliable connection for data
transfer.

- F2.1 Protocol: Utilize TCP on port 9999 (default) for reliable,
  ordered data stream transfer.

- F2.2 Management: Implement TCP keepalives and automatic reconnect
  logic for brief network interruptions.

# F3 Transfer Engine (Priority: Critical, Phase: 2)

The system must manage the efficient transfer of file data.

- **F3.1 Chunking:** Divide files into chunks, with a **4MB default
  chunk size**.

- **F3.2 Transfer Strategy:** Initially use a sequential, single-worker
  transfer for setup, transitioning to parallel transfer using worker
  pools (default 4 workers) for data transmission.

- **F3.3 Adaptive Sizing:** Future feature to dynamically adjust chunk
  size based on network conditions (latency/throughput).

# F4 Transfer Management (Priority: High, Phase: 2)

The system must allow users to control the transfer flow.

- **F4.1 Persistence:** Transfers must be resumable. State (completed
  chunks, peer info, session ID) must be persisted in a local JSON file.

- **F4.2 Control:** Provide CLI commands for initiating pause, resume,
  and cancel operations.

# F5 Integrity (Priority: Critical, Phase: 3)

The system must guarantee data correctness.

- **F5.1 Per-Chunk Hashing:** Calculate and verify a SHA-256 hash for
  every 4MB chunk.

- **F5.2 Auto Retry:** Automatically attempt to re-send a failed chunk
  up to **3 times** before reporting a fatal error.

- **F5.3 Full File Hashing:** Calculate and verify a full file SHA-256
  hash upon completion of the transfer.

# F6 Security (Priority: High, Phase: 3)

The system must protect data confidentiality and integrity from network
snooping.

- **F6.1 Encryption:** Encrypt all data using **AES-256-GCM**. A unique
  nonce must be generated for every chunk.

- **F6.2 Key Exchange:** Utilize a Diffie-Hellman (DH) key exchange
  method to establish a secure session key.

- **F6.3 Verification:** Implement secure peer authentication using
  **4-word verification codes** (passphrases) derived from a 2048-word
  dictionary, similar to Magic Wormhole.

# F7 Compression (Priority: Medium, Phase: 4)

The system should optimize transfer speed by reducing data size.

- **F7.1 Algorithm:** Use the **zstd** compression algorithm for its
  speed and efficiency.

- **F7.2 Heuristic:** Implement a heuristic to automatically detect
  files that are unlikely to benefit from compression (e.g., already
  compressed video/audio files) and disable compression for them.

# F8 CLI (Priority: Critical, Phase: 1)

The system\'s primary interface must be the command line.

- **F8.1 Framework:** Utilize the Cobra framework for robust command
  structure.

- **F8.2 Commands:** Provide atomic commands: send, receive, discover,
  status, cancel, pause, resume.

- **F8.3 Progress Display:** Implement a clear, real-time progress bar
  displaying completion percentage, transfer speed (MB/s), and Estimated
  Time of Arrival (ETA).

# F9 Future (Priority: Low, Phase: 5)

Potential extensions for future development.

- **F9.1 Swarm:** Extend P2P to N:N (many-to-many) transfers.

- **F9.2 Pairing:** Add QR code pairing for easy setup on
  mobile/desktop.

- **F9.3 Network:** Support protocols like WiFi Direct or mesh
  networking.

- **F9.4 Delta Sync:** Implement a block-level differential sync
  capability (like rsync).

4\. EXTERNAL INTERFACES

# 4.1 User Interfaces (CLI)

The primary interface is the CLI, built on Cobra. Commands must accept
the following flags (or similar short-flags):

  -----------------------------------------------------------------------
  Command                 Flag                    Description
  ----------------------- ----------------------- -----------------------
  send / receive          \--port                 Custom TCP port to
                                                  listen on.

  send                    \--chunk-size           Override the default
                                                  4MB chunk size.

  send                    \--parallel             Set the number of
                                                  parallel transfer
                                                  workers (default 4).

  send                    \--encrypt              Enable/Disable
                                                  encryption (default
                                                  enabled).

  send                    \--compress             Enable/Disable
                                                  compression (default
                                                  auto-detect).

  receive                 \--output               Specify the output
                                                  directory for received
                                                  files.

  receive                 \--peer                 Manually specify a peer
                                                  IP:Port for transfer.

  receive                 \--code                 Enter the verification
                                                  code for secure
                                                  connection.
  -----------------------------------------------------------------------

**Progress Display:**

The CLI must display progress in real-time, refreshing at 10Hz, using a
format similar
to:\[====================\>\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\--\] 60% @
85.5 MB/s, ETA 00:03:12

# 4.2 Configuration Interface (YAML)

P2Pxfer must support a persistent configuration file at
\~/.p2pxfer/config.yaml.

This file must manage:

- **Network Settings:** Default discovery/transfer ports, connection
  timeouts, broadcast intervals.

- **Transfer Settings:** Default chunk size, parallel workers, resume
  file location.

- **Security Settings:** Default encryption level, passphrase word list
  path.

- **Compression Settings:** Compression threshold, algorithm selection.

- **Display Settings:** Progress bar style, refresh rate.

5\. ARCHITECTURE

P2Pxfer must adopt a layered, modular architecture.

# 5.1 Architectural Layers

1.  **CLI Layer:** Handles user input, parsing flags, and formatting
    output (Cobra, BubbleTea).

2.  **Service Layer:** Translates CLI commands into core system calls,
    managing sessions and state (transfer persistence).

3.  **Core Modules Layer:** Contains primary logic for Discovery,
    Connection, Transfer Engine, Security, and Integrity.

4.  **Infrastructure Layer:** Handles networking (UDP/TCP sockets), file
    I/O, and cryptographic operations (Go stdlib, zerolog).

# 5.2 Node Diagram

\[Sender/receiver node diagram showing discovery via UDP broadcast then
a direct TCP transfer channel established between the two nodes.\]

The Sender (S) and Receiver (R) nodes operate as follows:

1.  **Discovery (S↔R):** Both nodes listen on UDP Port 9998. S sends a
    QUERY or ANNOUNCE. R responds with a RESPONSE.

2.  **Connection (S→R):** S initiates a TCP connection to R on Port
    9999.

3.  **Transfer (S↔R):** After authentication, S streams encrypted,
    chunked file data to R. R streams acknowledgements (ACKs) and NACKs
    (for retries).

# 5.3 Data Flows

  -----------------------------------------------------------------------
  Flow                                Processing Steps
  ----------------------------------- -----------------------------------
  **Send File**                       Parse CLI input → Discover peer →
                                      Calculate full file hash → Split
                                      file into chunks → Initiate DH key
                                      exchange → Parallel send chunks →
                                      Receive ACKs → Verify full transfer
                                      → Complete/Close.

  **Receive File**                    Listen on TCP port 9999 → Accept
                                      connection → DH Key Exchange →
                                      Authentication (verification code)
                                      → Receive file metadata → Parallel
                                      receive chunks → Verify chunk hash
                                      → Write to disk → Verify full file
                                      hash → Complete/Close.
  -----------------------------------------------------------------------

# 5.4 Module Interfaces (Go Pseudocode)

o\
// Discovery Module\
type PeerRegistry interface {\
Announce(metadata FileInfo)\
Discover() \[\]Peer\
GetPeer(id string) Peer\
}

// Transfer Engine Module\
type Session interface {\
StartSend(file FileInfo, peer Peer, key \[\]byte) error\
StartReceive(peer Peer, key \[\]byte) error\
Pause(sessionID string)\
Resume(sessionID string)\
Cancel(sessionID string)\
}

// Integrity Module\
type Hasher interface {\
HashChunk(data \[\]byte) \[\]byte // SHA-256\
VerifyFile(path string, expectedHash \[\]byte) bool\
}

\# 6. DETAILED REQUIREMENTS

This section details the critical functional requirements.

\## 6.1 Discovery Requirements

\| ID \| Priority \| Phase \| Requirement Description \|

\|\-\--\|\-\--\|\-\--\|\-\--\|

\| \*\*REQ-DISC-001\*\* \| Critical \| 1 \| The system must broadcast a
UDP \`ANNOUNCE\` message every 5 seconds on port 9998 to all hosts on
the local subnet. \|

\| \*\*REQ-DISC-002\*\* \| Critical \| 1 \| The system must maintain a
registry of discovered peers, removing inactive peers after 3 missed
announcements (15 seconds). \|

\| \*\*REQ-DISC-003\*\* \| High \| 1 \| The CLI must allow a user to
manually specify a peer\'s IP address and port, bypassing the discovery
mechanism. \|

\## 6.2 Connection Requirements

\| ID \| Priority \| Phase \| Requirement Description \|

\|\-\--\|\-\--\|\-\--\|\-\--\|

\| \*\*REQ-CONN-001\*\* \| Critical \| 1 \| The system must use TCP for
all data and control transfer on port 9999 (default). \|

\| \*\*REQ-CONN-002\*\* \| High \| 2 \| The TCP connection must be
configured with the \`TCP_NODELAY\` socket option for low-latency
transmission. \|

\| \*\*REQ-CONN-003\*\* \| High \| 2 \| TCP socket buffers
(\`SO_RCVBUF\` and \`SO_SNDBUF\`) must be set to \*\*256KB\*\* to
optimize high-speed data flow. \|

\| \*\*REQ-CONN-004\*\* \| High \| 2 \| If the connection is idle for 30
seconds, a keepalive must be sent. If no response, the connection must
attempt a re-connect up to 3 times before failing. \|

\## 6.3 Transfer Engine Requirements

\| ID \| Priority \| Phase \| Requirement Description \|

\|\-\--\|\-\--\|\-\--\|\-\--\|

\| \*\*REQ-TRAN-001\*\* \| Critical \| 2 \| Files must be split into
chunks of \*\*4MB\*\* (default). This size must be user-configurable. \|

\| \*\*REQ-TRAN-002\*\* \| Critical \| 2 \| The transfer must begin with
a sequential transfer of metadata, then switch to a parallel transfer
mode using a worker pool (default size 4). \|

\| \*\*REQ-TRAN-003\*\* \| Critical \| 2 \| The system must persist the
transfer state (metadata, completed chunks) to a local JSON file
(\`.p2pxfer/sessions/\<session_id\>.json\`). \|

\| \*\*REQ-TRAN-004\*\* \| High \| 2 \| The \`resume\` command must load
the state file (REQ-TRAN-003) and restart the transfer from the last
successfully acknowledged chunk. \|

\| \*\*REQ-TRAN-005\*\* \| High \| 2 \| The \`pause\` and \`cancel\`
commands must cleanly terminate I/O operations and update the transfer
state file immediately. \|

\## 6.4 Integrity Requirements

\| ID \| Priority \| Phase \| Inputs \| Processing Steps \| Outputs \|
Acceptance Criteria \|

\|\-\--\|\-\--\|\-\--\|\-\--\|\-\--\|\-\--\|\-\--\|

\| \*\*REQ-INTG-001\*\* \| Critical \| 3 \| File chunk (4MB), Chunk
Index (8B). \| Calculate SHA-256 hash of the chunk. Embed hash in the
transfer frame. \| Binary frame containing the chunk and its hash. \|
Hash calculation time \< 10ms per 4MB chunk. \|

\| \*\*REQ-INTG-002\*\* \| Critical \| 3 \| Received chunk, Expected
hash. \| Compare calculated hash of received chunk against the expected
hash. \| ACK or NACK (for retry). \| 100% data integrity verified for
all received chunks. \|

\| \*\*REQ-INTG-003\*\* \| High \| 3 \| NACK received. \| Upon receiving
a NACK, the sender must re-queue the chunk for transmission, up to \*\*3
total retries\*\*. \| Re-queued chunk. \| Transfer must fail gracefully
after 3 failed retries for any single chunk. \|

\| \*\*REQ-INTG-004\*\* \| Critical \| 3 \| Full file path. \| Calculate
and compare the SHA-256 hash of the complete output file against the
initial metadata hash. \| Final success/failure status. \| Transfer is
only declared complete if the full file hash matches. \|

\## 6.5 Security Requirements

\| ID \| Priority \| Phase \| Requirement Description \|

\|\-\--\|\-\--\|\-\--\|\-\--\|

\| \*\*REQ-SEC-001\*\* \| Critical \| 3 \| All chunk payloads must be
encrypted using \*\*AES-256-GCM\*\*. \|

\| \*\*REQ-SEC-002\*\* \| Critical \| 3 \| A unique, non-repeating nonce
must be generated for the encryption of \*\*every single chunk\*\*. \|

\| \*\*REQ-SEC-003\*\* \| Critical \| 3 \| Peer authentication must use
a \*\*4-word passphrase\*\* derived from a 2048-word list, which the
sender and receiver must manually verify via CLI prompt. \|

\| \*\*REQ-SEC-004\*\* \| High \| 3 \| A Diffie-Hellman key exchange
must be performed immediately after TCP connection establishment to
derive the session key. \|

\## 6.6 CLI/User Interface Requirements

\| ID \| Priority \| Phase \| Inputs \| Processing Steps \| Outputs \|
Acceptance Criteria \|

\|\-\--\|\-\--\|\-\--\|\-\--\|\-\--\|\-\--\|\-\--\|

\| \*\*REQ-CLI-001\*\* \| Critical \| 1 \| Transfer progress state. \|
Calculate percentage, MB/s, and ETA. \| Real-time progress bar (e.g.,
\`\[#\-\-\--\] 20% @ 50 MB/s, ETA 00:05:00\`). \| Progress display must
refresh at least 10 times per second (10Hz). \|

\| \*\*REQ-CLI-002\*\* \| High \| 1 \| User command with flags. \|
Utilize Cobra to parse command, validate flags, and execute the core
service call. \| Execution of the requested service
(send/receive/status). \| All required flags (4.1) must be implemented
and validated. \|

\# 7. NON-FUNCTIONAL REQUIREMENTS

\## 7.1 Performance

\* \*\*P-01 Throughput:\*\* On a GigE (1000BASE-T) network, the
sustained transfer rate must exceed \*\*100 MB/s\*\* (800 Mbps).

\* \*\*P-02 Memory:\*\* Peak operational memory consumption (excluding
file buffer) must be less than \*\*32MB\*\*.

\* \*\*P-03 CPU:\*\* Sustained CPU utilization must remain below
\*\*25%\*\* on a modern quad-core processor during peak transfer.

\* \*\*P-04 Latency:\*\* Peer discovery (first announcement to first
response) must complete within \*\*10 seconds\*\*.

\* \*\*P-05 Connect Time:\*\* TCP connection and initial handshake must
complete within \*\*3 seconds\*\*.

\## 7.2 Reliability

\* \*\*R-01 Integrity:\*\* Data integrity must be 100% verifiable
(REQ-INTG-004) under all conditions, including network interruption and
process crashes.

\* \*\*R-02 Crash Recovery:\*\* If the process is forcefully terminated,
the system must be able to recover the state file and successfully
resume the transfer upon next execution.

\* \*\*R-03 Error Handling:\*\* Network errors must utilize an
\*\*exponential backoff\*\* strategy for retries, starting at 1 second
and increasing up to 32 seconds, with a maximum of 3 failed retries per
chunk.

\* \*\*R-04 Error Rate:\*\* The non-recoverable transfer error rate must
be less than \*\*0.01%\*\* of total transfers.

\## 7.3 Usability

\* \*\*U-01 Configuration:\*\* The base system must function with
\*\*zero configuration\*\* required by the user.

\* \*\*U-02 Commands:\*\* Core functionality must be accessible via
simple, intuitive CLI commands (\`p2pxfer send file.iso\`).

\* \*\*U-03 Errors:\*\* All error messages must be clear, actionable,
and categorized (see Section 11).

\## 7.4 Compatibility

\* \*\*C-01 Cross-OS:\*\* Must be fully functional on all supported OS
platforms (2.4).

\* \*\*C-02 Network Agnostic:\*\* Must function equally well over Wi-Fi,
Ethernet, and future WiFi Direct implementations.

\* \*\*C-03 Protocol Versioning:\*\* The protocol must be versioned to
ensure backward compatibility with future versions.

\# 8. TECHNICAL REQUIREMENTS

\| Area \| Requirement \|

\|\-\--\|\-\--\|

\| \*\*Language/Runtime\*\* \| Go 1.21+ \|

\| \*\*CLI Framework\*\* \| Cobra \|

\| \*\*UI\*\* \| BubbleTea for rich terminal UI (e.g., progress bar). \|

\| \*\*Compression\*\* \| zstd library. \|

\| \*\*Crypto\*\* \| Go standard library (\`crypto/tls\`,
\`crypto/aes\`, \`crypto/sha256\`). \|

\| \*\*Sockets\*\* \| Use of \`net.Listen\` and \`net.Dial\` with
\`TCP_NODELAY\`, \`SO_RCVBUF\`, and \`SO_SNDBUF\` (256KB). \|

\| \*\*Concurrency\*\* \| Extensive use of Goroutines,
\`sync.WaitGroup\`, \`sync.Mutex\`, channels, and \`context.Context\`
for cancellation/timeouts. \|

\| \*\*I/O\*\* \| Buffered I/O with \*\*256KB\*\* buffers to minimize
system calls. \|

\| \*\*Logging\*\* \| Structured logging using \`zerolog\`. \|

\# 9. PROTOCOL SPECIFICATION

\## 9.1 Discovery Protocol (UDP)

\* \*\*Magic Number:\*\* All discovery messages must begin with the
4-byte magic number \`0x50325078\` (ASCII for \'P2Px\').

\* \*\*Types:\*\*

\* \`ANNOUNCE\`: Broadcast message indicating the sender is ready to
receive.

\* \`QUERY\`: Broadcast message asking for available peers.

\* \`RESPONSE\`: Unicast response to a \`QUERY\`.

\* \`GOODBYE\`: Broadcast message before shutting down.

\* \*\*Payload:\*\* All data is JSON-encoded and includes protocol
version, peer ID, listening TCP port, and optional status (e.g.,
\`is_transferring: true\`).

\## 9.2 Transfer Protocol (TCP)

The transfer sequence is a state machine:

1\. \*\*TCP Connect:\*\* Sender initiates connection.

2\. \*\*HELLO:\*\* Both nodes exchange protocol versions and
capabilities (e.g., encryption support).

3\. \*\*AUTH:\*\* Diffie-Hellman key exchange, followed by passphrase
verification (REQ-SEC-003).

4\. \*\*METADATA:\*\* Sender transmits file info (name, size, full
SHA-256 hash, chunk count).

5\. \*\*CHUNKS:\*\* Parallel transfer of binary data frames.

6\. \*\*VERIFY:\*\* Final full file hash check.

7\. \*\*CLOSE:\*\* Connection terminates.

\## 9.3 Binary Data Frame

The payload for chunk transfer is encapsulated in a fixed-structure
binary frame:

\| Field \| Size (Bytes) \| Description \|

\|\-\--\|\-\--\|\-\--\|

\| Type \| 2B \| Message type (e.g., DATA, ACK, NACK, METADATA) \|

\| Chunk Index \| 8B \| Index of the chunk (0 to N). \|

\| Length \| 4B \| Length of the payload (max 4MB + Overhead). \|

\| Flags \| 2B \| Bitmask: Compressed (0x01), Encrypted (0x02), Last
Chunk (0x04). \|

\| \*\*Payload\*\* \| Variable \| Encrypted chunk data (including GCM
tag) + SHA-256 hash. \|

\# 10. DATA STRUCTURES (Go)

\`\`\`go

// Peer represents a discovered P2Pxfer instance

type Peer struct {

ID string \`json:\"id\"\`

IP string \`json:\"ip\"\`

Port int \`json:\"port\"\`

LastSeen time.Time \`json:\"last_seen\"\`

IsBusy bool \`json:\"is_busy\"\`

}

// Session holds the runtime state of a transfer

type Session struct {

SessionID string

Peer Peer

Direction string // \"send\" or \"receive\"

File FileInfo

TransferState TransferState

ActiveWorkers sync.WaitGroup

EncryptionKey \[\]byte

}

// FileInfo holds static file properties

type FileInfo struct {

Name string

Size int64

FullHash \[\]byte // SHA-256 of entire file

ChunkSize int

ChunkCount int

}

// ChunkInfo records the status of a single chunk

type ChunkInfo struct {

Index int

Status string // \"pending\", \"transferring\", \"complete\", \"failed\"

Hash \[\]byte // SHA-256 of the unencrypted chunk

Offset int64

Length int

Retries int

}

// TransferState is the structure persisted for resume functionality

type TransferState struct {

SessionID string

LastUpdated time.Time

Chunks \[\]ChunkInfo // List of all chunks and their statuses

}

// Config represents the settings file (4.2)

type Config struct {

Network struct {

DiscoveryPort int

TransferPort int

}

Transfer struct {

DefaultChunkSize int

ParallelWorkers int

}

// \... other config fields

}

11\. ERROR HANDLING AND RECOVERY

Errors must be categorized and assigned unique codes for easy debugging
and user comprehension.

  ------------------------------------------------------------------------------
  Category             Code Range          Description       Recovery Strategy
  -------------------- ------------------- ----------------- -------------------
  **Network**          NET001-NET004       Connection        Exponential backoff
                                           failures,         retry (up to 3x).
                                           timeouts, address 
                                           errors.           

  **File**             FILE001-FILE005     I/O errors, file  Immediate halt,
                                           not found,        user intervention
                                           permission        required.
                                           denied, disk      
                                           full.             

  **Protocol**         PROTO001-PROTO003   Invalid protocol  Close connection,
                                           frame, unexpected log error, attempt
                                           message type.     reconnect if
                                                             appropriate.

  **Integrity**        CHUNK001-CHUNK003   Hash mismatch,    **3-retry max per
                                           incomplete chunk, chunk
                                           frame corruption. (REQ-INTG-003)**.

  **Authentication**   AUTH001-AUTH002     Failed passphrase Close connection,
                                           verification,     do not retry,
                                           invalid key       inform user of
                                           exchange.         security failure.

  **Session**          SESS001-SESS002     Failed to         Log critical error,
                                           load/save state   halt.
                                           file, invalid     
                                           session ID.       
  ------------------------------------------------------------------------------

12\. SECURITY REQUIREMENTS

# 12.1 Threat Model

- **Eavesdropping:** An attacker monitors network traffic (LAN access).

- **Man-in-the-Middle (MITM):** An attacker intercepts communication and
  attempts to relay/modify data.

- **Denial of Service (DoS):** An attacker floods the discovery or
  transfer ports.

# 12.2 Mitigations

1.  **Confidentiality:** All data is protected via **AES-256-GCM**
    (REQ-SEC-001). This prevents eavesdropping and provides
    authentication against payload modification.

2.  **Authentication:** The 4-word verification code (REQ-SEC-003)
    combined with the DH key exchange (REQ-SEC-004) mitigates MITM
    attacks. The user-verified code ensures the peer identity.

3.  **Integrity:** SHA-256 per chunk (REQ-INTG-001) and the GCM
    authentication tag ensure that data corruption or modification is
    immediately detected and rejected.

4.  **DoS:** Implement basic **rate limiting** on the TCP listening port
    to mitigate simple connection-flood attacks. Implement input
    validation on all received metadata to prevent buffer overruns or
    injection attacks.

13\. TESTING REQUIREMENTS

# 13.1 Unit Testing

- Coverage must exceed **90%** for core logic modules: chunking,
  hashing, encryption, and protocol state machine.

# 13.2 Integration Testing

- **End-to-End:** Full send/receive transfer on the same OS, cross-OS,
  and across different network types (Wi-Fi/Ethernet).

- **Interruption Recovery:** Simulate network cable pull, process
  termination, and power loss during transfer to verify successful
  resume capability.

- **Security:** Verify that tampering with data frames results in
  hash/GCM tag failures and correct error reporting.

# 13.3 Performance Testing

- **Throughput Benchmarks:** Conduct tests with 1GB, 10GB, and
  **100GB+** files to verify the 100MB/s throughput target (P-01).

- **Stress Test:** Execute multiple concurrent transfers (if F9.1 is
  implemented) and sustained 100GB+ transfers to monitor CPU (P-03) and
  Memory (P-02) usage over time.

14\. DEPLOYMENT

- **Artifacts:** The system must be cross-compiled into single binary
  executables for:

  - Linux (amd64, arm64)

  - macOS (amd64, arm64)

  - Windows (amd64)

- **Distribution:** Distribute through compressed archives (tar.gz,
  zip), a Windows MSI installer, and a Homebrew tap for macOS/Linux.

- **Versioning:** Use **Semantic Versioning (SemVer)**
  (MAJOR.MINOR.PATCH) for the software binary. The protocol must have a
  separate, internal protocol version number.

15\. APPENDICES

# 15.1 Glossary

(See Section 1.3)

# 15.2 Study Projects

Projects used for conceptual inspiration:

- **Syncthing:** Decentralized, P2P architecture.

- **croc:** CLI-first user experience and single-binary design.

- **Magic Wormhole:** Secure key exchange and verification code pattern.

- **rsync:** Delta-sync/optimization concepts.

# 15.3 Phase Timeline

  -----------------------------------------------------------------------
  Phase                   Duration                Focus
  ----------------------- ----------------------- -----------------------
  1                       2 Weeks                 INTRODUCTION, CLI,
                                                  Discovery (F1), Basic
                                                  Connection (F2).

  2                       2 Weeks                 Transfer Engine (F3),
                                                  Transfer Management
                                                  (F4), TCP Optimization
                                                  (REQ-CONN-002,
                                                  REQ-CONN-003).

  3                       2 Weeks                 Security (F6),
                                                  Integrity (F5),
                                                  Protocol
                                                  Implementation,
                                                  Testing.

  4                       Ongoing                 Compression (F7),
                                                  Configuration,
                                                  Documentation, Polish.

  5                       Future                  Swarm (F9.1), QR
                                                  Pairing (F9.2),
                                                  Advanced Features.
  -----------------------------------------------------------------------

# 15.4 Revision History

  -----------------------------------------------------------------------
  Version           Date              Author            Description of
                                                        Change
  ----------------- ----------------- ----------------- -----------------
  0.1               Date              Person            Initial draft of
                                                        the P2Pxfer SRS.

  -----------------------------------------------------------------------
