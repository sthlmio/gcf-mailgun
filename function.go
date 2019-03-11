package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (r Response) Json(w http.ResponseWriter) {
	res, _ := json.Marshal(r)

	if _, err := fmt.Fprintf(w, string(res)); err != nil {
		panic(err)
	}

	return
}

func Mail(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Set headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusForbidden)

		res := Response{
			Error: "Wrong Method, only accepting 'POST'",
		}

		res.Json(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)

		res := Response{
			Error: "Wrong Content-Type, only accepting 'application/json'",
		}

		res.Json(w)
		return
	}

	var d struct {
		Email   string `json:"email"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		panic(err)
		return
	}

	if d.Email == "" || d.Message == "" {
		w.WriteHeader(http.StatusBadRequest)

		res := Response{
			Error: "Empty data, 'email' or 'message'",
		}

		res.Json(w)
		return
	}

	mg := mailgun.NewMailgun(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_API_KEY"),
	)

	mg.SetAPIBase(os.Getenv("MAILGUN_API_BASE"))

	sender := d.Email
	subject := os.Getenv("MAIL_SUBJECT")
	body := d.Message
	recipient := os.Getenv("MAIL_RECIPIENT")

	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	_, _, err := mg.Send(ctx, message)

	if err != nil {
		log.Printf("Could not send email: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		res := Response{
			Error: err.Error(),
		}

		res.Json(w)
		return
	}

	w.WriteHeader(http.StatusOK)

	res := Response{
		Message: "All good cowboy!",
	}

	res.Json(w)
	return
}
