package cli

import (
	"fmt"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

// NewHeadingsCommand builds the headings command.
func NewHeadingsCommand(app *App) *cobra.Command {
	var dbPath string
	var projectID string
	var project string
	var includeTrashed bool
	var asJSON bool
	var noHeader bool

	cmd := &cobra.Command{
		Use:   "headings",
		Short: "List headings from the Things database",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, _, err := db.OpenDefault(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			if project != "" && projectID == "" {
				projectID, err = store.ResolveProjectID(project)
				if err != nil {
					return fmt.Errorf("Error: %s", err)
				}
			}

			headings, err := store.Headings(db.HeadingFilter{
				ProjectID:             projectID,
				IncludeTrashed:        includeTrashed,
				ExcludeTrashedContext: !includeTrashed,
			})
			if err != nil {
				return formatDBError(err)
			}
			return printHeadings(app.Out, headings, asJSON, noHeader)
		},
	}

	cmd.Flags().StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	cmd.Flags().StringVar(&dbPath, "database", "", "Alias for --db")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Filter by project ID")
	cmd.Flags().StringVar(&project, "project", "", "Filter by project title or ID")
	cmd.Flags().BoolVar(&includeTrashed, "include-trashed", false, "Include trashed headings")
	cmd.Flags().BoolVarP(&asJSON, "json", "j", false, "Output JSON")
	cmd.Flags().BoolVar(&noHeader, "no-header", false, "Suppress header row")

	return cmd
}
