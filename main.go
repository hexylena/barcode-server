package main

import (
    "fmt"
    "html"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/boombuler/barcode"
    "github.com/boombuler/barcode/code128"
    "image/png"
)

func BarcodeDisplayer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    url := fmt.Sprintf("../../i/%s/%s.png", vars["type"], vars["data"])
    page := fmt.Sprintf(
        "<html><head></head><body><img src=\"%s\"/></body></html>",
        html.EscapeString(url),
    )
    fmt.Fprintf(w, page)
}

func BarcodeEncoder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    if vars["type"] == "code128" {
        qrcode, err := code128.Encode(vars["data"])
        if err != nil {
            fmt.Println(err)
        } else {
            qrcode, err = barcode.Scale(qrcode, 220, 60)
            if err != nil {
                fmt.Println(err)
            } else {
                png.Encode(w, qrcode)
            }
        }
    } else {
        fmt.Fprintf(w, "Unsupported barcode type");
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/{type}/{data}", BarcodeDisplayer)
    r.HandleFunc("/i/{type}/{data}.png", BarcodeEncoder)
    //http.HandleFunc("/{type}", BarcodeEncoder)
    log.Fatal(http.ListenAndServe(":8080", r))
}
