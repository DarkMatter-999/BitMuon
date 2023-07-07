package muonengine

import (
	"bytes"
	"fmt"
	"io"
	"net"
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

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*HandShake, error) {	
		req := NewHandshake(infoHash, peerID)
		
		_, err := conn.Write(req.Serialize())
		if err != nil {
			return nil, err
		}

		res, err := Deserialize(conn)
		if err != nil {
			return nil, err
		}

		if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
			return nil, fmt.Errorf("Expected infohash %x but got %x", res.InfoHash, infoHash)
		}

		return res, nil
}
