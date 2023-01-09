package main

import (
	"flag"
	"fmt"
	"os"
	"typora-image-upload/src/upload/imgtp"
)

func after_upload(imageUrls []string) {
	fmt.Println("Upload Success:")
	for _, url := range imageUrls {
		fmt.Println(url)
	}
}

func main() {
	// Parse parses the command-line flags from os.Args[1:]
	flag.Parse()

	// typora give images file path
	image_paths := flag.Args()

	// no image
	if len(image_paths) == 0 {
		os.Exit(1)
	} else {
		// do upload
		//token := imgtp.GetToken()
		imgtp.Try()
	}

	// finish
	after_upload([]string{})
}
