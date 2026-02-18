package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rudraa2005/mic-website-main/backend/internal/db"
	"github.com/rudraa2005/mic-website-main/backend/internal/email"
	"github.com/rudraa2005/mic-website-main/backend/internal/handler"
	h "github.com/rudraa2005/mic-website-main/backend/internal/handler"
	"github.com/rudraa2005/mic-website-main/backend/internal/repository"
	r "github.com/rudraa2005/mic-website-main/backend/internal/router"
	"github.com/rudraa2005/mic-website-main/backend/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	log.Println("Starting MIC Website Backend Server...")
	log.Println("Connecting to database...")

	pool, err := db.NewPool()
	if err != nil {
		log.Fatal("DB connection failed: ", err)
	}
	defer pool.Close()

	log.Println("Database connection established successfully")

	userRepo := repository.NewAuthRepository(pool)
	profileRepo := repository.NewProfileRepo(pool)
	settingsRepo := repository.NewSettingsRepo(pool)
	startupRepo := repository.NewStartupRepository(pool)
	submissionRepo := repository.NewSubmissionsRepo(pool)
	feedbackRepo := repository.NewFeedbackRepo(pool)
	companyRepo := repository.NewCompanyRepo(pool)
	queryRepo := repository.NewQueryRepo(pool)
	queryService := service.NewQueryService(queryRepo)
	feedbackService := service.NewFeedbackService(feedbackRepo)

	facultyRepo := repository.NewFacultyRepository(pool)
	emailService := email.NewSMTPService(
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
	)
	notificationRepo := repository.NewNotificationRepository(pool)
	notificationService := service.NewNotificationService(notificationRepo, emailService)

	queryHandler := handler.NewQueryHandler(queryService)
	facultyReviewRepo := repository.NewFacultySubmissionRepo(pool)
	facultyReviewService := service.NewFacultyReviewService(facultyReviewRepo, notificationService)
	facultyReviewHandler := handler.NewFacultyReviewHandler(facultyReviewService)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)
	aiRepo := repository.NewAIRepo(pool)
	aiService := service.NewAIService(
		"http://localhost:9000",
		aiRepo,
		&http.Client{},
	)
	aiHandler := handler.NewAIHandler(aiService)
	authService := service.NewAuthService(userRepo, facultyRepo, profileRepo, settingsRepo)
	authHandler := handler.NewAuthHandler(authService)

	contentRepo := repository.NewContentRepository(pool)
	contentService := service.NewContentService(contentRepo)
	contentHandler := handler.NewContentHandler(contentService)
	facultyEventRepo := repository.NewEventInvitationRepository(pool)
	facultyProgressRepo := repository.NewFacultyProgressRepository(pool)
	facultyProgressService := service.NewFacultyProgressService(facultyProgressRepo)
	facultyProgressHandler := handler.NewFacultyProgressHandler(facultyProgressService)

	startupService := service.NewStartupService(startupRepo)
	startupHandler := h.NewStartupHandler(startupService)
	submissionService := service.NewSubmissionsService(notificationService, submissionRepo, profileRepo, aiService)
	submissionHandler := handler.NewSubmissionsHandler(submissionService)
	testEmailHandler := handler.NewTestEmailHandler(emailService)
	settingService := service.NewSettingService(settingsRepo)
	profileService := service.NewProfileService(profileRepo)
	profileHandler := handler.NewProfileHandler(profileService)
	settingsHandler := handler.NewSettingsHandler(settingService, profileService, authService)
	facultyEventService := service.NewEventInvitationService(facultyEventRepo)
	facultyEventHandler := handler.NewEventInvitationHandler(facultyEventService)

	adminFacultyRepo := repository.NewAdminFacultyRepository(pool)
	adminFacultyService := service.NewAdminFacultyService(adminFacultyRepo)
	adminFacultyHandler := handler.NewAdminFacultyHandler(adminFacultyService)

	adminSubmissionRepo := repository.NewAdminSubmissionRepo(pool)
	adminSubmissionHandler := handler.NewAdminSubmissionHandler(adminSubmissionRepo)

	adminWorkRepo := repository.NewAdminWorkRepo(pool)
	adminWorkHandler := handler.NewAdminWorkHandler(adminWorkRepo)

	facultyIncubationHandler := handler.NewFacultyIncubationHandler(facultyProgressService, companyRepo)
	workHandler := handler.NewWorkHandler(submissionRepo)

	router := r.NewRouter(startupHandler, authHandler, profileHandler, settingsHandler, submissionHandler, feedbackHandler, queryHandler, testEmailHandler, aiHandler, contentHandler, facultyReviewHandler, facultyEventHandler, facultyProgressHandler, adminFacultyHandler, adminSubmissionHandler, workHandler, facultyIncubationHandler, adminWorkHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", router)

}
