package server

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"goscouter/internal/subdomain"
)

type subdomainScanResponse struct {
	Domain      string                `json:"domain"`
	HasWildcard bool                  `json:"has_wildcard"`
	Count       int                   `json:"count"`
	Items       []subdomain.Subdomain `json:"items"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type scanResult struct {
	Domain      string
	HasWildcard bool
	Items       []subdomain.Subdomain
}

func subdomainScanHandler(finder *subdomain.Finder, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := strings.TrimSpace(c.Query("domain"))
		if domain == "" {
			c.JSON(http.StatusBadRequest, errorResponse{Error: "domain query parameter is required"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		result, err := runSubdomainScan(ctx, finder, domain)
		if err != nil {
			status := http.StatusBadGateway
			if errors.Is(err, subdomain.ErrInvalidDomain) {
				status = http.StatusBadRequest
			} else if errors.Is(err, context.DeadlineExceeded) {
				status = http.StatusGatewayTimeout
			}
			c.JSON(status, errorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, subdomainScanResponse{
			Domain:      result.Domain,
			HasWildcard: result.HasWildcard,
			Count:       len(result.Items),
			Items:       result.Items,
		})
	}
}

func runSubdomainScan(ctx context.Context, finder *subdomain.Finder, domain string) (scanResult, error) {
	results, hasWildcard, err := finder.Find(ctx, domain)
	if err != nil {
		return scanResult{}, err
	}

	items := make([]subdomain.Subdomain, 0, len(results))
	for _, item := range results {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	normalizedDomain := strings.TrimSuffix(strings.ToLower(domain), ".")
	return scanResult{
		Domain:      normalizedDomain,
		HasWildcard: hasWildcard,
		Items:       items,
	}, nil
}
