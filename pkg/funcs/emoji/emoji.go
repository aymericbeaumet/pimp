package emoji

import "text/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"EmojiReplace": Replace,
	}
}
