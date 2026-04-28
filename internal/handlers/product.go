package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"go-fiber-crud/internal/database"
	"go-fiber-crud/internal/models"
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
	success := c.Query("success")
	errorMessage := c.Query("error")
	products, err := h.DB.ListProducts(c.Context(), nameFilter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Render("index", fiber.Map{
		"Products":   products,
		"FilterName": nameFilter,
		"Success":    success,
		"Error":      errorMessage,
	}, "layouts/base")
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
	viewData := fiber.Map{
		"FormName":     name,
		"FormPrice":    priceStr,
		"FormQuantity": quantityStr,
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		viewData["Error"] = "invalid price"
		return c.Status(fiber.StatusBadRequest).Render("new", viewData, "layouts/base")
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		viewData["Error"] = "invalid quantity"
		return c.Status(fiber.StatusBadRequest).Render("new", viewData, "layouts/base")
	}

	if name == "" {
		viewData["Error"] = "name is required"
		return c.Status(fiber.StatusBadRequest).Render("new", viewData, "layouts/base")
	}

	if quantity < 0 {
		viewData["Error"] = "quantity must be 0 or greater"
		return c.Status(fiber.StatusBadRequest).Render("new", viewData, "layouts/base")
	}

	if err := h.DB.CreateProduct(c.Context(), name, price, quantity); err != nil {
		viewData["Error"] = "failed to create product"
		return c.Status(fiber.StatusInternalServerError).Render("new", viewData, "layouts/base")
	}
	return c.Redirect("/?success=Product%20created")
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
	viewData := fiber.Map{
		"FormName":     name,
		"FormPrice":    priceStr,
		"FormQuantity": quantityStr,
		"Product":      models.Product{ID: id, Name: name},
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		viewData["Error"] = "invalid price"
		return c.Status(fiber.StatusBadRequest).Render("edit", viewData, "layouts/base")
	}
	viewData["Product"] = models.Product{ID: id, Name: name, Price: price}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		viewData["Error"] = "invalid quantity"
		return c.Status(fiber.StatusBadRequest).Render("edit", viewData, "layouts/base")
	}
	viewData["Product"] = models.Product{ID: id, Name: name, Price: price, Quantity: quantity}

	if name == "" {
		viewData["Error"] = "name is required"
		return c.Status(fiber.StatusBadRequest).Render("edit", viewData, "layouts/base")
	}

	if quantity < 0 {
		viewData["Error"] = "quantity must be 0 or greater"
		return c.Status(fiber.StatusBadRequest).Render("edit", viewData, "layouts/base")
	}

	if err := h.DB.UpdateProduct(c.Context(), id, name, price, quantity); err != nil {
		viewData["Error"] = "failed to update product"
		return c.Status(fiber.StatusInternalServerError).Render("edit", viewData, "layouts/base")
	}
	return c.Redirect("/?success=Product%20updated")
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
	return c.Redirect("/?success=Product%20deleted")
}
