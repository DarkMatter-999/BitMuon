package muonengine

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string
	Info     bencodeInfo
}

func (i *bencodeInfo) hash() ([20]byte, error) {
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

func Open(s string) (*Torrent, error) {
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

type Torrent struct {
	Announce    string
	InfoHash    [20]byte
	PieceHash   [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func becToTorrent(beTorrent *bencodeTorrent) (Torrent, error) {
	infoHash, err := beTorrent.Info.hash()
	if err != nil {
		return Torrent{}, err
	}

	pieceHash, err2 := beTorrent.Info.splitPieceHashes()
	if err2 != nil {
		return Torrent{}, err
	}

	t := Torrent{
		Announce:    beTorrent.Announce,
		InfoHash:    infoHash,
		PieceHash:   pieceHash,
		PieceLength: beTorrent.Info.PieceLength,
		Length:      beTorrent.Info.Length,
		Name:        beTorrent.Info.Name,
	}

	return t, nil
}
