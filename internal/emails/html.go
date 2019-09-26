package emails

import (
	strip "github.com/grokify/html-strip-tags-go"
)

func stripHTML(html string) string {
	return strip.StripTags(html)
}
