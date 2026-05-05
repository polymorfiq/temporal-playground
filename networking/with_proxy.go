package networking

import (
	"log"
	"net/url"
	"os"

	"github.com/chromedp/chromedp"
)

func ProxifiedNavigate(unproxiedUrl string) []chromedp.Action {
	return []chromedp.Action{chromedp.Navigate(ProxifiedUrl(unproxiedUrl))}
}

func ProxifiedUrl(unproxiedUrl string) string {
	if os.Getenv("SCRAPER_PROXY_URL") != "" {
		proxyUrl, err := url.Parse(os.Getenv("SCRAPER_PROXY_URL"))
		if err != nil {
			log.Fatalf("Failed to parse proxy URL: %v", err)
		}

		q := proxyUrl.Query()
		q.Set("url", unproxiedUrl)
		proxyUrl.RawQuery = q.Encode()
		return proxyUrl.String()
	}

	return unproxiedUrl
}
