package main

import (
	"fmt"
	"html"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
	"github.com/codegangsta/cli"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func infoPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(
		w,
		`<html><head></head><body>
        <h1>Barcode Server</h1>
        Supported Codes:
        <ul>
            <li><a href="code128/12345">Code128</a></li>
            <li><a href="qr/12345?q=L">QR</a></li>
        </ul>
        </body></html>`,
	)
}

func barcodeDisplayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := fmt.Sprintf("../i/%s/%s.png", vars["type"], vars["data"])
	page := fmt.Sprintf(
		"<html><head></head><body><img src=\"%s\"/></body></html>",
		html.EscapeString(url),
	)
	fmt.Fprintf(w, page)
}

func barcodeEncoder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	width, err := strconv.Atoi(r.FormValue("w"))
	if err != nil {
		width = 200
	}

	height, err := strconv.Atoi(r.FormValue("h"))
	if err != nil {
		height = 200
	}

	if vars["type"] == "code128" {
		qrcode, err := code128.Encode(vars["data"])
		if err != nil {
			fmt.Println(err)
		} else {
			qrcode, err = barcode.Scale(qrcode, width, height)
			if err != nil {
				fmt.Println(err)
			} else {
				png.Encode(w, qrcode)
			}
		}
	} else if vars["type"] == "qr" {
		qualityLevel := r.FormValue("q")
		var encodingQuality qr.ErrorCorrectionLevel

		if qualityLevel == "M" {
			encodingQuality = qr.M
		} else if qualityLevel == "Q" {
			encodingQuality = qr.Q
		} else if qualityLevel == "H" {
			encodingQuality = qr.H
		} else {
			// default to lowest setting
			encodingQuality = qr.L
		}

		qrcode, err := qr.Encode(vars["data"], encodingQuality, qr.Auto)
		if err != nil {
			fmt.Println(err)
		} else {
			qrcode, err = barcode.Scale(qrcode, width, height)
			if err != nil {
				fmt.Println(err)
			} else {
				png.Encode(w, qrcode)
			}
		}
	} else {
		fmt.Fprintf(w, "Unsupported barcode type")
	}
}

func serve(listenAddr string) {
	r := mux.NewRouter()
	r.HandleFunc("/", infoPage)
	r.HandleFunc("/{type}/{data}", barcodeDisplayer)
	r.HandleFunc("/i/{type}/{data}.png", barcodeEncoder)
	fmt.Printf("Listening on %s\n", listenAddr)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(listenAddr, loggedRouter))
}

func main() {
	app := cli.NewApp()
	app.Name = "barcode-server"
	app.Usage = "dynamically generates + serves barcodes"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: "0.0.0.0:8080",
			Usage: "Address to listen on",
		},
	}
	app.Action = func(c *cli.Context) {
		serve(c.String("listen"))
	}

	app.Run(os.Args)
}
