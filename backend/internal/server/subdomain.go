package server

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"goscouter/backend/internal/subdomain"
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

type scanViewData struct {
	Domain        string
	Error         string
	HasWildcard   bool
	Count         int
	Items         []subdomain.Subdomain
	UniqueIPs     int
	UniqueIssuers int
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

func subdomainScanUIHandler(finder *subdomain.Finder, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := strings.TrimSpace(c.Query("domain"))
		if domain == "" {
			c.HTML(http.StatusBadRequest, "results.html", scanViewData{
				Error: "Domain query parameter is required.",
			})
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
			c.HTML(status, "results.html", scanViewData{
				Error: err.Error(),
			})
			return
		}

		uniqueIPs := make(map[string]struct{})
		uniqueIssuers := make(map[string]struct{})
		for _, item := range result.Items {
			for _, ip := range item.IPs {
				uniqueIPs[ip] = struct{}{}
			}
			if item.CertIssuer != "" {
				uniqueIssuers[item.CertIssuer] = struct{}{}
			}
		}

		c.HTML(http.StatusOK, "results.html", scanViewData{
			Domain:        result.Domain,
			HasWildcard:   result.HasWildcard,
			Count:         len(result.Items),
			Items:         result.Items,
			UniqueIPs:     len(uniqueIPs),
			UniqueIssuers: len(uniqueIssuers),
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
