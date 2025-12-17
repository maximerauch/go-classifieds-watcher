package rememberme

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/maximerauch/go-classifieds-watcher/internal/core"
)

type Provider struct {
	client    *http.Client
	searchURL string
}

func NewProvider(searchURL string) *Provider {
	return &Provider{
		client: &http.Client{
			// Increased timeout as HTML scraping is slower than JSON APIs
			Timeout: 30 * time.Second,
		},
		searchURL: searchURL,
	}
}

func (p *Provider) Name() string {
	return "remember-me-france"
}

func (p *Provider) FetchItems(ctx context.Context) ([]core.Item, error) {
	// Fetch Page 1 to discover total pages and initial items
	doc, err := p.fetchDocument(ctx, p.searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch first page: %w", err)
	}

	// Detect max page number from pagination
	maxPage := 1
	doc.Find("a.page-numbers").Each(func(i int, s *goquery.Selection) {
		txt := strings.TrimSpace(s.Text())
		if pageNum, err := strconv.Atoi(txt); err == nil {
			if pageNum > maxPage {
				maxPage = pageNum
			}
		}
	})

	fmt.Printf("DEBUG: Found max page: %d\n", maxPage)

	// Extract items from page 1
	allItems := p.extractItems(doc)

	if maxPage <= 1 {
		return allItems, nil
	}

	// Fan-Out: Scrape remaining pages in parallel
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	// RATE LIMITING: Crucial to avoid IP ban.
	// Use a buffered channel as a semaphore to limit concurrency to 5 workers.
	sem := make(chan struct{}, 5)

	for page := 2; page <= maxPage; page++ {
		wg.Add(1)
		go func(pNum int) {
			defer wg.Done()

			// Acquire token (blocks if 5 workers are already active)
			sem <- struct{}{}
			defer func() { <-sem }() // Release token

			if ctx.Err() != nil {
				return
			}

			// Construct the paginated URL dynamically based on the initial searchURL
			targetURL, err := p.buildPageURL(p.searchURL, pNum)
			if err != nil {
				fmt.Printf("⚠️ Error building URL for page %d: %v\n", pNum, err)
				return
			}

			pageDoc, err := p.fetchDocument(ctx, targetURL)
			if err != nil {
				// Best effort: Log error but continue processing other pages
				fmt.Printf("⚠️ Error scraping page %d: %v\n", pNum, err)
				return
			}

			items := p.extractItems(pageDoc)

			// Thread-safe append
			mu.Lock()
			allItems = append(allItems, items...)
			mu.Unlock()
		}(page)
	}

	wg.Wait()

	fmt.Printf("DEBUG: Total items found: %d\n", len(allItems))
	return allItems, nil
}

// buildPageURL injects "/page/N/" into the URL path for pagination.
// Input: https://site.org/pets/?filter=1 -> Output: https://site.org/pets/page/2/?filter=1
func (p *Provider) buildPageURL(originalURL string, pageNum int) (string, error) {
	u, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	// Append "page/X" to the existing path
	// path.Join cleans up the path, so we ensure it ends correctly
	newPath := path.Join(u.Path, "page", strconv.Itoa(pageNum))

	// WordPress often requires a trailing slash for pagination to work correctly before query params
	if !strings.HasSuffix(newPath, "/") {
		newPath += "/"
	}

	u.Path = newPath
	return u.String(), nil
}

func (p *Provider) fetchDocument(ctx context.Context, url string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent to mimic a standard browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ClassifiedsWatcher/1.0)")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("⚠️ Error closing response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

// extractItems parses the HTML based on the specific "article.pets" structure.
func (p *Provider) extractItems(doc *goquery.Document) []core.Item {
	var items []core.Item

	// CSS SELECTOR: We target <article> tags with class "pets"
	doc.Find("article.pets").Each(func(i int, s *goquery.Selection) {

		// ID: Get the ID from the <article id="pet-XXXX"> attribute
		id, _ := s.Attr("id")

		// Title: Inside <h3 class="pet-title">
		title := strings.TrimSpace(s.Find(".pet-title").Text())

		// URL: Inside <header class="pet-header"> -> <a>
		link, exists := s.Find(".pet-header a").First().Attr("href")
		if !exists {
			return
		}

		// Description: Inside <section class="pet-content">
		rawDesc := s.Find(".pet-content").Text()
		description := strings.TrimSpace(strings.ReplaceAll(rawDesc, "\n", " "))

		if id != "" && title != "" {
			items = append(items, core.Item{
				ID:          id,
				Title:       title,
				Url:         link,
				Price:       0,
				Currency:    "EUR",
				PublishedAt: time.Now(),
				Description: description,
			})
		}
	})

	return items
}
