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

	file, err := os.Create(os.Args[1] + "logFile.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetPrefix("LOG: ")
	log.SetFlags(log.Ldate | log.Ltime)

	torr, err := muonengine.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Downloading " + torr.Name)
	fmt.Printf("%v MB\n", torr.Length / (1024 * 1024))

	err = torr.DownloadTorrent()
	if err != nil {
		log.Fatal(err)
		return
	}

}
