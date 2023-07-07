package muonengine

import (
	"fmt"
	"net"
	"time"
)

const MAX_BLOCK_SIZE = 16384
const MAX_BACK_LOG = 5

type downloadClient struct {
	Conn     net.Conn
	Choked   bool
	Bitfield Bitfield
	peer     Peer
	infoHash [20]byte
	PeerID   [20]byte
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceStatus struct {
	index      int
	client     *downloadClient
	buf        []byte
	downloaded int
	requested  int
	backlog    int
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

func (t *p2pTorrent) getPieceBounds(index int) (int, int) {
	begin := index * t.PieceLength
	end := begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}

func (t *p2pTorrent) getPieceSize(index int) int {
	being, end := t.getPieceBounds(index)
	return end - being
}

func startDownloadManager(torr *p2pTorrent) {
	workQueue := make(chan *pieceWork, len(torr.PieceHashes))
	workResult := make(chan *pieceResult)

	for idx, hash := range torr.PieceHashes {
		length := torr.getPieceSize(idx)
		workQueue <- &pieceWork{idx, hash, length}
	}

	for _, peer := range torr.Peers {
		torr.downloadWorker(peer, workQueue, workResult)
	}
}

func (t *p2pTorrent) downloadWorker(peer Peer, workQueue chan *pieceWork, workResult chan *pieceResult) {
	c, err := newDownloadClient(peer, t.PeerID, t.InfoHash)
	if err != nil {
		fmt.Printf("Could not handshake with %s. Disconnecting\n", peer.IP)
		return
	}
	defer c.Conn.Close()
	fmt.Printf("Completed handshake with %s\n", peer.IP)

	c.SendUnchoke()
	c.SendInterested()

	for pw := range workQueue {
		if !c.Bitfield.HasPiece(pw.index) {
			workQueue <- pw
			continue
		}

		buf, err := downloadPiece(c, pw)
		if err != nil {
			workQueue <- pw
			return
		}

		fmt.Println(buf[:10])
	}
}

func (state *pieceStatus) readMessage() error {
	msg, err := state.client.Read()
	if err != nil {
		return err
	}

	// keep-alive
	if msg == nil { 
		return nil
	}

	switch msg.ID {
	case MsgUnchoke:
		state.client.Choked = false
	case MsgChoke:
		state.client.Choked = true
	case MsgHave:
		index, err := ParseHave(msg)
		if err != nil {
			return err
		}
		state.client.Bitfield.SetPiece(index)
	case MsgPiece:
		n, err := ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil
}

func downloadPiece(c *downloadClient, pw *pieceWork) ([]byte, error) {
	state := pieceStatus{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}

	c.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer c.Conn.SetDeadline(time.Time{})

	for state.downloaded < pw.length {
		if !state.client.Choked {
			for state.backlog < MAX_BACK_LOG && state.requested < pw.length {
				blockSize := MAX_BLOCK_SIZE
				if pw.length - state.requested < blockSize {
					blockSize = pw.length - state.requested
				}

				err := c.SendRequest(pw.index, state.requested, blockSize)
				if err != nil {
					return nil, err
				}
				state.backlog++
				state.requested += blockSize
			}
		}
		err := state.readMessage()
		if err != nil {
			return nil, err
		}
	}
	return state.buf, nil
}

