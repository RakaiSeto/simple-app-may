package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

func index_handler(w http.ResponseWriter, r *http.Request) {
    // MAIN SECTION HTML CODE
	img, err := os.Open("qr.jpeg")
    if err != nil {
        panic(err) // perhaps handle this nicer
    }

	w.Header().Set("Content-Type", "image/jpeg") // <-- set the content-type header
    io.Copy(w, img)
}


func after_handler(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")

	fmt.Printf("user agent is: %s \n", ua)
	w.Write([]byte("user agent is " + ua + "\n"))

	result := "no"

	if strings.Contains(ua, "mobile") {
		result = "yes"
	}

	fmt.Printf("user agent is a mobile: %v \n", strings.Contains(ua, "Mobile"))
	w.Write([]byte("user agent is a mobile:" + result + "\n"))
}

func main() {
	err := qrcode.WriteFile("https://a83f3038534822.lhrtunnel.link/after_qr", qrcode.Highest, 256, "qr.jpeg")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/after_qr", after_handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}