package asi67

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

// Constants for the API logic
const (
	itemsPerPage = 12
)

// --- JSON DATA STRUCTURES ---

type APIResponse struct {
	Data struct {
		ProdCount int                      `json:"prodCount"` // Total number of items
		ProdID    map[string]ListingDetail `json:"prodId"`    // Listings map
	} `json:"data"`
}

type ListingDetail struct {
	City         string  `json:"city"`
	Cp           string  `json:"cp"`
	Surface      float64 `json:"surface"`
	ProdRef      string  `json:"prod_ref"`
	RentTotal    float64 `json:"rent_total"`
	PricePrimary float64 `json:"pricePrimary"`
	Title        struct {
		Fr string `json:"fr"`
	} `json:"title"`
}

// --- PROVIDER IMPLEMENTATION ---

type Provider struct {
	apiURL string
	client *http.Client
}

func NewProvider(apiURL string) *Provider {
	return &Provider{
		apiURL: apiURL,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *Provider) Name() string {
	return "asi67 (api-client-v2)"
}

func (p *Provider) FetchListings(ctx context.Context) ([]core.Listing, error) {
	// Step 1: Fetch the first page synchronously to discover the total count
	firstPageListings, totalCount, err := p.fetchPage(ctx, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page 1: %w", err)
	}

	// If no items or only one page, return immediately
	if totalCount <= itemsPerPage {
		return firstPageListings, nil
	}

	// Step 2: Calculate total pages needed
	totalPages := int(math.Ceil(float64(totalCount) / float64(itemsPerPage)))
	fmt.Printf("DEBUG: Found %d total listings (%d pages). Starting parallel fetch...\n", totalCount, totalPages)

	// Step 3: Fan-out - Fetch remaining pages in parallel
	var (
		wg          sync.WaitGroup
		mu          sync.Mutex
		allListings = firstPageListings // Start with page 1 results
		errOccurred error
	)

	// Loop from page 2 to totalPages
	for page := 2; page <= totalPages; page++ {
		wg.Add(1)
		go func(pNum int) {
			defer wg.Done()

			// Check context cancellation before making request
			if ctx.Err() != nil {
				return
			}

			listings, _, err := p.fetchPage(ctx, pNum)
			if err != nil {
				// We log the error but don't stop the whole process (best effort)
				fmt.Printf("⚠️ Error fetching page %d: %v\n", pNum, err)
				return
			}

			// Thread-safe append
			mu.Lock()
			allListings = append(allListings, listings...)
			mu.Unlock()
		}(page)
	}

	wg.Wait()

	if nil != errOccurred {
		return nil, errOccurred
	}

	fmt.Printf("DEBUG: Successfully fetched %d listings from %d pages.\n", len(allListings), totalPages)
	return allListings, nil
}

// fetchPage handles the API call for a specific page number.
// It returns the listings, the total count (from metadata), and an error.
func (p *Provider) fetchPage(ctx context.Context, pageNum int) ([]core.Listing, int, error) {
	// 1. Prepare Request Payload with the specific page number
	requestBody := map[string]interface{}{
		"params": map[string]interface{}{
			"type_offer": "2",
			"prod_type":  "appt",
			"geo":        "strasbourg/67000",
			"query": map[string]interface{}{
				"page":                 strconv.Itoa(pageNum), // Inject page number here
				"prod.prod_type":       "appt",
				"prod.geo":             "strasbourg/67000",
				"prod.geo_radius":      "20",
				"prod.budget_rent_max": "1000",
			},
		},
	}

	jsonPayload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, 0, fmt.Errorf("json marshal error: %w", err)
	}

	// 2. Create HTTP Request
	req, err := http.NewRequestWithContext(ctx, "POST", p.apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, 0, fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ClassifiedsWatcher/1.0)")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 3. Execute
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http call error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("api status %d", resp.StatusCode)
	}

	// 4. Decode
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, 0, fmt.Errorf("json decode error: %w", err)
	}

	// 5. Map to Domain
	var listings []core.Listing
	for id, detail := range apiResp.Data.ProdID {
		fullURL := fmt.Sprintf("https://www.asi67.com/location/location,%s", id)

		title := detail.Title.Fr
		if title == "" {
			title = fmt.Sprintf("Apartment in %s (%s)", detail.City, detail.Cp)
		}

		price := detail.RentTotal
		if price == 0 {
			price = detail.PricePrimary
		}

		listings = append(listings, core.Listing{
			ID:          id,
			Title:       title,
			Url:         fullURL,
			Price:       price,
			Currency:    "EUR",
			Description: fmt.Sprintf("%.0f m² - %s", detail.Surface, detail.City),
			PublishedAt: time.Now(),
		})
	}

	return listings, apiResp.Data.ProdCount, nil
}
