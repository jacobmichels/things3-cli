package cli

import (
	"fmt"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
	"github.com/spf13/cobra"
)

// NewUpdateHeadingCommand builds the update-heading subcommand.
func NewUpdateHeadingCommand(app *App) *cobra.Command {
	var dbPath string
	var id string
	var projectID string
	var project string
	var newTitle string

	cmd := &cobra.Command{
		Use:   "update-heading [OPTIONS...] [--] [-|TITLE]",
		Short: "Update an existing heading",
		RunE: func(cmd *cobra.Command, args []string) error {
			rawInput, err := readInput(app.In, args)
			if err != nil {
				return err
			}
			currentTitle := extractTitle(rawInput, "")

			if id == "" && currentTitle == "" {
				return fmt.Errorf("Error: Must specify --id=ID or heading title")
			}
			if strings.TrimSpace(newTitle) == "" {
				return fmt.Errorf("Error: Must specify --title")
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
				headingID, err = store.ResolveHeadingID(currentTitle, resolvedProjectID)
				if err != nil {
					return fmt.Errorf("Error: %s", err)
				}
			}

			if err := store.RenameHeading(headingID, strings.TrimSpace(newTitle)); err != nil {
				return fmt.Errorf("Error: %s", err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&dbPath, "db", "d", "", "Path to Things database (overrides THINGSDB)")
	flags.StringVar(&dbPath, "database", "", "Alias for --db")
	flags.StringVar(&id, "id", "", "ID of the heading to update")
	flags.StringVar(&projectID, "project-id", "", "Project ID (required when targeting by title)")
	flags.StringVar(&project, "project", "", "Project title (required when targeting by title)")
	flags.StringVar(&newTitle, "title", "", "New title for the heading")

	return cmd
}
