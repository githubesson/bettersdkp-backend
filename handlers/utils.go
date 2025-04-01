package handlers

import (
	"regexp"
)

func ExtractPoints(text string) string {
	re := regexp.MustCompile(`\(max\. punktów (\d+)\)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return "0"
}

func SetCommonHeaders(headers map[string]string) map[string]string {
	commonHeaders := map[string]string{
		"Accept":             "*/*",
		"Accept-Language":    "pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":         "keep-alive",
		"Origin":             "https://sdkp.pjwstk.edu.pl",
		"Referer":            "https://sdkp.pjwstk.edu.pl/",
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
		"X-Requested-With":   "XMLHttpRequest",
		"sec-ch-ua":          `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
	}

	for k, v := range headers {
		commonHeaders[k] = v
	}

	return commonHeaders
}
