// +build !linux

package main

import "log"

func main() {
	log.Fatal("mailenc is only available for linux for now.")
}
