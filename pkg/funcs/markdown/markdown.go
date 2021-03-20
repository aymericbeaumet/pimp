package markdown

import "text/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"MarkdownRender": Render,
		"MarkdownTOC":    TOC,
	}
}
