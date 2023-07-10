## Package `muonengine`

The `muonengine` package provides functionality for a Go project related to a download client for a peer-to-peer network.

### Type `HandShake`

```go
type HandShake struct {
    Pstr     string
    InfoHash [20]byte
    PeerID   [20]byte
}
```

The `HandShake` struct represents a handshake message used in the BitTorrent protocol.

#### Fields

-   `Pstr`: The protocol string.
-   `InfoHash`: The info hash of the torrent.
-   `PeerID`: The ID of the peer.

### Function `NewHandshake`

```go
func NewHandshake(infoHash [20]byte, peerID [20]byte) *HandShake
```

The `NewHandshake` function creates a new handshake message with the provided info hash and peer ID.

#### Parameters

-   `infoHash`: The info hash of the torrent.
-   `peerID`: The ID of the peer.

#### Returns

-   A pointer to the created `HandShake` message.

### Method `Serialize` (HandShake)

```go
func (h *HandShake) Serialize() []byte
```

The `Serialize` method serializes the handshake message into a byte slice.

#### Returns

-   The serialized handshake message as a byte slice.

### Function `Deserialize`

```go
func Deserialize(r io.Reader) (*HandShake, error)
```

The `Deserialize` function deserializes a handshake message from an `io.Reader` and returns a `HandShake` object.

#### Parameters

-   `r`: The `io.Reader` to read the handshake message from.

#### Returns

-   A pointer to the deserialized `HandShake` object.
-   An error if deserialization fails.

### Function `completeHandshake`

```go
func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*HandShake, error)
```

The `completeHandshake` function performs a complete handshake with a peer.

#### Parameters

-   `conn`: The network connection to the peer.
-   `infoHash`: The info hash of the torrent.
-   `peerID`: The ID of the peer.

#### Returns

-   A pointer to the received `HandShake` object.
-   An error if the handshake fails or the received info hash does not match the expected info hash.
