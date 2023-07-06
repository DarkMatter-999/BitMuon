package main

import (
	"bitmuon/muonengine"
	"fmt"
)

func main() {
	var torr, _ = muonengine.Open("file.torrent")

	fmt.Println(torr.Name)
}
