## Package `muonengine`

The `muonengine` package provides functionality for a Go project related to a download client for a peer-to-peer network.

### Types

#### Type `bencodeFiles`

```go
type bencodeFiles struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}
```

The `bencodeFiles` struct represents the files in the torrent.

##### Fields

-   `Length`: The length of the file.
-   `Path`: The path of the file.

#### Type `bencodeInfo`

```go
type bencodeInfo struct {
	Pieces      string         `bencode:"pieces"`
	PieceLength int            `bencode:"piece length"`
	Length      int            `bencode:"length"`
	Name        string         `bencode:"name"`
	Files       []bencodeFiles `bencode:"files"`
}
```

The `bencodeInfo` struct represents the information about the torrent.

##### Fields

-   `Pieces`: The raw binary representation of the piece hashes.
-   `PieceLength`: The length of each piece.
-   `Length`: The total length of the data.
-   `Name`: The name of the torrent.
-   `Files`: The files in the torrent.

#### Type `bencodeTorrent`

```go
type bencodeTorrent struct {
	Announce     string        `bencode:"announce"`
	AnnounceList [][]string    `bencode:"announce-list"`
	Info         bencodeInfo   `bencode:"info"`
}
```

The `bencodeTorrent` struct represents the bencoded torrent file.

##### Fields

-   `Announce`: The tracker URL.
-   `AnnounceList`: The list of tracker URLs.
-   `Info`: The information about the torrent.

#### Type `bencodeInfo_v1`

```go
type bencodeInfo_v1 struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}
```

The `bencodeInfo_v1` struct represents the version 1 information about the torrent.

##### Fields

-   `Pieces`: The raw binary representation of the piece hashes.
-   `PieceLength`: The length of each piece.
-   `Length`: The total length of the data.
-   `Name`: The name of the torrent.

#### Type `bencodeInfo_v2`

```go
type bencodeInfo_v2 struct {
	Pieces      string         `bencode:"pieces"`
	PieceLength int            `bencode:"piece length"`
	Files       []bencodeFiles `bencode:"files"`
	Name        string         `bencode:"name"`
}
```

The `bencodeInfo_v2` struct represents the version 2 information about the torrent.

##### Fields

-   `Pieces`: The raw binary representation of the piece hashes.
-   `PieceLength`: The length of each piece.
-   `Files`: The files in the torrent.
-   `Name`: The name of the torrent.

#### Type `TorrentFile`

```go
type TorrentFile struct {
	Announce    []string
	InfoHash    [20]byte
	PieceHash   [][20]byte
	PieceLength int
	Length      int
	Name        string
	Files       []bencodeFiles
}
```

The `TorrentFile` struct represents a torrent file.

##### Fields

-   `Announce`: The tracker URLs.
-   `InfoHash`: The SHA1 hash of the torrent information.
-   `PieceHash`: The SHA1 hashes of the pieces.
-   `PieceLength`: The length of each piece.
-   `Length`: The total length of the data.
-   `Name`: The name of the torrent.
-   `Files`: The files in the torrent.

### Functions

#### Function `Openfile`

```go
func Openfile(r io.Reader) (*bencodeTorrent, error)
```

The `Openfile` function opens a torrent file from an `io.Reader` and returns the parsed `bencodeTorrent` object.

##### Parameters

-   `r`: The `io.Reader` to read the torrent file from.

##### Returns

-   A pointer to the parsed `bencodeTorrent` object.
-   An error if parsing the file fails.

#### Function `Open`

```go
func Open(s string) (*TorrentFile, error)
```

The `Open` function opens a torrent file specified by the file path and returns the parsed `TorrentFile` object.

##### Parameters

-   `s`: The file path of the torrent file.

##### Returns

-   A pointer to the parsed `TorrentFile` object.
-   An error if opening or parsing the file fails.

#### Function `becToTorrent`

```go
func becToTorrent(beTorrent *bencodeTorrent) (TorrentFile, error)
```

The `becToTorrent` function converts a `bencodeTorrent` object to a `TorrentFile` object.

##### Parameters

-   `beTorrent`: The `bencodeTorrent` object to convert.

##### Returns

-   The converted `TorrentFile` object.
-   An error if conversion fails.

#### Method `becTo_v1` (bencodeTorrent)

```go
func (t *bencodeTorrent) becTo_v1() bencodeInfo_v1
```

The `becTo_v1` method converts the `bencodeTorrent` object to a `bencodeInfo_v1` object.

##### Returns

-   The converted `bencodeInfo_v1` object.

#### Method `becTo_v2` (bencodeTorrent)

```go
func (t *bencodeTorrent) becTo_v2() bencodeInfo_v2
```

The `becTo_v2` method converts the `bencodeTorrent` object to a `bencodeInfo_v2` object.

#### Returns

-   The converted `bencodeInfo_v2` object.

#### Method `hash_v1` (bencodeInfo_v1)

```go
func (i *bencodeInfo_v1) hash_v1() ([20]byte, error)
```

The `hash_v1` method calculates the SHA1 hash of the `bencodeInfo_v1` object.

#### Returns

-The calculated SHA1 hash.
-An error if the hash calculation fails.

#### Method `hash_v2` (bencodeInfo_v2)

```go
func (i *bencodeInfo_v2) hash_v2() ([20]byte, error)
```

The `hash_v2` method calculates the SHA1 hash of the `bencodeInfo_v2` object.

#### Returns

-   The calculated SHA1 hash.
-   An error if the hash calculation fails.

#### Method `splitPieceHashes` (bencodeInfo)

```go
func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error)
```

The `splitPieceHashes` method splits the raw binary representation of piece hashes into individual hashes.

#### Returns

-   A slice of `[20]byte` representing the individual piece hashes.
-   An error if the splitting fails.
