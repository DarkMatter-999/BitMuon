## Package `muonengine`

The `muonengine` package provides functionality for a Go project related to a download client for a peer-to-peer network.

### Constants

-   `MAX_BLOCK_SIZE`: The maximum block size for downloading.
-   `MAX_BACK_LOG`: The maximum back log value.

### Type `downloadClient`

```go
type downloadClient struct {
    Conn     net.Conn
    Choked   bool
    Bitfield Bitfield
    peer     Peer
    infoHash [20]byte
    PeerID   [20]byte
}
```

The `downloadClient` struct represents a client for downloading files from a peer in the peer-to-peer network.

#### Fields

-   `Conn`: The network connection.
-   `Choked`: A flag indicating whether the client is choked.
-   `Bitfield`: The bitfield indicating the availability of pieces.
-   `peer`: The peer information.
-   `infoHash`: The info hash of the torrent.
-   `PeerID`: The ID of the client.

### Type `pieceWork`

```go
type pieceWork struct {
    index  int
    hash   [20]byte
    length int
}
```

The `pieceWork` struct represents a piece of work for downloading a specific piece from a peer.

#### Fields

-   `index`: The index of the piece.
-   `hash`: The hash of the piece.
-   `length`: The length of the piece.

### Type `pieceResult`

```go
type pieceResult struct {
    index int
    buf   []byte
}
```

The `pieceResult` struct represents the result of downloading a piece from a peer.

#### Fields

-   `index`: The index of the piece.
-   `buf`: The downloaded data buffer.

### Type `pieceStatus`

```go
type pieceStatus struct {
    index      int
    client     *downloadClient
    buf        []byte
    downloaded int
    requested  int
    backlog    int
}
```

The `pieceStatus` struct represents the status of downloading a piece from a peer.

#### Fields

-   `index`: The index of the piece.
-   `client`: The download client.
-   `buf`: The data buffer for the piece.
-   `downloaded`: The amount of data downloaded.
-   `requested`: The amount of data requested.
-   `backlog`: The backlog of requested data.

### Function `newDownloadClient`

```go
func newDownloadClient(peer Peer, peerID, infoHash [20]byte) (*downloadClient, error)
```

The `newDownloadClient` function creates a new download client and establishes a connection with the specified peer.

#### Parameters

-   `peer`: The peer information.
-   `peerID`: The ID of the client.
-   `infoHash`: The info hash of the torrent.

#### Returns

-   A pointer to the created `downloadClient`.
-   An error if the connection or handshake fails.

### Method `Read` (downloadClient)

```go
func (c *downloadClient) Read() (*Message, error)
```

The `Read` method reads a message from the download client's connection.

#### Returns

-   The read `Message`.
-   An error if reading fails.

### Method `SendInterested` (downloadClient)

```go
func (c *downloadClient) SendInterested() error
```

The `SendInterested` method sends an "Interested" message to the peer.

#### Returns

-   An error if sending fails.

### Method `SendNotInterested` (downloadClient)

```go
func (c *downloadClient) SendNotInterested() error
```

The `SendNotInterested` method sends a "Not Interested" message to the peer.

#### Returns

-   An error if sending fails.

### Method `SendUnchoke` (downloadClient)

```go
func (c *downloadClient) SendUnchoke() error
```

The `SendUnchoke` method sends an "Unchoke" message to the peer.

#### Returns

-   An error if sending fails.

### Method `SendHave` (downloadClient)

```go
func (c *downloadClient) SendHave(index int) error
```

The `SendHave` method sends a "Have" message to the peer, indicating that the client has a specific piece.

#### Parameters

-   `index`: The index of the piece.

#### Returns

-   An error if sending fails.

### Method `SendRequest` (downloadClient)

```go
func (c *downloadClient) SendRequest(index, begin, length int) error
```

The `SendRequest` method sends a "Request" message to the peer, requesting a specific block of data from a piece.

#### Parameters

-   `index`: The index of the piece.
-   `begin`: The beginning offset of the requested block.
-   `length`: The length of the requested block.

#### Returns

-   An error if sending fails.

### Method `getPieceBounds` (p2pTorrent)

```go
func (t *p2pTorrent) getPieceBounds(index int) (int, int)
```

The `getPieceBounds` method returns the beginning and ending offsets of a piece in the torrent.

#### Parameters

-   `index`: The index of the piece.

#### Returns

-   The beginning offset.
-   The ending offset.

### Method `getPieceSize` (p2pTorrent)

```go
func (t *p2pTorrent) getPieceSize(index int) int
```

The `getPieceSize` method returns the size of a piece in the torrent.

#### Parameters

-   `index`: The index of the piece.

#### Returns

-   The size of the piece.

### Function `checkShaSum`

```go
func checkShaSum(pw *pieceWork, buf []byte) error
```

The `checkShaSum` function checks the integrity of a downloaded piece by comparing its SHA-1 hash.

#### Parameters

-   `pw`: The pieceWork representing the piece.
-   `buf`: The downloaded data buffer.

#### Returns

-   An error if the integrity check fails.

### Function `startDownloadManager`

```go
func startDownloadManager(torr *p2pTorrent) error
```

The `startDownloadManager` function starts the download manager for the given torrent.

#### Parameters

-   `torr`: The p2pTorrent representing the torrent.

#### Returns

-   An error if the download manager encounters an error.

### Method `downloadWorker` (p2pTorrent)

```go
func (t *p2pTorrent) downloadWorker(peer Peer, workQueue chan *pieceWork, workResult chan *pieceResult)
```

The `downloadWorker` method represents a worker that downloads pieces from a specific peer.

#### Parameters

-   `peer`: The peer to download from.
-   `workQueue`: The channel for receiving piece work.
-   `workResult`: The channel for sending piece results.

### Method `readMessage` (pieceStatus)

```go
func (state *pieceStatus) readMessage() error
```

The `readMessage` method reads a message from the download client and updates the piece status accordingly.

#### Returns

-   An error if reading or parsing the message fails.

### Function `downloadPiece`

```go
func downloadPiece(c *downloadClient, pw *pieceWork) ([]byte, error)
```

The `downloadPiece` function downloads a specific piece from a peer.

### Parameters

-   `c`: The download client.
-   `pw`: The pieceWork representing the piece to download.

### Returns

-   The downloaded data buffer.
-   An error if downloading or integrity checking fails.
