// Command pdf is a chromedp example demonstrating how to capture a pdf of a
// page.
package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	maxConcurrency = 0 // keep concurrent renders at bay
	limiter        = make(chan nothing, 1+maxConcurrency)
	httpPort       = os.Getenv("PORT")

	chromeCtx context.Context
)

type nothing struct{}

func main() {
	if httpPort == "" {
		httpPort = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	mux.HandleFunc("/pdf/", HandlePrintPDFRequest)      // used by whatever calls the service ( but we would use kafka likely )
	mux.HandleFunc("/pdf/{document_id}", HandleConvert) //internal

	host := strings.Join([]string{"0.0.0.0", httpPort}, ":")
	server := &http.Server{}
	server.Addr = host
	server.Handler = mux
	server.ReadHeaderTimeout = 10 * time.Second

	// var cancel func()
	chromeCtx, _ = chromedp.NewContext(context.TODO()) // listen for signals
	// chromeCtx, cancel = chromedp.NewContext(context.TODO()) // listen for signals
	// defer cancel()

	log.Printf("listening on %s", host)
	log.Println(server.ListenAndServe())
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	// default 200
}

func HandlePrintPDFRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("request")
	requestCtx := r.Context()
	select {
	case limiter <- nothing{}:
		defer func() { <-limiter }()
	case <-requestCtx.Done():
		http.Error(w, requestCtx.Err().Error(), http.StatusInternalServerError)
		return
	}

	// TODO: experiment with new browser vs. new tab performance

	w.Header().Add("Content-Type", "application/pdf")
	requestPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("payload read")
	log.Println("chromeCtx err:", chromeCtx.Err())

	// _, err = chromedp.RunResponse(chromeCtx, printToPDF(requestPayload, w))
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	res1 := []byte{}
	res2 := []byte{}
	res3 := []byte{}
	// fmt.Println(string(requestPayload))
	if err := chromedp.Run(chromeCtx,
		setContent(requestPayload),
		chromedp.FullScreenshot(&res1, 100),
		chromedp.Sleep(1*time.Second),
		chromedp.FullScreenshot(&res2, 100),
		printToPDF(w),
		chromedp.FullScreenshot(&res3, 100),
	); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}

	ioutil.WriteFile("screenshot1.png", res1, 0644)
	ioutil.WriteFile("screenshot2.png", res2, 0644)
	ioutil.WriteFile("screenshot3.png", res3, 0644)

}

func setContent(payload []byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, string(payload)).Do(ctx)
		}),
	}
}

// print a specific pdf page.
func printToPDF(res io.Writer) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithLandscape(false).
				WithPaperHeight(297.0 / 25.4).
				WithPaperWidth(210 / 25.4).
				WithDisplayHeaderFooter(false).
				WithMarginBottom(0.0).
				WithMarginLeft(0.0).
				WithMarginTop(0.0).
				WithMarginRight(0.0).
				WithPreferCSSPageSize(true).
				Do(ctx)
			if err != nil {
				return err
			}

			n, err := res.Write(buf)
			log.Printf("written %v bytes", n)
			return err
		}),
	}
}
