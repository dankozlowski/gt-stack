package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dankoz/gt-stacks/internal/core"
	"github.com/spf13/cobra"
)

func interactivePicker(prompt string, options []string) (string, error) {
	fmt.Println(prompt)
	for i, o := range options {
		fmt.Printf("  %d) %s\n", i+1, o)
	}
	fmt.Print("> ")
	r := bufio.NewReader(os.Stdin)
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	idx, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil || idx < 1 || idx > len(options) {
		return "", fmt.Errorf("invalid selection %q", strings.TrimSpace(line))
	}
	return options[idx-1], nil
}

func newUpCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:     "up [n]",
		Aliases: []string{"u"},
		Short:   "Checkout child (interactive picker on multiple)",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n := 1
			if len(args) == 1 {
				v, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				n = v
			}
			out, err := c.CheckoutUp(cmd.Context(), n, interactivePicker)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "→", out)
			return nil
		},
	}
}

func newDownCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:     "down [n]",
		Aliases: []string{"d"},
		Short:   "Checkout parent",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n := 1
			if len(args) == 1 {
				v, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				n = v
			}
			out, err := c.CheckoutDown(cmd.Context(), n)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "→", out)
			return nil
		},
	}
}

func newCheckoutCmd(c *core.Core) *cobra.Command {
	return &cobra.Command{
		Use:     "checkout [branch]",
		Aliases: []string{"co"},
		Short:   "Checkout a branch (interactive picker if not specified)",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return c.Checkout(cmd.Context(), args[0])
			}
			s, err := c.LoadStack(cmd.Context())
			if err != nil {
				return err
			}
			var opts []string
			for name, br := range s.Branches {
				if br.Tracked {
					opts = append(opts, name)
				}
			}
			chosen, err := interactivePicker("checkout which branch?", opts)
			if err != nil {
				return err
			}
			return c.Checkout(cmd.Context(), chosen)
		},
	}
}
