package muonengine

import (
	"fmt"
	"net"
	"time"
)

const DELAY = 3

type p2pTorrent struct {
	Peers       []Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (t *TorrentFile) DownloadTorrent()  error {
	torr, err := Download(t)
	if err != nil {
		return err
	}
	
	startDownloadManager(torr)	

	return nil

}

func recvBitField(conn net.Conn) (Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{}) 

	msg, err := MessageDeserialize(conn)
	if err != nil {
		return nil, err
	}

	if msg.ID != MsgBitfield {
		return nil, fmt.Errorf("Expected bitfield but got ID %d", msg.ID)
	}
	return msg.Data, nil
}
