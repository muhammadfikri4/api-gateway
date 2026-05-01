package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"api-gateway/repository"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	serviceRepo *repository.ServiceRepository
}

func NewProxyHandler(serviceRepo *repository.ServiceRepository) *ProxyHandler {
	return &ProxyHandler{serviceRepo: serviceRepo}
}

func (p *ProxyHandler) Handle(c *gin.Context) {
	path := c.Request.URL.Path

	// Extract the route prefix (first path segment)
	segments := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
	if len(segments) == 0 || segments[0] == "" {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "No route specified"})
		return
	}

	prefix := "/" + segments[0]

	service, err := p.serviceRepo.FindByRoutePrefix(prefix)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No service found for route: " + prefix,
		})
		return
	}

	targetURL, err := url.Parse(service.TargetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Invalid target URL"})
		return
	}

	// Strip the route prefix from the path
	remainingPath := ""
	if len(segments) > 1 {
		remainingPath = "/" + segments[1]
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = targetURL.Path + remainingPath
		req.URL.RawQuery = c.Request.URL.RawQuery
		req.Host = targetURL.Host
		req.Header = c.Request.Header.Clone()
		req.Header.Set("X-Forwarded-Host", c.Request.Host)
		req.Header.Set("X-Forwarded-For", c.ClientIP())
		req.Header.Set("X-Gateway-Service", service.Name)
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error for service %s: %v", service.Name, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"success":false,"message":"Service unavailable: ` + service.Name + `"}`))
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
