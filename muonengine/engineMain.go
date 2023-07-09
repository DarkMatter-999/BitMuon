package muonengine

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
)

func (t *TorrentFile) BuildTrackerURL(peerID [20]byte, port uint16, announceUrl string) (string, error) {
	base, err := url.Parse(announceUrl)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash": []string{string(t.InfoHash[:])},
		"peer_id": []string{string(peerID[:])},
		"port": []string{string(strconv.Itoa(int(port)))},
		"uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) GetHostPortUDP(announceUrl string) (*url.URL, error) {
	host, err := url.ParseRequestURI(announceUrl)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func ParseMagnetLink(magnetLink string) (*TorrentFile, error) {
	torrent := &TorrentFile{}

	u, err := url.Parse(magnetLink)
	if err != nil {
		return nil, err
	}

	infoHash, err := hex.DecodeString(u.Query().Get("xt")[9:])
	if err != nil {

		return nil, err
	}
	copy(torrent.InfoHash[:], infoHash)

	torrent.Announce = []string{u.Query().Get("tr")}
	torrent.Name = u.Query().Get("dn")
	torrent.PieceLength = 0 // Placeholder for the piece length
	torrent.Length = 0      // Placeholder for the total length

	fmt.Println(torrent)
	/*
	// Retrieve the hash of each piece
	pieces := u.Query().Get("xt")[20:]
	if len(pieces)%40 != 0 {
		return nil, fmt.Errorf("invalid magnet link")
	}
	numPieces := len(pieces) / 40
	torrent.PieceHash = make([][20]byte, numPieces)
	for i := 0; i < numPieces; i++ {
		piece, err := hex.DecodeString(pieces[i*40 : (i+1)*40])
		if err != nil {
			return nil, err
		}
		copy(torrent.PieceHash[i][:], piece)
	}
	*/

	return torrent, nil
}

