package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/queensaver/packages/logger"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var httpServerPort = flag.String("http_server_port", "8080", "HTTP server port")
var allowedOrigins = flag.String("allowed_origins", "https://dev.pado.mayrwoeger.com", "Allow-listed domains that will be receiving the Access-Control-Allow-Origin header.")
var senderName = flag.String("sender_name", "Carolina Reitmann", "Sender Full Name")
var senderAddress = flag.String("sender_address", "office@rundumsbauchgefuehl.at", "E-Mail Address of the sender")
var ownerAddress = flag.String("owner_address", "wogri@wogri.com", "E-Mail Address of owner of the system")

const siteVerifyURL = "https://www.google.com/recaptcha/api/siteverify"
const plainMail = `Vielen Dank für Deine Bestellung bei Mutterliebe!

Wir kümmern uns möglichst rasch um Deine Bestellung und melden uns sobald es etwas Neues gibt!`

const htmlMail = `<strong>Vielen Dank für Deine Bestellung bei Mutterliebe!</strong>

Wir kümmern uns möglichst rasch um Deine Bestellung und melden uns sobald es etwas Neues gibt!`

const orderMail = `Neue Bestellung!

%s
%s
%s
%s
`

type data struct {
	Name           string `json:name`
	EmailAddress   string `json:emailAddress`
	Address        string `json:address`
	Comment        string `json:comment`
	ReCaptchaToken string `json:reCaptchaToken`
}

type SiteVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
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

	logger.Info("received request", "req", d)
	secret, ok := os.LookupEnv("RECAPTCHA_SECRET")
	if ok {
		logger.Info("ReCaptcha Secret Set!")
	} else {
		logger.Info("ReCaptcha Secret Not Set!", "secret", secret)
	}
	err = CheckRecaptcha(secret, d.ReCaptchaToken)
	if err != nil {
		logger.Error("recaptcha error", "error", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	from := mail.NewEmail(*senderName, *senderAddress)
	subject := "Vielen Dank für Deine Bestellung!"
	to := mail.NewEmail(d.Name, d.EmailAddress)
	message := mail.NewSingleEmail(from, subject, to, plainMail, htmlMail)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		logger.Error("sendgrid error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		logger.Info("Sendgrid report",
			"status", response.StatusCode,
			"body", response.Body,
			"headers", response.Headers)
	}

	from = mail.NewEmail(*senderName, *senderAddress)
	subject = "Mutterliebe: Neue Bestellung"
	to = mail.NewEmail("Carolina Reitmann", *ownerAddress)
	order := fmt.Sprintf(orderMail, d.Name, d.EmailAddress, d.Address, d.Comment)
	message = mail.NewSingleEmailPlainText(from, subject, to, order)
	client = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err = client.Send(message)
	if err != nil {
		logger.Error("sendgrid owner mail error", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		logger.Info("Sendgrid report",
			"status", response.StatusCode,
			"body", response.Body,
			"headers", response.Headers)
	}

	w.WriteHeader(http.StatusOK)
}

func CheckRecaptcha(secret, response string) error {
	req, err := http.NewRequest(http.MethodPost, siteVerifyURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("secret", secret)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var body SiteVerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return err
	}

	// Check recaptcha verification success.
	if !body.Success {
		return errors.New(fmt.Sprintf("unsuccessful recaptcha verify request: %s", body))
	}

	// Check response score.
	if body.Score < 0.5 {
		return errors.New("lower received score than expected")
	}

	// Check response action.
	if body.Action != "verify_bauchgefuehl" {
		return errors.New("mismatched recaptcha action")
	}

	return nil
}

func main() {
	flag.Parse()
	http.HandleFunc("/api/HttpTrigger", dataHandler)
	port := *httpServerPort
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		port = val
	}
	logger.Info("starting up")
	logger.Fatal("ListenAndServe Error", "error", http.ListenAndServe(":"+port, nil))
}
