package service

import (
	"context"
	"log"

	"github.com/rudraa2005/mic-website-main/backend/internal/email"
	"github.com/rudraa2005/mic-website-main/backend/internal/model"
)

type NotificationRepo interface {
	NotifyStatusChange(ctx context.Context, notification *model.Notification) error
	GetNotificationsByUser(ctx context.Context, userID string) ([]model.Notification, error)
	MarkNotificationAsRead(ctx context.Context, id string, userID string) error
	GetUnreadCountByUser(ctx context.Context, userID string) (int, error)
}

type NotificationService struct {
	notificationRepo NotificationRepo
	emailService     email.Service
}

func NewNotificationService(notificationRepo NotificationRepo, emailService email.Service) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		emailService:     emailService,
	}
}

func (ns *NotificationService) createNotification(ctx context.Context, notification *model.Notification) error {

	return ns.notificationRepo.NotifyStatusChange(ctx, notification)
}

func (ns *NotificationService) GetNotificationsByUser(ctx context.Context, userID string) ([]model.Notification, error) {
	return ns.notificationRepo.GetNotificationsByUser(ctx, userID)
}

func (ns *NotificationService) MarkNotificationAsRead(ctx context.Context, id string, userID string) error {
	return ns.notificationRepo.MarkNotificationAsRead(ctx, id, userID)
}

func (ns *NotificationService) GetUnreadCountByUser(ctx context.Context, userID string) (int, error) {
	return ns.notificationRepo.GetUnreadCountByUser(ctx, userID)
}

func (ns *NotificationService) NotifyStatusChange(ctx context.Context, userID string, email string, submissionID string, oldStatus string, newStatus string) error {

	log.Println(
		"[NOTIFICATION SERVICE HIT]",
		"userID=", userID,
		"email=", email,
		"submissionID=", submissionID,
		"old=", oldStatus,
		"new=", newStatus,
	)

	var message string
	var emailSubject string
	var emailBody string

	switch {
	case oldStatus == "draft" && newStatus == "submitted":
		message = "Your submission has been successfully submitted."
		emailSubject = "Submission Submitted"
		emailBody = "Dear User,\n\nYour submission with ID " + submissionID + " has been successfully submitted.\n\nBest regards,\nTeam"
	case oldStatus == "submitted" && newStatus == "approved":
		message = "Your submission is now approved."
		emailSubject = "Submission Approved"
		emailBody = "Dear User,\n\nYour submission with ID " + submissionID + " is now approved.\n\nBest regards,\nTeam MIC"

	default:
		return nil
	}

	err := ns.createNotification(ctx, &model.Notification{
		UserID: userID,
		Type:   "status_change",
		Title:  emailSubject,
		Body:   message,
	})
	if err != nil {
		return err
	}

	go func() {
		err := ns.emailService.Send(email, emailSubject, emailBody)
		if err != nil {
			log.Println("[EMAIL FAILED]", err)
		} else {
			log.Println("[EMAIL SENT SUCCESS]", email)
		}
	}()

	return nil

}
func (ns *NotificationService) SendSubmissionStatusUpdate(ctx context.Context, email, title, status string) error {
	subject := "Submission Update: " + title
	body := "Dear User,\n\nYour idea '" + title + "' has been " + status + " by the faculty review committee.\n\nBest regards,\nMAHE Innovation Centre"

	go func() {
		err := ns.emailService.Send(email, subject, body)
		if err != nil {
			log.Println("[EMAIL FAILED]", err)
		}
	}()

	return nil
}
