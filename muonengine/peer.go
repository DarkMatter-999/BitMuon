package muonengine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

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
	var peers []Peer
	_, err := rand.Read(peerId[:])
	if err != nil {
		return nil, err
	}

	for _, announce := range torr.Announce {
		log.Printf("Trying Tracker: %v", announce)
		fmt.Printf("Trying Tracker: %v \n", announce)

		if announce[:4] == "http" {
			url, err := torr.BuildTrackerURL(peerId, PORT, announce)
			if err != nil {
				log.Println(err)
				continue
			}

			newpeers, err := requestPeer(url)
			if err != nil {
				log.Println(err)
				continue
			}

			peers = append(peers, newpeers...)

		} else if announce[:3] == "udp" {
			newpeers, err := requestPeerUDP(torr, peerId, announce)
			if err != nil {
				log.Println(err)
				continue
			}

			peers = append(peers, newpeers...)
		}
	}

	p2ptorr := p2pTorrent{
		Peers: peers,
		PeerID: peerId,
		InfoHash: torr.InfoHash,
		PieceHashes: torr.PieceHash,
		PieceLength: torr.PieceLength,
		Length: torr.Length,
		Name: torr.Name,
		Files: torr.Files,
	}
	log.Printf("Got peers : %v", peers)

	return &p2ptorr, nil
}

func requestPeer (url string) ([]Peer, error){
	log.Printf("Using URL: %s", url)
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

func requestPeerUDP(t *TorrentFile, peerId [20]byte, announceUrl string) ([]Peer, error) {
	log.Printf("Using URL: %s", announceUrl)
	url, err := t.GetHostPortUDP(announceUrl)
	if err != nil {
		return nil, err
	}

	sock, err := net.DialTimeout("udp", url.Host, 3 * time.Second)
	if err != nil {
		return nil, err
	}

	defer sock.Close()

	sock.SetDeadline(time.Now().Add(3 * time.Second))
	sock.SetDeadline(time.Time{})
	
	transactionID := rand.Uint32()

	// Generate Connect request
	protocolID := uint64(0x41727101980)
	action := uint32(0) // Connect action
	packet := make([]byte, 16)

	binary.BigEndian.PutUint64(packet[:8], protocolID)
	binary.BigEndian.PutUint32(packet[8:12], action)
	binary.BigEndian.PutUint32(packet[12:16], transactionID)

	_, err = sock.Write(packet)
	if err != nil {
		sock.Close()
		return nil, err
	}
	sock.SetDeadline(time.Now().Add(3 * time.Second))

	response := make([]byte, 16)

	_, err = sock.Read(response)
	if err != nil {
		sock.Close()
		return nil, err
	}
	sock.SetDeadline(time.Now().Add(3 * time.Second))
	sock.SetDeadline(time.Time{})

	action = binary.BigEndian.Uint32(response[0:4])
	receivedTransactionID := binary.BigEndian.Uint32(response[4:8])
	connectionID := binary.BigEndian.Uint64(response[8:16])

	log.Printf("Got Connection ID %v", connectionID)

	if action != 0 || receivedTransactionID != transactionID {
		sock.Close()
		return nil, fmt.Errorf("Invalid connect response.")
	}

/*
	sock.Close()

	sock, err = net.Dial("udp", url.Host)
	if err != nil {
		return nil, err
	}
	defer sock.Close()

*/
	port, err := strconv.ParseUint(url.Port(), 10, 16)	
	if err != nil {
		sock.Close()
		return nil, err
	}

	// Generate Announce request
	action = uint32(1) // Announce action
	packet = make([]byte, 98)
	event := uint64(0)
	binary.BigEndian.PutUint64(packet[:8], connectionID)
	binary.BigEndian.PutUint32(packet[8:12], action)
	binary.BigEndian.PutUint32(packet[12:16], transactionID)
	copy(packet[16:36], t.InfoHash[:])
	copy(packet[36:56], peerId[:])
	binary.BigEndian.PutUint64(packet[56:64], 0)
	binary.BigEndian.PutUint64(packet[64:72], uint64(t.Length))
	binary.BigEndian.PutUint64(packet[72:80], 0)
	binary.BigEndian.PutUint32(packet[80:84], uint32(event))
	binary.BigEndian.PutUint32(packet[84:88], 0)
	binary.BigEndian.PutUint32(packet[88:92], 50)
	binary.BigEndian.PutUint32(packet[92:96], ^uint32(0))
	binary.BigEndian.PutUint16(packet[96:98], uint16(port))

	_, err =	sock.Write(packet)
	if err != nil {
		sock.Close()
		return nil, err
	}

	response = make([]byte, 4096)
	_, err = sock.Read(response)
	if err != nil {
		sock.Close()
		return nil, err
	}

	_ = binary.BigEndian.Uint32(response[0:4]) // Recieved Action
	receivedTransactionID = binary.BigEndian.Uint32(response[4:8])
	if receivedTransactionID == transactionID {
		offset := 20
		peers := make([]Peer, 0)
		for offset+6 < len(response) {
			if !bytes.Equal(response[offset:offset+4], []byte{0,0,0,0}) {
				ip := net.IP(response[offset : offset+4])
				port := binary.BigEndian.Uint16(response[offset+4 : offset+6])
				peer := Peer{IP: ip, Port: port}
				peers = append(peers, peer)
			}
			offset += 6
		}

		return peers, nil
	}


	sock.Close()

	return nil, fmt.Errorf("Connection closed by %v", url)
}
