package main

import (
	"fmt"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func write(ch chan<- int) {

	respons, err := http.Get("http://138.201.177.104:3040/ping")
	if err != nil {
		fmt.Println("err : Get Number")
		return
	}
	defer respons.Body.Close()
	fmt.Println(respons)

	defer wg.Done()
	close(ch)
}

func main() {
	ch := make(chan int)

	wg.Add(1)
	go write(ch)
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
