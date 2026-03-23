package cli

import (
	"fmt"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

// NewDeleteHeadingCommand builds the delete-heading subcommand.
func NewDeleteHeadingCommand(app *App) *cobra.Command {
	var dbPath string
	var id string
	var projectID string
	var project string
	var confirm string

	cmd := &cobra.Command{
		Use:   "delete-heading [OPTIONS...] [--] [-|TITLE]",
		Short: "Delete an existing heading",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}
			title := extractTitle(rawInput, "")

			if id == "" && title == "" {
				return fmt.Errorf("Error: Must specify --id=ID or heading title")
			}

			target := deleteConfirmTarget(id, rawInput)
			if err := confirmDelete(app, "heading", target, confirm); err != nil {
				return err
			}

			store, _, err := db.OpenDefaultWritable(dbPath)
			if err != nil {
				return formatDBError(err)
			}
			defer store.Close()

			headingID := id
			if headingID == "" {
				resolvedProjectID := projectID
				if project != "" && resolvedProjectID == "" {
					resolvedProjectID, err = store.ResolveProjectID(project)
					if err != nil {
						return fmt.Errorf("Error: %s", err)
					}
				}
				headingID, err = store.ResolveHeadingID(title, resolvedProjectID)
				if err != nil {
					return fmt.Errorf("Error: %s", err)
				}
			}

			if err := store.DeleteHeading(headingID); err != nil {
				return fmt.Errorf("Error: %s", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	flags.StringVar(&dbPath, "database", "", "Alias for --db")
	flags.StringVar(&id, "id", "", "ID of the heading to delete")
	flags.StringVar(&projectID, "project-id", "", "Project ID (required when targeting by title)")
	flags.StringVar(&project, "project", "", "Project title (required when targeting by title)")
	flags.StringVar(&confirm, "confirm", "", "Confirm deletion by typing the heading ID or title")

	return cmd
}
