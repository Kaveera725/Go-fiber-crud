package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"go-fiber-crud/internal/database"
)

type ProductHandler struct {
	DB *database.DB
}

// NewProductHandler wires a ProductHandler with its dependencies.
func NewProductHandler(db *database.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

// List renders the product list page.
func (h *ProductHandler) List(c *fiber.Ctx) error {
	nameFilter := c.Query("name")
	products, err := h.DB.ListProducts(c.Context(), nameFilter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Render("index", fiber.Map{"Products": products, "FilterName": nameFilter}, "layouts/base")
}

// NewForm renders the create-product page.
func (h *ProductHandler) NewForm(c *fiber.Ctx) error {
	return c.Render("new", fiber.Map{}, "layouts/base")
}

// Create validates input and inserts a new product.
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	quantityStr := c.FormValue("quantity")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid price")
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid quantity")
	}

	if name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}

	if quantity < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "quantity must be 0 or greater")
	}

	if err := h.DB.CreateProduct(c.Context(), name, price, quantity); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/")
}

// EditForm loads a product and renders the edit page.
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

// Update validates input and updates an existing product.
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	name := c.FormValue("name")
	priceStr := c.FormValue("price")
	quantityStr := c.FormValue("quantity")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid price")
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid quantity")
	}

	if name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}

	if quantity < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "quantity must be 0 or greater")
	}

	if err := h.DB.UpdateProduct(c.Context(), id, name, price, quantity); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Redirect("/")
}

// Delete removes a product and returns to the list.
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
