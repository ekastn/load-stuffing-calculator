package handler

import (
	"net/http"
	"strconv"

	"github.com/ekastn/load-stuffing-calculator/internal/dto"
	"github.com/ekastn/load-stuffing-calculator/internal/response"
	"github.com/ekastn/load-stuffing-calculator/internal/service"
	"github.com/gin-gonic/gin"
)

type ContainerHandler struct {
	containerSvc service.ContainerService
}

func NewContainerHandler(containerSvc service.ContainerService) *ContainerHandler {
	return &ContainerHandler{containerSvc: containerSvc}
}

// CreateContainer godoc
// @Summary      Create a new container
// @Description  Creates a new container type. Requires admin privileges.
// @Tags         containers
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateContainerRequest true "Container Creation Data"
// @Success      201  {object}  response.APIResponse{data=dto.ContainerResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /containers [post]
func (h *ContainerHandler) CreateContainer(c *gin.Context) {
	var req dto.CreateContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	resp, err := h.containerSvc.CreateContainer(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create container: "+err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetContainer godoc
// @Summary      Get a container by ID
// @Description  Retrieves container details by ID.
// @Tags         containers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Container ID"
// @Success      200  {object}  response.APIResponse{data=dto.ContainerResponse}
// @Failure      400  {object}  response.APIResponse
// @Failure      404  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /containers/{id} [get]
func (h *ContainerHandler) GetContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Container ID is required")
		return
	}

	resp, err := h.containerSvc.GetContainer(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Container not found")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// ListContainers godoc
// @Summary      List containers
// @Description  Retrieves a paginated list of containers.
// @Tags         containers
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number" default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Success      200  {object}  response.APIResponse{data=[]dto.ContainerResponse}
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /containers [get]
func (h *ContainerHandler) ListContainers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.containerSvc.ListContainers(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list containers")
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// UpdateContainer godoc
// @Summary      Update a container
// @Description  Updates an existing container. Requires admin privileges.
// @Tags         containers
// @Accept       json
// @Produce      json
// @Param        id      path      string                 true  "Container ID"
// @Param        request body      dto.UpdateContainerRequest  true  "Container Update Data"
// @Success      200     {object}  response.APIResponse
// @Failure      400     {object}  response.APIResponse
// @Failure      500     {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /containers/{id} [put]
func (h *ContainerHandler) UpdateContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Container ID is required")
		return
	}

	var req dto.UpdateContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	err := h.containerSvc.UpdateContainer(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update container: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}

// DeleteContainer godoc
// @Summary      Delete a container
// @Description  Deletes a container by ID. Requires admin privileges.
// @Tags         containers
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Container ID"
// @Success      200  {object}  response.APIResponse
// @Failure      400  {object}  response.APIResponse
// @Failure      500  {object}  response.APIResponse
// @Security     BearerAuth
// @Router       /containers/{id} [delete]
func (h *ContainerHandler) DeleteContainer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, http.StatusBadRequest, "Container ID is required")
		return
	}

	err := h.containerSvc.DeleteContainer(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete container: "+err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}
