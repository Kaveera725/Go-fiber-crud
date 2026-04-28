package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"go-fiber-crud/internal/database"
)

type ProductHandler struct {
	DB *database.DB
}

func NewProductHandler(db *database.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (h *ProductHandler) List(c *fiber.Ctx) error {
	products, err := h.DB.ListProducts(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Render("index", fiber.Map{"Products": products}, "layouts/base")
}

func (h *ProductHandler) NewForm(c *fiber.Ctx) error {
	return c.Render("new", fiber.Map{}, "layouts/base")
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid price")
	}
	if name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}

	if err := h.DB.CreateProduct(c.Context(), name, price); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/")
}

func (h *ProductHandler) EditForm(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	product, err := h.DB.GetProduct(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "product not found")
	}

	return c.Render("edit", fiber.Map{"Product": product}, "layouts/base")
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid price")
	}
	if name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}

	if err := h.DB.UpdateProduct(c.Context(), id, name, price); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/")
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := h.DB.DeleteProduct(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/")
}
