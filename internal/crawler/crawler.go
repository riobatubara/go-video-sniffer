package crawler

import (
    "net/http"
    "golang.org/x/net/html"
    "io"
    "strings"
    "log"
)

func CrawlPage(url string) ([]string, error) {
    res, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != 200 {
        return nil, err
    }

    return extractVideoLinks(res.Body)
}

func extractVideoLinks(body io.Reader) ([]string, error) {
    var videoLinks []string

    tokenizer := html.NewTokenizer(body)
    for {
        tt := tokenizer.Next()
        switch tt {
        case html.ErrorToken:
            return videoLinks, nil
        case html.StartTagToken, html.SelfClosingTagToken:
            tagName, hasAttr := tokenizer.TagName()
            tag := string(tagName)
            if tag == "video" || tag == "source" {
                for hasAttr {
                    key, val, more := tokenizer.TagAttr()
                    if string(key) == "src" {
                        src := string(val)
                        if isValidVideoLink(src) {
                            videoLinks = append(videoLinks, src)
                        }
                    }
                    hasAttr = more
                }
            }
        }
    }
}

func isValidVideoLink(link string) bool {
    if strings.Contains(link, "ads.") || strings.Contains(link, "doubleclick") {
        return false
    }
    if strings.HasSuffix(link, ".mp4") || strings.Contains(link, ".m3u8") || strings.Contains(link, ".mpd") {
        return true
    }
    return false
}
