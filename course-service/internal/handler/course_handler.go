package handler

import (
	"main/course-service/internal/dto"
	"main/course-service/internal/middleware"
	"main/course-service/internal/model"
	"main/course-service/internal/service"
	"main/course-service/internal/validation"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CourseHandler struct {
	service service.CourseService
}

func NewCourseHandler(service service.CourseService) *CourseHandler {
	return &CourseHandler{
		service: service,
	}
}
func (h *CourseHandler) Create(c fiber.Ctx) error {
	var req dto.CourseRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	course, err := h.service.Create(c.Context(), actorFromCtx(c), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(service.ToCourseDTO(course))

}
func (h *CourseHandler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	course, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(service.ToCourseDTO(course))
}
func (h *CourseHandler) GetManagedCourse(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	course, err := h.service.GetManagedCourse(c.Context(), actorFromCtx(c), id)
	if err != nil {
		return err
	}
	return c.JSON(service.ToCourseDTO(course))
}
func (h *CourseHandler) ListManagedLessons(c fiber.Ctx) error {
	courseID, err := uuid.Parse(c.Params("courseID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	lessons, err := h.service.ListManagedLessons(c.Context(), actorFromCtx(c), courseID)
	if err != nil {
		return err
	}
	return c.JSON(ToLessonDTOList(lessons))
}
func (h *CourseHandler) Update(c fiber.Ctx) error {
	var req dto.UpdateCourseRequest
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	course, err := h.service.Update(c.Context(), actorFromCtx(c), id, req)
	if err != nil {
		return err
	}
	return c.JSON(service.ToCourseDTO(course))
}
func (h *CourseHandler) Publish(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	course, err := h.service.Publish(c.Context(), actorFromCtx(c), id)
	if err != nil {
		return err
	}
	return c.JSON(service.ToCourseDTO(course))
}
func (h *CourseHandler) Archive(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	course, err := h.service.Archive(c.Context(), actorFromCtx(c), id)
	if err != nil {
		return err
	}
	return c.JSON(service.ToCourseDTO(course))
}
func (h *CourseHandler) CreateLesson(c fiber.Ctx) error {
	var req dto.CreateLessonRequest
	courseID, err := uuid.Parse(c.Params("courseID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	lesson, err := h.service.CreateLesson(c.Context(), actorFromCtx(c), courseID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(service.ToLessonDTO(lesson))
}
func (h *CourseHandler) UpdateLesson(c fiber.Ctx) error {
	var req dto.UpdateLessonRequest
	lessonID, err := uuid.Parse(c.Params("lessonID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	lesson, err := h.service.UpdateLesson(c.Context(), actorFromCtx(c), lessonID, req)
	if err != nil {
		return err
	}
	return c.JSON(service.ToLessonDTO(lesson))
}
func (h *CourseHandler) ListLessons(c fiber.Ctx) error {
	courseID, err := uuid.Parse(c.Params("courseID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	lessons, err := h.service.ListLessons(c.Context(), courseID)
	if err != nil {
		return err
	}
	return c.JSON(ToLessonDTOList(lessons))
}
func (h *CourseHandler) DeleteLesson(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("lessonID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := h.service.DeleteLesson(c.Context(), actorFromCtx(c), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *CourseHandler) ListPublic(c fiber.Ctx) error {
	search := c.Query("search")
	limit := fiber.Query[int](c, "limit", 10)
	offset := fiber.Query[int](c, "offset", 0)
	courses, err := h.service.ListPublic(c.Context(), search, limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(service.ToCourseDTOList(courses))
}
func (h *CourseHandler) ListCourses(c fiber.Ctx) error {
	limit := fiber.Query[int](c, "limit", 10)
	offset := fiber.Query[int](c, "offset", 0)
	courses, err := h.service.ListMyCourses(c.Context(), actorFromCtx(c), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(service.ToCourseDTOList(courses))
}
func actorFromCtx(c fiber.Ctx) service.Actor {
	return service.Actor{
		UserID: fiber.Locals[uuid.UUID](c, middleware.UserIDKey),
		Role:   fiber.Locals[string](c, middleware.RoleKey),
	}
}
func ToLessonDTOList(lessons []model.Lesson) []dto.LessonResponse {
	res := make([]dto.LessonResponse, 0, len(lessons))
	for i := range lessons {
		res = append(res, service.ToLessonDTO(&lessons[i]))
	}
	return res
}
