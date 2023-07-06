package main

import (
	"bitmuon/muonengine"
	"fmt"
)

func main() {
	var torr, _ = muonengine.Open("file.torrent")
	fmt.Println(torr.Name)

	fmt.Println(torr.BuildTrackerURL([20]byte{45, 84, 82, 50, 57, 52, 48, 45, 107, 56, 104, 106, 48, 119, 103, 101, 106, 54, 99, 104}, 6881))
	fmt.Println(torr.InfoHash)
}
