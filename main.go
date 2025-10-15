package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var wg sync.WaitGroup

// write function sends HTTP GET requests to a specified URL every second
func write(ctx context.Context, ch chan<- []byte) {
	defer close(ch)
	defer wg.Done()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("write: Context canceled")
			return
		case <-ticker.C:
			resp, err := http.Get("http://138.201.177.104:3040/ping")
			if err != nil {
				log.Println("write: HTTP request error:", err)
				return
			}

			func() {
				defer resp.Body.Close()
				var buffer bytes.Buffer
				if _, err := io.Copy(&buffer, resp.Body); err != nil {
					log.Println("write: Error copying body:", err)
					return
				}

				ch <- buffer.Bytes()
			}()
		}
	}
}

// reader function processes data received from the write function
func reader(ctx context.Context, ch <-chan []byte, db *leveldb.DB) {
	defer wg.Done()

	file, err := os.OpenFile("Number.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("reader: File open error:", err)
		return
	}
	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("reader: Context canceled")
			return
		case num, ok := <-ch:
			if !ok {
				log.Println("reader: Channel closed")
				return
			}

			numberStr := string(num)
			nint, err := strconv.Atoi(numberStr)
			if err != nil {
				log.Println("reader: Cannot convert bytes to int:", err)
				continue
			}
			fmt.Println(nint)

			if _, err := fmt.Fprintln(file, nint); err != nil {
				log.Println("reader: Error writing to file:", err)
			}

			key := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
			if err := db.Put(key, num, nil); err != nil {
				log.Println("reader: Error writing to DB:", err)
			} else {
				data, err := db.Get(key, nil)
				if err != nil {
					log.Println("reader: Error reading from DB:", err)
				} else {
					fmt.Printf("reader: Value: %s, Key: %s\n", data, key)
				}
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan []byte)

	db, err := leveldb.OpenFile("start-goroutines", nil)
	if err != nil {
		log.Fatal("main: Error opening LevelDB:", err)
	}
	defer db.Close()

	wg.Add(2)
	go write(ctx, ch)
	go reader(ctx, ch, db)

	wg.Wait()
}

// The program will run indefinitely until manually terminated.
