package routes

import (
	"main/course-service/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterCourseRoutes(app *fiber.App, h *handler.CourseHandler, jwtMiddleware fiber.Handler) {
	courses := app.Group("/courses")
	courses.Get("", h.ListPublic)
	courses.Get("/:courseID/lessons", h.ListLessons)
	courses.Get("/my", jwtMiddleware, h.ListCourses)
	courses.Get("/my/:courseID/lessons", jwtMiddleware, h.ListManagedLessons)
	courses.Get("/my/:id", jwtMiddleware, h.GetManagedCourse)
	courses.Get("/:id", h.GetByID)
	courses.Post("/:courseID/lessons", jwtMiddleware, h.CreateLesson)
	courses.Patch("/lessons/:lessonID", jwtMiddleware, h.UpdateLesson)
	courses.Delete("/lessons/:lessonID", jwtMiddleware, h.DeleteLesson)
	courses.Post("", jwtMiddleware, h.Create)
	courses.Patch("/:id", jwtMiddleware, h.Update)
	courses.Post("/:id/publish", jwtMiddleware, h.Publish)
	courses.Post("/:id/archive", jwtMiddleware, h.Archive)

}
