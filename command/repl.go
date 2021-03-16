package command

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/sahilm/fuzzy"
	"github.com/urfave/cli/v2"
)

var identifierRegexp = regexp.MustCompilePOSIX("^[a-zA-Z][a-zA-Z0-9]*$")

var replCommand = &cli.Command{
	Hidden: true,
	Action: func(c *cli.Context) error {
		var sb strings.Builder

		executor := func(input string) {
			input = strings.TrimSpace(input)
			if len(input) == 0 {
				return
			}

			sb.Reset()
			sb.WriteString(c.String("ldelim"))
			sb.WriteRune(' ')
			sb.WriteString(strings.TrimSpace(input))
			sb.WriteRune(' ')
			sb.WriteString(c.String("rdelim"))

			rendered, err := render(sb.String(), c.String("ldelim"), c.String("rdelim"))
			if err != nil {
				fmt.Fprintln(c.App.ErrWriter, err)
			} else {
				fmt.Fprint(c.App.Writer, rendered)
				if !strings.HasSuffix(rendered, "\n") {
					fmt.Fprint(c.App.Writer, "\n")
				}
			}
		}

		completer := func(document prompt.Document) []prompt.Suggest {
			if strings.HasSuffix(document.Text, " ") {
				return nil
			}

			args := strings.Fields(document.Text)
			if len(args) == 0 {
				return nil
			}
			lastArg := args[len(args)-1]
			if !identifierRegexp.MatchString(lastArg) {
				return nil
			}

			pattern := lastArg
			data := make([]string, 0, len(fm))
			for fn := range fm {
				data = append(data, fn)
			}
			matches := fuzzy.Find(pattern, data)

			out := make([]prompt.Suggest, 0, len(matches))
			for _, match := range matches {
				out = append(out, prompt.Suggest{
					Text:        match.Str,
					Description: reflect.TypeOf(fm[match.Str]).String(),
				})
			}
			return out
		}

		p := prompt.New(
			executor,
			completer,
			prompt.OptionTitle("pimp"),
			prompt.OptionPrefix("pimp> "),
		)
		p.Run()

		return nil
	},
}
