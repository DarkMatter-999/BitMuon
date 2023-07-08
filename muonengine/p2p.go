package muonengine

import (
	"fmt"
	"net"
	"os"
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

func (t *TorrentFile) DownloadTorrent() error {
	torr, err := Download(t)
	if err != nil {
		return err
	}
	
	buf, err := startDownloadManager(torr)
	if err != nil {
		return err
	}
	outFile, err := os.Create(t.Name)
	if err != nil {
		return err
	}
	defer outFile.Close()
	
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}

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
