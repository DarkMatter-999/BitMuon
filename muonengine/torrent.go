package muonengine

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jackpal/bencode-go"
)

type bencodeFiles struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	Files		[]bencodeFiles `bencode:"files"`
}

type bencodeTorrent struct {
	Announce		string `bencode:"announce"`
	AnnounceList	[][]string `bencode:"announce-list"`
	Info			bencodeInfo `bencode:"info"`
}

type bencodeInfo_v1 struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeInfo_v2 struct {
	Pieces      string			`bencode:"pieces"`
	PieceLength int				`bencode:"piece length"`
	Files		[]bencodeFiles	`bencode:"files"`
	Name		string			`bencode:"name"`
}

func (t *bencodeTorrent) becTo_v1() (i bencodeInfo_v1) {
	v1info := bencodeInfo_v1 {
		Pieces: t.Info.Pieces,
		PieceLength: t.Info.PieceLength,
		Length: t.Info.Length,
		Name: t.Info.Name,
	}
	return v1info
}

func (t *bencodeTorrent) becTo_v2() (i bencodeInfo_v2) {
	v2info := bencodeInfo_v2 {
		Pieces: t.Info.Pieces,
		PieceLength: t.Info.PieceLength,
		Files: t.Info.Files,
		Name: t.Info.Name,
	}
	return v2info
}

func (i *bencodeInfo_v1) hash_v1() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}

	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo_v2) hash_v2() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i) 
	if err != nil {
		return [20]byte{}, err
	}

	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20

	buf := []byte(i.Pieces)

	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Incorrect piece length of %v", len(buf))
		return nil, err
	}

	numHash := len(buf) / hashLen
	hashes := make([][20]byte, numHash)

	for i := 0; i < numHash; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}

	return hashes, nil

}

func Openfile(r io.Reader) (*bencodeTorrent, error) {
	torr := bencodeTorrent{}
	err := bencode.Unmarshal(r, &torr)
	if err != nil {
		return nil, err
	} else {
		return &torr, nil
	}
}

func Open(s string) (*TorrentFile, error) {
	file, err := os.Open(s)
	if err != nil {
		return nil, err
	}

	betorr, err := Openfile(file)
	if err != nil {
		return nil, err
	}

	torr, err := becToTorrent(betorr)
	if err != nil {
		return nil, err
	}

	return &torr, nil

}

type TorrentFile struct {
	Announce    []string
	InfoHash    [20]byte
	PieceHash   [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func becToTorrent(beTorrent *bencodeTorrent) (TorrentFile, error) {
	var infoHash [20]byte
	length := beTorrent.Info.Length
	if len(beTorrent.Info.Files) == 0 {
		v1info := beTorrent.becTo_v1()
		infohash, err := v1info.hash_v1()
		if err != nil {
			return TorrentFile{}, err
	}
		infoHash = infohash
	} else {
		v2info := beTorrent.becTo_v2()
		infohash, err := v2info.hash_v2()
		if err != nil {
			return TorrentFile{}, err
		}
		infoHash = infohash
		for _, file := range beTorrent.Info.Files {
			length += file.Length
		}
	}

	log.Printf("Got infoHash: %x", infoHash)

	pieceHash, err := beTorrent.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}

	announceArray := []string{beTorrent.Announce}

	for _, url := range beTorrent.AnnounceList {
		announceArray = append(announceArray, url[0])
	}

	t := TorrentFile{
		Announce:    announceArray,
		InfoHash:    infoHash,
		PieceHash:   pieceHash,
		PieceLength: beTorrent.Info.PieceLength,
		Length:      length,
		Name:        beTorrent.Info.Name,
	}

	return t, nil
}
