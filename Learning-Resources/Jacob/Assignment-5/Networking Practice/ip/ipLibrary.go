/* IP

Provide ip v4 or v6 as argument
*/

package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println(os.Args[1])
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	name := os.Args[1]
	addr := net.ParseIP(name)

	mask := addr.DefaultMask()
	network := addr.Mask(mask)
	ones, bits := mask.Size()

	if addr == nil {
		fmt.Println("Invalid address")
	} else {
		fmt.Println("Address is ", addr.String(),
			" Default mask length is ", bits,
			"Leading ones count is ", ones,
			"Mask is (hex) ", mask.String(),
			" Network is ", network.String())
	}
	os.Exit(0)
}
