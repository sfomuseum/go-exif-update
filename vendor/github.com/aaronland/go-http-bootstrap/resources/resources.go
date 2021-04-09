package resources

// this will likely get moved in to it's own package or merged
// with go-http-rewrite below (20190723/thisisaaronland)

import (
	"fmt"
	"github.com/aaronland/go-http-rewrite"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	_ "log"
	"net/http"
)

type AppendResourcesOptions struct {
	JS             []string
	CSS            []string
	DataAttributes map[string]string
}

func AppendResourcesHandler(next http.Handler, opts *AppendResourcesOptions) http.Handler {

	var cb rewrite.RewriteHTMLFunc

	cb = func(n *html.Node, w io.Writer) {

		if n.Type == html.ElementNode && n.Data == "head" {

			for _, js := range opts.JS {

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

			for _, css := range opts.CSS {
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

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			cb(c, w)
		}
	}

	return rewrite.RewriteHTMLHandler(next, cb)
}
