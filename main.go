package main

import (
	"bitmuon/muonengine"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: bitmuon <file.torrent>")
		return
	}
	var torr, _ = muonengine.Open(os.Args[1])
	fmt.Println("Downloading " + torr.Name)
	file, err := os.Create(torr.Name + "logFile.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Ltime)

	err = torr.DownloadTorrent()
	if err != nil {
		return
	}
}
