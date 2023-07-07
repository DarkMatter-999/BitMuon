package muonengine

import (
	"fmt"
	"net"
	"time"
)

type downloadClient struct {
	Conn     net.Conn
	Choked   bool
	Bitfield []byte
	peer     Peer
	infoHash [20]byte
	PeerID   [20]byte
}

func newDownloadClient(peer Peer, peerID, infoHash [20]byte) (*downloadClient, error){
	conn, err := net.DialTimeout("tcp", peer.String(), DELAY * time.Second)
	if err != nil {
		return nil, err
	}

	defer conn.SetDeadline(time.Time{})
	_, err = completeHandshake(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}
		
	bf, err := recvBitField(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &downloadClient{
		Conn: conn,
		Choked: true,
		Bitfield: bf,
		peer: peer,
		infoHash: infoHash,
		PeerID: peerID,
	}, nil
}
 
func (c *downloadClient) Read() (*Message, error) {
	msg, err := MessageDeserialize(c.Conn)
	return msg, err
}

func (c *downloadClient) SendInterested() error {
	msg := Message{ID: MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *downloadClient) SendNotInterested() error {
	msg := Message{ID: MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *downloadClient) SendUnchoke() error {
	msg := Message{ID: MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *downloadClient) SendHave(index int) error {
	msg := FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

func (c *downloadClient) SendRequest(index, begin, length int) error {
	req := FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

func startDownloadManager(torr *p2pTorrent) {
	ph := torr.PieceHashes
	for _, peer := range torr.Peers {
		torr.downloadWorker(peer, ph[0])
	}
}

func (t *p2pTorrent) downloadWorker(peer Peer, ph [20]byte) {
	c, err := newDownloadClient(peer, t.PeerID, t.InfoHash)
	if err != nil {
		fmt.Printf("Could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}
	defer c.Conn.Close()
	fmt.Printf("Completed handshake with %s\n", peer.IP)

	c.SendUnchoke()
	c.SendInterested()

	fmt.Println(ph)
}

