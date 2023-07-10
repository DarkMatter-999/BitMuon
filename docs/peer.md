## Package `muonengine`

The `muonengine` package provides functionality for a Go project related to a download client for a peer-to-peer network.

### Constants

-   `PORT`: The port number used for communication (6868).

### Type `Peer`

```go
type Peer struct {
    IP   net.IP
    Port uint16
}
```

The `Peer` struct represents a peer in the peer-to-peer network.

#### Fields

-   `IP`: The IP address of the peer.
-   `Port`: The port number of the peer.

### Method `String` (Peer)

```go
func (p Peer) String() string
```

The `String` method returns the string representation of the peer in the format "IP:Port".

#### Returns

-   The string representation of the peer.

### Type `bencodeTrackerResp`

```go
type bencodeTrackerResp struct {
    Interval int    `bencode:"interval"`
    Peers    string `bencode:"peers"`
}
```

The `bencodeTrackerResp` struct represents the response from the tracker in bencode format.

#### Fields

-   `Interval`: The interval in seconds between tracker requests.
-   `Peers`: The raw binary representation of peers.

### Function `getPeers`

```go
func getPeers(peersBin []byte) ([]Peer, error)
```

The `getPeers` function parses the raw binary representation of peers and returns a slice of `Peer` objects.

#### Parameters

-   `peersBin`: The raw binary representation of peers.

#### Returns

-   A slice of `Peer` objects representing the parsed peers.
-   An error if parsing fails.

### Function `Download`

```go
func Download(torr *TorrentFile) (*p2pTorrent, error)
```

The `Download` function downloads the torrent specified by the `TorrentFile` and returns a `p2pTorrent` object.

#### Parameters

-   `torr`: The `TorrentFile` object representing the torrent to download.

#### Returns

-   A pointer to the downloaded `p2pTorrent` object.
-   An error if downloading or the tracker requests encounter an error.

### Function `requestPeer`

```go
func requestPeer(url string) ([]Peer, error)
```

The `requestPeer` function sends an HTTP GET request to the tracker URL and returns the list of peers obtained from the response.

#### Parameters

-   `url`: The tracker URL to send the request to.

#### Returns

-   A slice of `Peer` objects representing the obtained peers.
-   An error if the request or parsing fails.

### Function `requestPeerUDP`

```go
func requestPeerUDP(t *TorrentFile, peerId [20]byte, announceUrl string) ([]Peer, error)
```

The `requestPeerUDP` function sends a UDP request to the tracker and returns the list of peers obtained from the response.

#### Parameters

-   `t`: The `TorrentFile` object representing the torrent.
-   `peerId`: The ID of the client.
-   `announceUrl`: The announce URL of the tracker.

#### Returns

-   A slice of `Peer` objects representing the obtained peers.
-   An error if the request or parsing fails.
