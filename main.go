package main

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/codegangsta/cli"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"html"
	"image/png"
	"log"
	"net/http"
	"os"
)

func infoPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(
		w,
		`<html><head></head><body>
        <h1>Barcode Server</h1>
        Supported Codes:
        <ul>
            <li><a href="code128/12345">Code128</a></li>
        </ul>
        </body></html>`,
	)
}

func barcodeDisplayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := fmt.Sprintf("../../i/%s/%s.png", vars["type"], vars["data"])
	page := fmt.Sprintf(
		"<html><head></head><body><img src=\"%s\"/></body></html>",
		html.EscapeString(url),
	)
	fmt.Fprintf(w, page)
}

func barcodeEncoder(w http.ResponseWriter, r *http.Request) {
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
