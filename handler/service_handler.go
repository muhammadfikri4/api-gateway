package handler

import (
	"net/http"
	"strconv"
	"strings"

	"api-gateway/model"
	"api-gateway/repository"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	serviceRepo *repository.ServiceRepository
}

func NewServiceHandler(serviceRepo *repository.ServiceRepository) *ServiceHandler {
	return &ServiceHandler{serviceRepo: serviceRepo}
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var req model.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: err.Error()})
		return
	}

	req.RoutePrefix = normalizePrefix(req.RoutePrefix)

	service := &model.Service{
		Name:        req.Name,
		Description: req.Description,
		TargetURL:   strings.TrimRight(req.TargetURL, "/"),
		RoutePrefix: req.RoutePrefix,
		IsActive:    true,
	}

	if err := h.serviceRepo.Create(service); err != nil {
		c.JSON(http.StatusConflict, model.Response{Success: false, Message: "Service name or route prefix already exists"})
		return
	}

	c.JSON(http.StatusCreated, model.Response{
		Success: true,
		Message: "Service created successfully",
		Data:    service,
	})
}

func (h *ServiceHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	services, total, err := h.serviceRepo.FindAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to fetch services"})
		return
	}

	c.JSON(http.StatusOK, model.PaginatedResponse{
		Success: true,
		Message: "Services fetched successfully",
		Data:    services,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: "Invalid ID"})
		return
	}

	service, err := h.serviceRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{Success: false, Message: "Service not found"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Success: true, Data: service})
}

func (h *ServiceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: "Invalid ID"})
		return
	}

	service, err := h.serviceRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.Response{Success: false, Message: "Service not found"})
		return
	}

	var req model.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: err.Error()})
		return
	}

	if req.Name != "" {
		service.Name = req.Name
	}
	if req.Description != "" {
		service.Description = req.Description
	}
	if req.TargetURL != "" {
		service.TargetURL = strings.TrimRight(req.TargetURL, "/")
	}
	if req.RoutePrefix != "" {
		service.RoutePrefix = normalizePrefix(req.RoutePrefix)
	}
	if req.IsActive != nil {
		service.IsActive = *req.IsActive
	}

	if err := h.serviceRepo.Update(service); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Success: true, Message: "Service updated successfully", Data: service})
}

func (h *ServiceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Response{Success: false, Message: "Invalid ID"})
		return
	}

	if err := h.serviceRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{Success: false, Message: "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, model.Response{Success: true, Message: "Service deleted successfully"})
}

func normalizePrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	prefix = strings.TrimRight(prefix, "/")
	return prefix
}
