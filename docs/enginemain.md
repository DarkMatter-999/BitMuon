# `muonengine` Package Documentation

## Imports

-   **encoding/hex:** This package provides functions for encoding and decoding hexadecimal strings.
-   **fmt:** This package implements formatted I/O with functions similar to C's printf and scanf.
-   **net/url:** This package provides URL parsing and formatting.

## Functions

### `BuildTrackerURL`

```go
func (t *TorrentFile) BuildTrackerURL(peerID [20]byte, port uint16, announceUrl string) (string, error)
```

The `BuildTrackerURL` function builds and returns a tracker URL for the given torrent file. It takes the following parameters:

-   `peerID [20]byte`: The peer ID used for identifying the client.
-   `port uint16`: The port number on which the client listens for connections.
-   `announceUrl string`: The announce URL of the tracker.

### `GetHostPortUDP`

```go
func (t *TorrentFile) GetHostPortUDP(announceUrl string) (*url.URL, error)
```

The `GetHostPortUDP` function parses the announce URL and returns the host and port information as a `url.URL` object. It takes the following parameter:

-   `announceUrl string`: The announce URL of the tracker.

### `ParseMagnetLink`

```go
func ParseMagnetLink(magnetLink string) (*TorrentFile, error)
```

The `ParseMagnetLink` function parses a magnet link and returns a `TorrentFile` object. It takes the following parameter:

-   `magnetLink string`: The magnet link to parse.

This function retrieves various information from the magnet link, such as the info hash, announce URL, name, and more, and populates a `TorrentFile` object accordingly.
