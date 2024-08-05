package rewrite

import (
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	_ "log"
	"net/http"
)

// AppendResourcesOptions is a struct containing configuration options for the `AppendResourcesHandler` method.
type AppendResourcesOptions struct {
	// A list of JavaScript URIs to append to an HTML document's `<head>` element as `<script>` tags.
	JavaScript []string
	// A list of CSS URIs to append to an HTML document's `<head>` element as `<link rel="stylesheet">` tags.
	Stylesheets []string
	// A dictionary of key and value pairs to append to an HTML document's <body> element as `data-{KEY}="{VALUE}` attributes.
	DataAttributes map[string]string
	// AppendJavaScriptAtEOF is a boolean flag to append JavaScript markup at the end of an HTML document
	// rather than in the <head> HTML element. Default is false
	AppendJavaScriptAtEOF bool
}

// AppendResourcesHandler() creates a `RewriteHTMLFunc` callback function, configured by 'opts', and uses that
// callback function and 'previous_handler' to invoke the `RewriteHTMLHandler` function. All of this will cause
// the output of 'previous_handler' to be rewritten to append headers and data attributes defined in 'opts'.
func AppendResourcesHandler(previous_handler http.Handler, opts *AppendResourcesOptions) http.Handler {

	var cb RewriteHTMLFunc

	cb = func(n *html.Node, w io.Writer) {

		if n.Type == html.ElementNode && n.Data == "head" {

			if !opts.AppendJavaScriptAtEOF {
				appendJS(n, opts)
			}

			for _, css := range opts.Stylesheets {
				link_type := html.Attribute{"", "type", "text/css"}
				link_rel := html.Attribute{"", "rel", "stylesheet"}
				link_href := html.Attribute{"", "href", css}

				link := html.Node{
					Type:      html.ElementNode,
					DataAtom:  atom.Link,
					Data:      "link",
					Namespace: "",
					Attr:      []html.Attribute{link_type, link_rel, link_href},
				}

				n.AppendChild(&link)
			}
		}

		if n.Type == html.ElementNode && n.Data == "body" {

			for k, v := range opts.DataAttributes {

				data_ns := ""
				data_key := fmt.Sprintf("data-%s", k)
				data_value := v

				data_attr := html.Attribute{data_ns, data_key, data_value}
				n.Attr = append(n.Attr, data_attr)
			}

		}

		if n.Type == html.ElementNode && n.Data == "html" {

			if opts.AppendJavaScriptAtEOF {
				appendJS(n, opts)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cb(c, w)
		}
	}

	return RewriteHTMLHandler(previous_handler, cb)
}

func appendJS(n *html.Node, opts *AppendResourcesOptions) {

	for _, js := range opts.JavaScript {

		script_type := html.Attribute{"", "type", "text/javascript"}
		script_src := html.Attribute{"", "src", js}

		script := html.Node{
			Type:      html.ElementNode,
			DataAtom:  atom.Script,
			Data:      "script",
			Namespace: "",
			Attr:      []html.Attribute{script_type, script_src},
		}

		n.AppendChild(&script)
	}

}
