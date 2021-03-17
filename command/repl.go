package command

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/aymericbeaumet/pimp/script"
	prompt "github.com/c-bata/go-prompt"
	"github.com/sahilm/fuzzy"
	"github.com/urfave/cli/v2"
)

var identifierRegexp = regexp.MustCompilePOSIX("^[a-zA-Z][a-zA-Z0-9]*$")

// https://golang.org/pkg/text/template/#hdr-Functions
var nativeTemplateFunctions = []string{
	"and",
	"call",
	"eq",
	"ge",
	"gt",
	"html",
	"index",
	"js",
	"le",
	"len",
	"lt",
	"ne",
	"not",
	"or",
	"print",
	"printf",
	"println",
	"slice",
	"urlquery",
}

var replCommand = &cli.Command{
	Hidden: true,
	Action: func(c *cli.Context) error {
		completionCandidates := make([]string, 0, len(nativeTemplateFunctions)+len(funcmap))
		completionCandidates = append(completionCandidates, nativeTemplateFunctions...)
		for templateFunc := range funcmap {
			completionCandidates = append(completionCandidates, templateFunc)
		}

		executor := func(input string) {
			input = strings.TrimSpace(input)
			if len(input) == 0 {
				return
			}

			var out strings.Builder
			if err := script.Run(&out, input, c.String("ldelim"), c.String("rdelim"), funcmap); err != nil {
				fmt.Fprintln(c.App.ErrWriter, err)
				return
			}

			if !strings.HasSuffix(out.String(), "\n") {
				fmt.Fprintln(c.App.Writer, out.String())
			} else {
				fmt.Fprint(c.App.Writer, out.String())
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
			matches := fuzzy.Find(pattern, completionCandidates)

			out := make([]prompt.Suggest, 0, len(matches))
			for _, match := range matches {
				var description string
				if templateFunc, ok := funcmap[match.Str]; ok {
					description = reflect.TypeOf(templateFunc).String()
				} else {
					description = "predefined global function"
				}
				out = append(out, prompt.Suggest{
					Text:        match.Str,
					Description: description,
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
