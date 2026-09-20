/******************************************************************************
 * Total Commander FTP Password Recovery usage example.
 *
 * Version        : v1.0.1
 * Language       : Go
 * Author         : Bartosz Wójcik
 * Web page       : https://www.pelock.com
 *
 *****************************************************************************/

package main

import (
	"encoding/hex"
	"fmt"
	"os"

	totalcommanderftppasswordrecovery "github.com/PELock/Total-Commander-FTP-Password-Recovery-Go"
)

func main() {
	const cipherHex = "00112233445566778899aabbccddeeff"
	decoder := totalcommanderftppasswordrecovery.New()
	plain, err := decoder.DecryptPassword(cipherHex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Decode failed.")
		os.Exit(1)
	}
	fmt.Println(string(plain))
	fmt.Println("hex: " + hex.EncodeToString(plain))
}
