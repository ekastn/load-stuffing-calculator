package handler

import (
	"net/http"

	"load-stuffing-calculator/internal/dto"
	"load-stuffing-calculator/internal/response"
	"load-stuffing-calculator/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	product, err := h.svc.Create(
		c.Request.Context(),
		req.Label,
		req.SKU,
		req.LengthMm,
		req.WidthMm,
		req.HeightMm,
		req.WeightKg,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.ProductResponse{
		ID:       product.ID.String(),
		Label:    product.Label,
		SKU:      product.Sku,
		LengthMm: product.LengthMm,
		WidthMm:  product.WidthMm,
		HeightMm: product.HeightMm,
		WeightKg: product.WeightKg,
	}

	response.Success(c, http.StatusCreated, resp)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	product, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	resp := dto.ProductResponse{
		ID:       product.ID.String(),
		Label:    product.Label,
		SKU:      product.Sku,
		LengthMm: product.LengthMm,
		WidthMm:  product.WidthMm,
		HeightMm: product.HeightMm,
		WeightKg: product.WeightKg,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list products")
		return
	}

	resp := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		resp = append(resp, dto.ProductResponse{
			ID:       product.ID.String(),
			Label:    product.Label,
			SKU:      product.Sku,
			LengthMm: product.LengthMm,
			WidthMm:  product.WidthMm,
			HeightMm: product.HeightMm,
			WeightKg: product.WeightKg,
		})
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	product, err := h.svc.Update(
		c.Request.Context(),
		id,
		req.Label,
		req.SKU,
		req.LengthMm,
		req.WidthMm,
		req.HeightMm,
		req.WeightKg,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp := dto.ProductResponse{
		ID:       product.ID.String(),
		Label:    product.Label,
		SKU:      product.Sku,
		LengthMm: product.LengthMm,
		WidthMm:  product.WidthMm,
		HeightMm: product.HeightMm,
		WeightKg: product.WeightKg,
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	err = h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}
