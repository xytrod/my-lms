package handler

import (
	"main/user-service/internal/dto"
	"main/user-service/internal/model"
	"main/user-service/internal/service"
	"main/user-service/internal/validation"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}
func ToUserDTO(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
func ToUserDTOList(users []model.User) []dto.UserResponse {
	resp := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		resp = append(resp, ToUserDTO(&users[i]))
	}
	return resp
}
func (h *UserHandler) Create(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	user, err := h.service.CreateUser(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(ToUserDTO(user))

}
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	idString := c.Params("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest,
			"invalid id")
	}
	user, err := h.service.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(ToUserDTO(user))
}
func (h *UserHandler) Update(c fiber.Ctx) error {
	idString := c.Params("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var req dto.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	req.ID = id
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	user, err := h.service.UpdateUser(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(ToUserDTO(user))
}
func (h *UserHandler) Delete(c fiber.Ctx) error {
	idString := c.Params("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	err = h.service.DeleteUser(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *UserHandler) List(c fiber.Ctx) error {
	limit := fiber.Query[int](c, "limit", 10)
	offset := fiber.Query[int](c, "offset", 0)
	users, err := h.service.List(c.Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(ToUserDTOList(users))
}
