# Total Commander FTP Password Recovery — Go SDK

Offline decoder for FTP passwords stored by Total Commander in `wcx_ftp.ini`.

Checksum validation present in Total Commander is not implemented here.

## Installation

```bash
go get github.com/PELock/Total-Commander-FTP-Password-Recovery-Go
```

## Usage

```go
package main

import (
	"encoding/hex"
	"fmt"

	totalcommanderftppasswordrecovery "github.com/PELock/Total-Commander-FTP-Password-Recovery-Go"
)

func main() {
	decoder := totalcommanderftppasswordrecovery.New()
	plain, err := decoder.DecryptPassword("00112233445566778899aabbccddeeff")
	if err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(plain))
}
```

Golden vector: `00112233445566778899aabbccddeeff` → `fdf3b350e9b8c5fbe82d478d`

See `examples/` and `TestDecryptPassword`.

## License

Apache-2.0. Copyright Bartosz Wójcik / PELock.
