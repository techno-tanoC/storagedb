package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
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
	log.Print("Downloaded!")

	dbmap, err := initDb()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/count", func(c *gin.Context) {
		count, err := count(dbmap)
		if err != nil {
			log.Print(err)
		}

		c.JSON(200, gin.H{
			"count": count,
		})
	})
	r.GET("/increment", func(c *gin.Context) {
		err := increment(dbmap)
		if err != nil {
			log.Print(err)
		}

		c.JSON(204, nil)
	})
	r.GET("/decrement", func(c *gin.Context) {
		err := decrement(dbmap)
		if err != nil {
			log.Print(err)
		}

		c.JSON(204, nil)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-sigs

	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Print("Shutdown...")
	err = srv.Shutdown(c)
	if err != nil {
		log.Fatal(err)
	}

	err = dbmap.Db.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Uploading...")
	err = uploadFile(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Stopped!")
	fmt.Println("")
}

func count(dbmap *gorp.DbMap) (int64, error) {
	count, err := dbmap.SelectInt("select count from counter")
	if err != nil {
		return 0, err
	}

	log.Printf("count: %d", count)
	return count, nil
}

func increment(dbmap *gorp.DbMap) error {
	_, err := dbmap.Exec("update counter set count = count + 1")
	if err != nil {
		return err
	}

	log.Printf("increment")
	return nil
}

func decrement(dbmap *gorp.DbMap) error {
	_, err := dbmap.Exec("update counter set count = count - 1")
	if err != nil {
		return err
	}

	log.Printf("decrement")
	return nil
}
