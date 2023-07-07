package muonengine

import (
	"fmt"
	"io"
	"net"
	"time"
)

type HandShake struct {
	Pstr string
	InfoHash [20]byte
	PeerID [20]byte
}

func NewHandshake(infoHash [20]byte, peerID [20]byte) *HandShake {
	return &HandShake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

func (h *HandShake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	curr := 0
	buf[0] = byte(len(h.Pstr))
	curr = 1
	curr += copy(buf[curr:], []byte(h.Pstr))
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func Deserialize(r io.Reader) (*HandShake, error) {
	length := make([]byte, 1)
	_, err := io.ReadFull(r, length)
	if err != nil {
		return nil, err
	}

	pstrlen := int(length[0])
	if pstrlen == 0 {
		return nil, fmt.Errorf("Found empty Pstr")
	}

	handShakeBuffer := make([]byte, pstrlen+48)
	_, err = io.ReadFull(r, handShakeBuffer)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handShakeBuffer[pstrlen+8 : pstrlen+8+20])
	copy(peerID[:], handShakeBuffer[pstrlen+8+20 :])

	handShake := HandShake{
		Pstr: string(handShakeBuffer[:pstrlen]),
		InfoHash: infoHash,
		PeerID: peerID,
	}

	return &handShake, nil
}

func (t *TorrentFile) DownloadToFile()  error {
	torr, err := Download(t)
	if err != nil {
		return err
	}

	for i:=0; i < len(torr.Peers); i++ {
		conn, err := net.DialTimeout("tcp", torr.Peers[i].String(), 3 * time.Second)
		if err != nil {
			continue
		}
		defer conn.SetDeadline(time.Time{})
		
		req := NewHandshake(torr.InfoHash, torr.PeerID)
		
		_, err = conn.Write(req.Serialize())
		if err != nil {
			continue
			//return err
		}

		res, err := Deserialize(conn)
		if err != nil {
			continue
			//return err
		}

		fmt.Println(*res)
	}

	return nil

}

type p2pTorrent struct {
	Peers       []Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}
