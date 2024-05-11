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

func write(ch chan<- int, db *leveldb.DB) {

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

		key := []byte(strconv.FormatInt(time.Now().Unix(), 10))
		value := buffer.Bytes()

		err = db.Put(key, value, nil)
		if err != nil {
			fmt.Println("not send key & value in the databace")
			return
		}

		data, err := db.Get([]byte(key), nil)
		if err != nil {
			log.Fatal("noooooot")
		}
		log.Printf("Value: %s\n", data)

		number := buffer.String()
		nint, err := strconv.Atoi(number)
		if err != nil {
			fmt.Println("not chang 'byte' --> 'int' ")
		}
		ch <- nint
		time.Sleep(time.Second * 1)
		defer respons.Body.Close()
	}
	defer close(ch)
	defer wg.Done()
}

func reader(ch <-chan int) {

	// part one & two
	file, err := os.OpenFile("Number.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err : not write in file")
		return
	}

	for num := range ch {
		fmt.Fprintln(file, num)
		fmt.Println(num)
	}

	defer file.Close()
	defer wg.Done()
}

func main() {
	ch := make(chan int)

	db, err := leveldb.OpenFile("start-gorotins", nil)
	if err != nil {
		fmt.Println("err in open library", err)
		return
	}
	defer db.Close()

	wg.Add(2)
	go write(ch, db)
	go reader(ch)
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
