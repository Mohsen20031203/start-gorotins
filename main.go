package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

var wg sync.WaitGroup

func write(ch chan<- []byte) {

	for {
		respons, err := http.Get("http://138.201.177.104:3040/ping")
		if err != nil {
			fmt.Println("err : Get Number")
			return
		}

		var buffer bytes.Buffer

		_, err = io.Copy(&buffer, respons.Body)
		if err != nil {
			fmt.Println("not send & read body respons", err)
			return
		}

		ch <- buffer.Bytes()
		time.Sleep(time.Second * 1)
		defer respons.Body.Close()
	}
	defer close(ch)
	defer wg.Done()
}

func reader(ch <-chan []byte, db *leveldb.DB) {

	// part two
	file, err := os.OpenFile("Number.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err : not write in file")
		return
	}

	for num := range ch {
		// part one
		number := string(num)
		nint, err := strconv.Atoi(number)
		if err != nil {
			fmt.Println("not convert bytes to int")
			return
		}
		fmt.Println(nint)

		// part two
		fmt.Fprintln(file, nint)

		// part three
		value := num
		key := []byte(strconv.FormatInt(time.Now().Unix(), 10))

		err = db.Put(key, value, nil)
		if err != nil {
			fmt.Println("not send key & value in the databace")
			return
		}

		data, err := db.Get(key, nil)
		if err != nil {
			log.Fatal("dont get value databace")
		}
		log.Printf("Value: %s\n", data)
	}

	defer file.Close()
	defer wg.Done()
}

func main() {
	ch := make(chan []byte)

	db, err := leveldb.OpenFile("start-gorotins", nil)
	if err != nil {
		fmt.Println("err in open library", err)
		return
	}
	defer db.Close()

	wg.Add(2)
	go write(ch)
	go reader(ch, db)
	wg.Wait()
}

/*
برنامه‌ای بنویسید که دو go routine موازی با هم اجرا شوند.
یکی از این go routine ها باید یک http request به این ادرس بزند: http://138.201.177.104:3040/ping و عددی که دریافت می‌شود را به صورت مکرر به go routine دوم منتقل کند. این کار باید هر یک ثانیه تکرار شود.
در go routine دوم، سه فاز تعریف شده است که پس از پایان هر کدام به فاز بعدی می‌رویم. این go routine باید همیشه منتظر دریافت داده از go routine اول باشد و:
۱. در فاز اول، عدد دریافتی را در کنسول و صفحه نمایش چاپ کنید.
۲. در فاز دوم، عدد را در یک فایل ذخیره کنید.
۳. در فاز سوم، یک دیتابیس Leveldb بسازید و عدد دریافتی از سرور را در آن ذخیره کنید، هر کلید زمان درخواست و مقدار خروجی دریافتی از سرور است.
*/
