package router

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rudraa2005/mic-website-main/backend/internal/handler"
	appmw "github.com/rudraa2005/mic-website-main/backend/internal/middleware"
)

func NewRouter(sh *handler.StartupHandler, ah *handler.AuthHandler, ph *handler.ProfileHandler, seh *handler.SettingsHandler, subh *handler.SubmissionsHandler, fh *handler.FeedbackHandler, qh *handler.QueryHandler, th *handler.TestEmailHandler, aih *handler.AIHandler, ch *handler.ContentHandler, frh *handler.FacultyReviewHandler, feh *handler.EventInvitationHandler, fph *handler.FacultyProgressHandler, afh *handler.AdminFacultyHandler, ash *handler.AdminSubmissionHandler, workh *handler.WorkHandler, fih *handler.FacultyIncubationHandler, awh *handler.AdminWorkHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/test-email", th.SendTestEmail)
	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "frontend/static"))
	FileServer(r, "/static", filesDir)

	// Serve chatbot assets from frontend/templates/assets/
	assetsDir := http.Dir(filepath.Join(workDir, "frontend/templates/assets"))
	FileServer(r, "/assets", assetsDir)

	r.Route("/api", func(r chi.Router) {
		// Chatbot proxy — no auth required (public-facing)
		r.Post("/chat", func(w http.ResponseWriter, r *http.Request) {
			proxyURL := "http://localhost:9000/chat"
			proxyReq, err := http.NewRequest("POST", proxyURL, r.Body)
			if err != nil {
				log.Printf("[CHAT PROXY] Failed to create request: %v", err)
				http.Error(w, `{"error":"proxy error"}`, http.StatusBadGateway)
				return
			}
			proxyReq.Header.Set("Content-Type", "application/json")
			if sid := r.Header.Get("X-Session-ID"); sid != "" {
				proxyReq.Header.Set("X-Session-ID", sid)
			}

			client := &http.Client{}
			resp, err := client.Do(proxyReq)
			if err != nil {
				log.Printf("[CHAT PROXY] Python backend error: %v", err)
				http.Error(w, `{"error":"AI service unavailable"}`, http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		})

		r.Post("/login", ah.Login)
		r.Post("/signup", ah.Signup)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		r.Get("/submissions/incubation", workh.GetIncubationPipeline)

		// Content management routes (accessible by ADMIN role)
		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRoles("ADMIN", "FACULTY"))

			r.Get("/contents", ch.GetAllContent)
			r.Post("/create-content", ch.CreateContent)
			r.Put("/contents/{id}", ch.UpdateContent)
			r.Delete("/contents/{id}", ch.DeleteContent)
		})

		// Admin submission review routes
		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("ADMIN"))

			r.Get("/admin/submissions", ash.GetPendingSubmissions)
			r.Get("/admin/submissions/all", ash.GetAllSubmissions)
			r.Post("/admin/submissions/{id}/decision", ash.DecideSubmission)

			// Faculty assignment routes
			r.Post("/admin/submissions/{id}/assign-faculty", ash.AssignFaculty)
			r.Delete("/admin/submissions/{id}/assign-faculty/{faculty_id}", ash.RemoveFaculty)
			r.Get("/admin/submissions/{id}/faculty", ash.GetAssignedFaculty)

			// Tags management
			r.Put("/admin/submissions/{id}/tags", ash.UpdateTags)
		})

		// Admin faculty management routes
		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRoles("ADMIN", "FACULTY"))

			r.Get("/admin/faculty", afh.GetAllFaculty)
			r.Post("/admin/faculty", afh.CreateFaculty)
			r.Put("/admin/faculty/{id}", afh.UpdateFaculty)
			r.Delete("/admin/faculty/{id}", afh.DeleteFaculty)
		})

		// Admin work/pipeline management routes
		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("ADMIN"))

			r.Get("/admin/work", awh.GetAllWork)
			r.Put("/admin/work/{id}", awh.UpdateWork)
			r.Delete("/admin/work/{id}", awh.DeleteWork)

			r.Get("/admin/companies", awh.GetCompanies)
			r.Post("/admin/companies", awh.AddCompany)
			r.Delete("/admin/companies/{id}", awh.DeleteCompany)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("FACULTY"))

			r.Get("/faculty/reviews", frh.GetSubmitted)
			r.Get("/faculty/reviews/{id}", frh.GetByID)
			r.Post("/faculty/reviews/{id}/decision", frh.Decide)

			r.Get("/faculty/events/invitations", feh.GetMyInvitations)
			r.Post("/faculty/events/invitations/{invitation_id}/rsvp", feh.UpdateRSVP)
			r.Get("/faculty/progress", fph.GetMyProgress)
			r.Get("/faculty/progress/{submission_id}", fph.GetProgressBySubmission)
			r.Post("/faculty/feedback", fh.Create)
			r.Get("/submissions/{submission_id}/file", subh.DownloadSubmissionFile)

			// Incubation Portfolio
			r.Get("/faculty/incubation", fih.GetPortfolio)
			r.Post("/faculty/incubation/{submission_id}", fih.UpdateProgress)
			r.Get("/faculty/companies", fih.GetCompanies)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("STUDENT"))

			r.Post("/startups", sh.Create)
			r.Get("/startups/mine", sh.GetMine)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)

			r.Get("/profile/me", ph.Me)
			r.Post("/profile/photo", ph.UploadPhoto)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("STUDENT"))

			r.Get("/settings", seh.GetSettings)
			r.Post("/settings/update", seh.UpdateSettings)
			r.Post("/settings/update-password", ah.ChangePassword)
			r.Post("/settings/profile-photo", ph.UploadPhoto)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("STUDENT"))

			r.Post("/submissions/create", subh.CreateSubmission)
			r.Post("/submissions/submit/{submission_id}", subh.SubmitSubmission)
			r.Put("/submissions/{submission_id}", subh.UpdateSubmission)
			r.Get("/submissions/mine", subh.GetByUserID)
			r.Get("/submissions/{submission_id}", subh.GetBySubmissionID)
			r.Delete("/submissions/{submission_id}", subh.DeleteSubmission)
			r.Post("/submissions/{submission_id}/attach-file", subh.UploadSubmissionFile)
			r.Get("/submissions/{submission_id}/file", subh.DownloadSubmissionFile)
		})
		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("STUDENT"))

			r.Get("/feedbacks", fh.GetMyFeedbacks)

			r.Post("/queries", qh.CreateQuery)
			r.Get("/queries/mine", qh.GetMyQueries)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmw.AuthMiddleware)
			r.Use(appmw.RequireRole("STUDENT"))

			r.Post("/ai/analyze", aih.AnalyzeDraft)
			r.Get("/ai/insights/{submission_id}", aih.GetInsights)
		})

		r.Route("/content", func(r chi.Router) {
			r.Get("/about/cards", ch.GetAboutCards)
			r.Get("/about/features", ch.GetAboutFeatures)
			r.Get("/about/team", ch.GetTeamMembers)
			r.Get("/about/testimonials", ch.GetTestimonials)
			r.Get("/about/stats", ch.GetStats)

			r.Get("/resources", ch.GetResources)
			r.Get("/resources/top", ch.GetTopResources)

			r.Get("/events/upcoming", ch.GetUpcomingEvents)
			r.Get("/events/all", ch.GetAllEvents)
		})
	})

	templatesBase := filepath.Join(workDir, "frontend/templates")

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		path := r.URL.Path

		if path == "/" {
			path = "/index.html"
		} else if !strings.Contains(filepath.Base(path), ".") {
			path = path + ".html"
		}

		fullPath := filepath.Join(templatesBase, path)
		http.ServeFile(w, r, fullPath)
	})

	return r
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
