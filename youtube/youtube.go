package youtube

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

func GetFirstVideoID(query string) (string, error) {
	searchURL := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", url.QueryEscape(query+" physiotherapy"))

	resp, err := http.Get(searchURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Look for non-ad video renderers
	re := regexp.MustCompile(`"videoRenderer":{"videoId":"(\w+)"`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	for _, match := range matches {
		if len(match) >= 2 {
			// Check if this video section contains ad indicators
			adCheck := regexp.MustCompile(`"adBadge"|"ytd-ad-slot-renderer"|"ytd-promoted-video-renderer"|"sponsored-tag"`)
			if !adCheck.MatchString(match[0]) {
				return match[1], nil
			}
		}
	}

	return "", fmt.Errorf("no non-ad video ID found")
}
