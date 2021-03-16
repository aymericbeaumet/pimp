package command

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var replCommand = &cli.Command{
	Hidden: true,
	Action: func(c *cli.Context) error {
		const prompt = "pimp> "
		var sb strings.Builder

		fmt.Fprint(c.App.Writer, prompt)

		scanner := bufio.NewScanner(c.App.Reader)
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())

			if len(text) == 0 {
				fmt.Fprint(c.App.Writer, prompt)
				continue
			}

			sb.Reset()
			sb.WriteString(c.String("ldelim"))
			sb.WriteRune(' ')
			sb.WriteString(scanner.Text())
			sb.WriteRune(' ')
			sb.WriteString(c.String("rdelim"))

			rendered, err := render(sb.String(), c.String("ldelim"), c.String("rdelim"))
			if err != nil {
				fmt.Fprintln(c.App.ErrWriter, err)
				fmt.Fprint(c.App.Writer, prompt)
				continue
			}

			fmt.Fprint(c.App.Writer, rendered)

			if !strings.HasSuffix(rendered, "\n") {
				fmt.Fprint(c.App.Writer, "\n")
			}

			fmt.Fprint(c.App.Writer, prompt)
		}

		return scanner.Err()
	},
}
