package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/queensaver/packages/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var httpServerPort = flag.String("http_server_port", "8080", "HTTP server port")
var allowedOrigins = flag.String("allowed_origins", "https://dev.pado.mayrwoeger.com", "Allow-listed domains that will be receiving the Access-Control-Allow-Origin header.")
var senderName = flag.String("sender_name", "Carolina Reitmann", "Sender Full Name")
var senderAddress = flag.String("sender_address", "office@rundumsbauchgefuehl.at", "E-Mail Address of the sender")

type data struct {
	Name           string `json:name`
	EmailAddress   string `json:emailAddress`
	Address        string `json:address`
	Comment        string `json:comment`
	reCaptchaToken string `json:reCaptchaToken`
}

func dataHandler(w http.ResponseWriter, req *http.Request) {
	origins := strings.Split(*allowedOrigins, ",")
	for _, origin := range origins {
		if origin == req.Header.Get("Origin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			break
		}
	}

	var d data
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&d)
	if err != nil {
		logger.Error("Decoder Error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	from := mail.NewEmail(*senderName, *senderAddress)
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", "wogri@wogri.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		logger.Error("sendrgrid error", "error", err)
	} else {
		logger.Info("Sendgrid report",
			"status", response.StatusCode,
			"body", response.Body,
			"headers", response.Headers)
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	http.HandleFunc("/api/send", dataHandler)
	logger.Fatal("ListenAndServe Error", "error", http.ListenAndServe(":"+*httpServerPort, nil))
}
