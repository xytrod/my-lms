package routes

import (
	"main/enrollment_progress-service/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterEnrollmentRoutes(app *fiber.App, h *handler.EnrollmentHandler, jwtMiddleware fiber.Handler) {
	enrollments := app.Group("/enrollments", jwtMiddleware)
	enrollments.Get("/my", h.ListMy)
	enrollments.Post("/:courseID", h.Enroll)
	enrollments.Post("/:courseID/lessons/:lessonID/complete", h.CompletedLesson)
	enrollments.Get("/:courseID/progress", h.GetProgress)

}
