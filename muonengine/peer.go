package muonengine

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/jackpal/bencode-go"
)

const PORT uint16 = 6868

type Peer struct {
	IP net.IP
	Port uint16
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func getPeers(peersBin []byte) ([]Peer, error) {
	const peerSize = 6

	numPeer := len(peersBin) / peerSize
	if len(peersBin) % peerSize != 0 {
		return nil, fmt.Errorf("Invalid sized peers recieved of length %v", len(peersBin))
	}

	peers := make([]Peer, numPeer)
	
	for i:=0; i<numPeer; i++ {
		offset := i*peerSize
		peers[i].IP = net.IP(peersBin[offset: offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4:offset+6])
	}
	return peers, nil
}

func Download(torr *TorrentFile) (*p2pTorrent, error) {
	var peerId [20]byte;

	_, err := rand.Read(peerId[:])
	if err != nil {
		return nil, err
	}

	url, err := torr.BuildTrackerURL(peerId, PORT)
	if err != nil {
		return nil, err
	}

	peers, err := requestPeer(url)
	if err != nil {
		return nil, err
	}

	p2ptorr := p2pTorrent{
		Peers: peers,
		PeerID: peerId,
		InfoHash: torr.InfoHash,
		PieceHashes: torr.PieceHash,
		PieceLength: torr.PieceLength,
		Length: torr.Length,
		Name: torr.Name,
	}

	return &p2ptorr, nil

}

func requestPeer (url string) ([]Peer, error){
	fmt.Println(url)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Incorrect StatusCode recieved %v", res.StatusCode)
	}

	defer res.Body.Close()
	trackerRes := bencodeTrackerResp{}

	err = bencode.Unmarshal(res.Body, &trackerRes)
	if err != nil {
		return nil, err
	}

	peers, err := getPeers([]byte(trackerRes.Peers))
	if err != nil {
		return nil, err
	}

	return peers, err
}
