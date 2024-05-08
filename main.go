package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	/*"io/ioutil"
	"net/http"
	"strconv"*/)

var wg sync.WaitGroup

func write(ch chan<- int) {
	defer wg.Done()

	/*resp, err := http.Get("http://138.201.177.104:3040/ping") // <-- This URL appears to be blocked
	time.Sleep(time.Second * 1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	numStr := string(body)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}*/

	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}

func read(ch <-chan int) error {

	// part two
	file, err := os.Create("Number.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for {
		num, ok := <-ch
		if !ok {
			break
		}
		_, err := fmt.Fprintln(writer, num)
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	wg.Done()
	return nil
}

func main() {
	ch := make(chan int)

	wg.Add(2)
	go write(ch)
	go read(ch)
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
