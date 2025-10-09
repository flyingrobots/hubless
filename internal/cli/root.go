package cli

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	mockdata "github.com/flyingrobots/hubless/internal/mock"
	tui "github.com/flyingrobots/hubless/internal/ui/tui/mock"
)

// NewRootCommand wires the Fang-powered Cobra tree for mocked Hubless flows.
func NewRootCommand() *cobra.Command {
	env := loadMockEnvironment()

	cmd := &cobra.Command{
		Use:   "hubless",
		Short: "Hubless CLI (mocked wireframes)",
		Long:  "Hubless CLI playground showcasing mocked data and TUI wireframes.",
	}

	cmd.AddCommand(newListCommand(env))
	cmd.AddCommand(newViewCommand(env))
	cmd.AddCommand(newKanbanCommand(env))
	cmd.AddCommand(newTUICommand())
	cmd.AddCommand(newSyncCommand())
	cmd.AddCommand(newAssignCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newCommentCommand())
	cmd.AddCommand(newCreateCommand())

	return cmd
}

type mockEnvironment struct {
	Issues []mockdata.Issue
	Board  []mockdata.BoardColumn
	Now    time.Time
}

func loadMockEnvironment() mockEnvironment {
	now := time.Now()
	return mockEnvironment{
		Issues: mockdata.MockCatalog(now),
		Board:  mockdata.MockBoard(),
		Now:    now,
	}
}

func newListCommand(env mockEnvironment) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List mocked issues",
		Run: func(cmd *cobra.Command, args []string) {
			out := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
			fmt.Fprintln(out, "ID\tSTATUS\tPRIORITY\tASSIGNEE\tUPDATED")
			for _, issue := range env.Issues {
				fmt.Fprintf(out, "%s\t%s\t%s\t%s\t%s\n",
					issue.ID,
					issue.Status,
					issue.Priority,
					issue.Assignee,
					formatAgo(issue.LastUpdated),
				)
			}
			_ = out.Flush()
		},
	}
	return cmd
}

func newViewCommand(env mockEnvironment) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show mocked issue detail",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			for _, issue := range env.Issues {
				if issue.ID == id {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\nStatus: %s\nPriority: %s\nAssignee: %s\nUpdated: %s\n\n%s\n",
						issue.Title,
						issue.Status,
						issue.Priority,
						issue.Assignee,
						formatAgo(issue.LastUpdated),
						issue.Body,
					)
					fmt.Fprintln(cmd.OutOrStdout(), "\nTimeline:")
					for _, evt := range issue.Events {
						fmt.Fprintf(cmd.OutOrStdout(), "- %-16s %-10s %s\n", evt.Label, formatAgo(evt.Timestamp), evt.Note)
					}
					return nil
				}
			}
			return fmt.Errorf("unknown issue: %s", id)
		},
	}
	return cmd
}

func newKanbanCommand(env mockEnvironment) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kanban",
		Short: "Render mocked kanban columns",
		Run: func(cmd *cobra.Command, args []string) {
			for _, column := range env.Board {
				fmt.Fprintf(cmd.OutOrStdout(), "\n%s (limit %d)\n", column.Name, column.Limit)
				for _, card := range column.Issues {
					fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %-32s @%-10s %s\n", card.ID, card.Title, card.Assignee, card.Priority)
				}
			}
		},
	}
	return cmd
}

func newTUICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch the mocked Bubbletea UI",
		RunE: func(cmd *cobra.Command, args []string) error {
			program := tui.NewProgram()
			_, err := program.Run()
			return err
		},
	}
	return cmd
}

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Pretend to synchronise events",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "[mock] fetch   · refs/hubless/issues → ok")
			fmt.Fprintln(out, "[mock] apply   · catalog indexes      → ok")
			fmt.Fprintln(out, "[mock] project · snapshots rebuilt   → ok")
			fmt.Fprintln(out, "Mock sync complete")
		},
	}
	return cmd
}

func newAssignCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign <id> --to <user>",
		Short: "Mock assignment mutation",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			to, _ := cmd.Flags().GetString("to")
			fmt.Fprintf(cmd.OutOrStdout(), "Would assign %s to %s (mock)\n", args[0], to)
		},
	}
	cmd.Flags().String("to", "", "target assignee")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <id> --to <state>",
		Short: "Mock status transition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			to, _ := cmd.Flags().GetString("to")
			fmt.Fprintf(cmd.OutOrStdout(), "Would set status of %s to %s (mock)\n", args[0], to)
		},
	}
	cmd.Flags().String("to", "", "target status")
	_ = cmd.MarkFlagRequired("to")
	return cmd
}

func newCommentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment <id>",
		Short: "Mock comment creation",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Would open $EDITOR for comment on %s (mock)\n", args[0])
		},
	}
	return cmd
}

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Mock issue creation",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "Would launch $EDITOR and emit issue:created (mock)")
		},
	}
	return cmd
}

func formatAgo(t time.Time) string {
	if t.IsZero() {
		return "n/a"
	}
	diff := time.Since(t)
	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		return fmt.Sprintf("%dm", int(diff/time.Minute))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh", int(diff/time.Hour))
	}
	return fmt.Sprintf("%dd", int(diff/(24*time.Hour)))
}
