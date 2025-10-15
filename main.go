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
func write(ch chan<- []byte, ctx context.Context) {
	ctx, cancle := context.WithCancel(ctx)

	defer close(ch)
	defer wg.Done()
	for {
		select {
		default:
			respons, err := http.Get("http://138.201.177.104:3040/ping")
			if err != nil {
				log.Printf("write: Error occurred:")
				cancle()
				return
			}

			var buffer bytes.Buffer
			_, err = io.Copy(&buffer, respons.Body)
			if err != nil {
				log.Printf("write: Error while copying body:")
				return
			}

			ch <- buffer.Bytes()
			time.Sleep(time.Second * 1)
			defer respons.Body.Close()
		}
	}
}

// reader function processes data received from the write function
func reader(ch <-chan []byte, db *leveldb.DB, ctx context.Context) {
	defer wg.Done()

	file, err := os.OpenFile("Number.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("err : not write in file", err)
		return
	}
	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			log.Printf("reader: Context canceled")
			return
		case num := <-ch:
			number := string(num)
			nint, err := strconv.Atoi(number)
			if err != nil {
				log.Println("not convert bytes to int", err)
				return
			}
			fmt.Println(nint)

			fmt.Fprintln(file, nint)

			key := []byte(strconv.FormatInt(time.Now().Unix(), 10))
			err = db.Put(key, num, nil)
			if err != nil {
				log.Println("not send key & value in the database", err)
				return
			}

			data, err := db.Get(key, nil)
			if err != nil {
				log.Println("reader: Error while getting data from database:", err)
				return
			}
			fmt.Printf("reader: Value: %s, key: %s \n", data, key)
		}
	}
}

// main function initializes the program
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan []byte)

	db, err := leveldb.OpenFile("start-gorotins", nil)
	if err != nil {
		log.Println("err in open library", err)
		return
	}
	defer db.Close()

	wg.Add(2)
	go write(ch, ctx)
	go reader(ch, db, ctx)

	wg.Wait()
}
