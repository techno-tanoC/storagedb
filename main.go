package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	ctx := context.Background()

	err := downloadFile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			func() {
				scanner := bufio.NewScanner(os.Stdin)
				ok := scanner.Scan()
				if !ok {
					log.Fatal("scanner error")
				}
				line := scanner.Text()
				log.Print(line)

				f, err := os.Create(objectPath)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				_, err = fmt.Fprint(f, line+"\n")
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}()

	<-sigs

	log.Print("")
	log.Print("Uploading...")
	err = uploadFile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Shutdown!")
}
