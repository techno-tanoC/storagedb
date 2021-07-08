package main

import (
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

	log.Print("Downloading...")
	err := downloadFile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		dbmap, err := initDb()
		if err != nil {
			log.Fatal(err)
		}
		defer dbmap.Db.Close()

		count, err := dbmap.SelectInt("select count from counter")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("count: %d", count)

		log.Print("incrementing...")
		_, err = dbmap.Exec("update counter set count = count + 1")
		if err != nil {
			log.Fatal(err)
		}

		count, err = dbmap.SelectInt("select count from counter")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("new count: %d", count)
	}()

	<-sigs

	fmt.Println("")
	log.Print("Uploading...")
	err = uploadFile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Shutdown!")
}
