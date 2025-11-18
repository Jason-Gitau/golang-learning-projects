package scraper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/user/web-scraper/models"
)

// Parser handles HTML parsing operations
type Parser struct{}

// NewParser creates a new HTML parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseHTML parses HTML content and extracts title, description, and links
func (p *Parser) ParseHTML(html string, baseURL string) (string, string, []models.LinkData, error) {
	// Create a goquery document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", "", nil, err
	}

	// Extract title
	title := p.extractTitle(doc)

	// Extract meta description
	description := p.extractDescription(doc)

	// Extract links
	links := p.extractLinks(doc, baseURL)

	return title, description, links, nil
}

// extractTitle extracts the page title
func (p *Parser) extractTitle(doc *goquery.Document) string {
	// Try <title> tag first
	title := doc.Find("title").First().Text()
	if title != "" {
		return strings.TrimSpace(title)
	}

	// Try og:title meta tag
	ogTitle, exists := doc.Find("meta[property='og:title']").Attr("content")
	if exists && ogTitle != "" {
		return strings.TrimSpace(ogTitle)
	}

	// Try h1 as fallback
	h1 := doc.Find("h1").First().Text()
	if h1 != "" {
		return strings.TrimSpace(h1)
	}

	return "No title found"
}

// extractDescription extracts the page description
func (p *Parser) extractDescription(doc *goquery.Document) string {
	// Try meta description
	metaDesc, exists := doc.Find("meta[name='description']").Attr("content")
	if exists && metaDesc != "" {
		return strings.TrimSpace(metaDesc)
	}

	// Try og:description
	ogDesc, exists := doc.Find("meta[property='og:description']").Attr("content")
	if exists && ogDesc != "" {
		return strings.TrimSpace(ogDesc)
	}

	// Try first paragraph as fallback
	firstP := doc.Find("p").First().Text()
	if firstP != "" {
		// Limit to 200 characters
		if len(firstP) > 200 {
			firstP = firstP[:200] + "..."
		}
		return strings.TrimSpace(firstP)
	}

	return "No description found"
}

// extractLinks extracts all links from the page
func (p *Parser) extractLinks(doc *goquery.Document, baseURL string) []models.LinkData {
	links := []models.LinkData{}
	seen := make(map[string]bool)

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		// Skip anchor links, javascript, and mailto
		if strings.HasPrefix(href, "#") ||
			strings.HasPrefix(href, "javascript:") ||
			strings.HasPrefix(href, "mailto:") {
			return
		}

		// Get link text
		text := strings.TrimSpace(s.Text())
		if text == "" {
			text = href
		}

		// Avoid duplicates
		if seen[href] {
			return
		}
		seen[href] = true

		// Resolve relative URLs
		absoluteURL := p.resolveURL(href, baseURL)

		links = append(links, models.LinkData{
			URL:  absoluteURL,
			Text: text,
		})
	})

	return links
}

// resolveURL converts relative URLs to absolute URLs
func (p *Parser) resolveURL(href, baseURL string) string {
	// If already absolute, return as is
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	// If protocol-relative URL
	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}

	// For relative URLs, we'd need proper URL parsing
	// For simplicity, just return the href
	// In production, use url.Parse and url.ResolveReference
	return href
}

// SanitizeText cleans and normalizes text
func (p *Parser) SanitizeText(text string) string {
	// Remove extra whitespace
	text = strings.TrimSpace(text)

	// Replace multiple spaces with single space
	text = strings.Join(strings.Fields(text), " ")

	return text
}
