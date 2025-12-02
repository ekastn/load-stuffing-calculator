package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productSvc service.ProductService
}

func NewProductHandler(productSvc service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Creates a new product. Requires admin privileges.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProductRequest true "Product Creation Data"
// @Success      201  {object}  response.APIResponse{data=dto.ProductResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.productSvc.CreateProduct(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create product: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetProduct godoc
// @Summary      Get a product by ID
// @Description  Retrieves product details by ID.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  response.APIResponse{data=dto.ProductResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	resp, err := h.productSvc.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Product not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListProducts godoc
// @Summary      List products
// @Description  Retrieves a paginated list of products.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Success      200  {object}  response.APIResponse{data=[]dto.ProductResponse}
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.productSvc.ListProducts(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list products")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Updates an existing product. Requires admin privileges.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Product ID"
// @Param        request body      dto.UpdateProductRequest  true  "Product Update Data"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.productSvc.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update product: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Deletes a product by ID. Requires admin privileges.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Product ID is required")
		return
	}

	err := h.productSvc.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete product: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}
