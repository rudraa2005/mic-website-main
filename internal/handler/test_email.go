package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/rudraa2005/mic-website-main/backend/internal/email"
)

type TestEmailHandler struct {
	emailService email.Service
}

func NewTestEmailHandler(emailService email.Service) *TestEmailHandler {
	return &TestEmailHandler{
		emailService: emailService,
	}
}

func (h *TestEmailHandler) SendTestEmail(w http.ResponseWriter, r *http.Request) {
	err := h.emailService.Send(
		"rudranil.mitblr2024@learner.manipal.edu",
		"SMTP Test",
		"If you got this, SMTP works.",
	)
	log.Println(
		"SMTP CHECK",
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		len(os.Getenv("SMTP_PASSWORD")),
	)
	if err != nil {
		log.Println("SMTP TEST FAILED:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed"))
		return
	}

	log.Println("SMTP TEST SUCCESS")
	w.Write([]byte("sent"))
}
