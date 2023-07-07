package main

import (
	"bitmuon/muonengine"
	"fmt"
)

func main() {
	var torr, _ = muonengine.Open("file.torrent")
	fmt.Println(torr.Name)

	err := torr.DownloadTorrent()
	if err != nil {
		return
	}

	fmt.Println(torr.InfoHash)
}
