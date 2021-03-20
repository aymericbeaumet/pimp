package emoji

import (
	"github.com/kyokomi/emoji/v2"
)

func Replace(in string) string {
	return emoji.Sprint(in)
}
