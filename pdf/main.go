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
)

type nothing struct{}

func main() {
	if httpPort == "" {
		httpPort = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	mux.HandleFunc("/pdf", HandlePrintPDFRequest)

	host := strings.Join([]string{"0.0.0.0", httpPort}, ":")
	server := &http.Server{}
	server.Addr = host
	server.Handler = mux
	server.ReadHeaderTimeout = 10 * time.Second

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
	chromeCtx, cancel := chromedp.NewContext(requestCtx)
	defer cancel()

	w.Header().Add("Content-Type", "application/pdf")
	err := chromedp.Run(chromeCtx, printToPDF(r.Body, w))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// print a specific pdf page.
func printToPDF(contents io.Reader, res io.Writer) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			html, err := ioutil.ReadAll(contents)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, string(html)).Do(ctx); err != nil {
				return err
			}

			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithLandscape(false).
				WithPaperHeight(297.0 / 2.54).
				WithPaperWidth(210 / 2.54).
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

			_, err = res.Write(buf)
			return err
		}),
	}
}
