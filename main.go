package main

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/qr"
	"github.com/codegangsta/cli"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"html"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	version     = "1.3"
	hostname, _ = os.Hostname()
	builddate   string
	gitrev      string
	prefix      string
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
			<li><a href="info/meta?q=L">Info</a></li>
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

func qrQuality(qualityLevel string) qr.ErrorCorrectionLevel {
	if qualityLevel == "M" {
		return qr.M
	} else if qualityLevel == "Q" {
		return qr.Q
	} else if qualityLevel == "H" {
		return qr.H
	}
	// default to lowest setting
	return qr.L
}

func redir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, prefix, http.StatusSeeOther)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
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
		encodingQuality := qrQuality(qualityLevel)

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
	} else if vars["type"] == "info" {
		qualityLevel := r.FormValue("q")
		encodingQuality := qrQuality(qualityLevel)

		qrcode, err := qr.Encode(
			fmt.Sprintf("barcode-server version:v%s host:%s builddate:%s gitrev:%s", version, hostname, builddate, gitrev),
			encodingQuality,
			qr.Auto,
		)
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

func serve(listenAddr, prefix string) {
	r := mux.NewRouter()
	r.HandleFunc("/", redir)
	r.HandleFunc(prefix, infoPage)
	r.HandleFunc(prefix+"/{type}/{data}", barcodeDisplayer)
	r.HandleFunc(prefix+"/i/{type}/{data}.png", barcodeEncoder)
	r.HandleFunc(prefix+"/healthcheck", healthCheck)
	fmt.Printf("Listening on %s%s\n", listenAddr, prefix)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(listenAddr, loggedRouter))
}

func main() {
	app := cli.NewApp()
	app.Name = "barcode-server"
	app.Usage = "dynamically generates + serves barcodes"
	app.Version = fmt.Sprintf("%s (%s)", version, builddate)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Value: "0.0.0.0:8080",
			Usage: "Address to listen on",
		},
		cli.StringFlag{
			Name:  "prefix, p",
			Value: "/barcodes",
			Usage: "URL Path prefix",
		},
	}
	app.Action = func(c *cli.Context) {
		prefix = c.String("prefix")
		serve(
			c.String("listen"),
			c.String("prefix"),
		)
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
