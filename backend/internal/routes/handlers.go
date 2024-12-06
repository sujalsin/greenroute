package routes

import (
	"greenroute/internal/models"
	"greenroute/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RouteHandler handles HTTP requests for route calculations
type RouteHandler struct {
	routeService *services.RouteService
}

// NewRouteHandler creates a new instance of RouteHandler
func NewRouteHandler(routeService *services.RouteService) *RouteHandler {
	return &RouteHandler{
		routeService: routeService,
	}
}

// RegisterRoutes registers all route endpoints
func (h *RouteHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/routes/calculate", h.CalculateRoute)
		v1.GET("/routes/:id", h.GetRoute)
	}
}

// RouteRequest represents the incoming request for route calculation
type RouteRequest struct {
	StartLocation models.Location         `json:"start_location" binding:"required"`
	EndLocation   models.Location         `json:"end_location" binding:"required"`
	Preferences   models.RoutePreferences `json:"preferences"`
}

// CalculateRoute handles the route calculation request
func (h *RouteHandler) CalculateRoute(c *gin.Context) {
	var req RouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	route, err := h.routeService.CalculateRoute(
		req.StartLocation,
		req.EndLocation,
		req.Preferences,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, route)
}

// GetRoute retrieves a previously calculated route
func (h *RouteHandler) GetRoute(c *gin.Context) {
	routeID := c.Param("id")
	// TODO: Implement route retrieval from database
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented yet"})
}
