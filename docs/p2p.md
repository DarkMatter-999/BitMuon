## Package `muonengine`

The `muonengine` package provides functionality for a Go project related to a download client for a peer-to-peer network.

### Constants

-   `DELAY`: The delay value used for network connections (in seconds).

### Type `p2pTorrent`

```go
type p2pTorrent struct {
    Peers       []Peer
    PeerID      [20]byte
    InfoHash    [20]byte
    PieceHashes [][20]byte
    PieceLength int
    Length      int
    Name        string
    Files       []bencodeFiles
}
```

The `p2pTorrent` struct represents a torrent with information necessary for downloading.

#### Fields

-   `Peers`: The list of peers.
-   `PeerID`: The ID of the client.
-   `InfoHash`: The info hash of the torrent.
-   `PieceHashes`: The list of piece hashes.
-   `PieceLength`: The length of each piece.
-   `Length`: The total length of the torrent.
-   `Name`: The name of the torrent.
-   `Files`: The list of files in the torrent.

### Method `DownloadTorrent` (TorrentFile)

```go
func (t *TorrentFile) DownloadTorrent() error
```

The `DownloadTorrent` method downloads the torrent specified by the `TorrentFile`.

#### Returns

-   An error if downloading or the download manager encounters an error.

### Function `recvBitField`

```go
func recvBitField(conn net.Conn) (Bitfield, error)
```

The `recvBitField` function receives a bitfield message from a network connection.

#### Parameters

-   `conn`: The network connection to receive the bitfield from.

#### Returns

-   The received bitfield as a `Bitfield`.
-   An error if receiving or parsing fails.
