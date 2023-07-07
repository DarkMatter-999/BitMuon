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

type p2pClient struct {
	Conn     net.Conn
	Choked   bool
	Bitfield []byte
	peer     Peer
	infoHash [20]byte
	PeerID   [20]byte
}

func (t *TorrentFile) DownloadTorrent()  error {
	torr, err := Download(t)
	if err != nil {
		return err
	}

	for i:=0; i < len(torr.Peers); i++ {
		conn, err := net.DialTimeout("tcp", torr.Peers[i].String(), DELAY * time.Second)
		if err != nil {
			fmt.Println(err)
			continue
		}

		defer conn.SetDeadline(time.Time{})
		res, err := completeHandshake(conn, torr.InfoHash, torr.PeerID)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		fmt.Println(*res)
	}

	return nil

}

