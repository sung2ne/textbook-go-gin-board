package sanitize

import (
	"github.com/microcosm-cc/bluemonday"
)

var (
	strictPolicy = bluemonday.StrictPolicy()
	ugcPolicy    = bluemonday.UGCPolicy()
)

func Strict(s string) string {
	return strictPolicy.Sanitize(s)
}

func UGC(s string) string {
	return ugcPolicy.Sanitize(s)
}

func NewCustomPolicy() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowElements("p", "br", "b", "i", "u", "strong", "em")
	p.AllowElements("ul", "ol", "li")
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowAttrs("href").OnElements("a")
	p.RequireNoReferrerOnLinks(true)
	p.AllowAttrs("src", "alt").OnElements("img")
	p.AllowURLSchemes("https")
	return p
}
