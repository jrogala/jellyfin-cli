package cmd

import (
	"fmt"
	"strings"

	"github.com/jrogala/jellyfin-cli/internal/cmdutil"
	"github.com/jrogala/jellyfin-cli/pkg/ops"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(librariesCmd, moviesCmd, searchCmd, infoCmd, updateCmd,
		identifyCmd, refreshCmd, scanCmd, sessionsCmd, itemsCmd)

	searchCmd.Flags().Int("limit", 20, "Max results")

	itemsCmd.Flags().String("library", "", "Library ID")
	itemsCmd.Flags().String("type", "", "Item type (Movie, Episode, Audio, Series)")
	itemsCmd.Flags().Int("limit", 50, "Max results")

	updateCmd.Flags().String("name", "", "New title")
	updateCmd.Flags().Int("year", 0, "Production year")
	updateCmd.Flags().String("overview", "", "Description/overview")

	identifyCmd.Flags().String("name", "", "Title to search for")
	identifyCmd.Flags().Int("year", 0, "Year")
	identifyCmd.Flags().String("imdb", "", "IMDB ID (e.g. tt1234567)")
	identifyCmd.Flags().String("tmdb", "", "TMDB ID")
	identifyCmd.MarkFlagRequired("name")
}

var librariesCmd = &cobra.Command{
	Use:     "libraries",
	Aliases: []string{"libs"},
	Short:   "List libraries",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		libs, err := ops.ListLibraries(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, libs, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tNAME\tTYPE")
			for _, l := range libs {
				fmt.Fprintf(w, "%s\t%s\t%s\n", l.ID, l.Name, l.CollectionType)
			}
			w.Flush()
		})
		return nil
	},
}

var moviesCmd = &cobra.Command{
	Use:   "movies",
	Short: "List all movies",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		result, err := ops.ListMovies(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, result, func() {
			fmt.Printf("Total: %d\n\n", result.TotalRecordCount)
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tNAME\tYEAR\tRATING\tPATH")
			for _, item := range result.Items {
				id := item.ID
				if len(id) > 12 {
					id = id[:12]
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					id, item.Name, ops.FormatYear(item.ProductionYear),
					ops.FormatRating(item.CommunityRating), item.Path)
			}
			w.Flush()
		})
		return nil
	},
}

var itemsCmd = &cobra.Command{
	Use:   "items",
	Short: "List items in a library",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		libraryID, _ := cmd.Flags().GetString("library")
		itemType, _ := cmd.Flags().GetString("type")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := ops.ListItems(c, ops.ListItemsOptions{
			LibraryID: libraryID,
			ItemType:  itemType,
			Limit:     limit,
		})
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, result, func() {
			fmt.Printf("Total: %d\n\n", result.TotalRecordCount)
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tTYPE\tNAME\tYEAR\tPATH")
			for _, item := range result.Items {
				id := item.ID
				if len(id) > 12 {
					id = id[:12]
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					id, item.Type, item.Name, ops.FormatYear(item.ProductionYear), item.Path)
			}
			w.Flush()
		})
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for media",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		query := strings.Join(args, " ")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := ops.Search(c, query, limit)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, result, func() {
			fmt.Printf("Results: %d\n\n", result.TotalRecordCount)
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tTYPE\tNAME\tYEAR")
			for _, h := range result.Hints {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					h.ID, h.Type, h.Name, ops.FormatYear(h.ProductionYear))
			}
			w.Flush()
		})
		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info <id>",
	Short: "Show item details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		item, err := ops.GetItemInfo(c, args[0])
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, item, func() {
			fmt.Printf("Name:     %s\n", item.Name)
			fmt.Printf("Type:     %s\n", item.Type)
			fmt.Printf("ID:       %s\n", item.ID)
			if item.ProductionYear > 0 {
				fmt.Printf("Year:     %d\n", item.ProductionYear)
			}
			if item.CommunityRating > 0 {
				fmt.Printf("Rating:   %.1f\n", item.CommunityRating)
			}
			if item.OfficialRating != "" {
				fmt.Printf("Rated:    %s\n", item.OfficialRating)
			}
			if item.Runtime != "" {
				fmt.Printf("Runtime:  %s\n", item.Runtime)
			}
			fmt.Printf("Path:     %s\n", item.Path)
			if item.Overview != "" {
				fmt.Printf("Overview: %s\n", item.Overview)
			}
			if len(item.ProviderIds) > 0 {
				for k, v := range item.ProviderIds {
					fmt.Printf("%-9s %s\n", k+":", v)
				}
			}
		})
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update item metadata",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		updates := map[string]any{}
		if v, _ := cmd.Flags().GetString("name"); v != "" {
			updates["Name"] = v
		}
		if v, _ := cmd.Flags().GetInt("year"); v > 0 {
			updates["ProductionYear"] = v
		}
		if v, _ := cmd.Flags().GetString("overview"); v != "" {
			updates["Overview"] = v
		}
		if len(updates) == 0 {
			return fmt.Errorf("no updates specified (use --name, --year, --overview)")
		}
		if err := ops.UpdateItem(c, args[0], updates); err != nil {
			return err
		}
		fmt.Println("Item updated")
		return nil
	},
}

var identifyCmd = &cobra.Command{
	Use:   "identify <id>",
	Short: "Identify/match item with metadata providers",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		name, _ := cmd.Flags().GetString("name")
		year, _ := cmd.Flags().GetInt("year")
		providerIDs := map[string]string{}
		if v, _ := cmd.Flags().GetString("imdb"); v != "" {
			providerIDs["Imdb"] = v
		}
		if v, _ := cmd.Flags().GetString("tmdb"); v != "" {
			providerIDs["Tmdb"] = v
		}
		if err := ops.IdentifyItem(c, args[0], name, year, providerIDs); err != nil {
			return err
		}
		fmt.Println("Item identified and metadata updated")
		return nil
	},
}

var refreshCmd = &cobra.Command{
	Use:   "refresh <id>",
	Short: "Re-fetch metadata from providers",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		if err := ops.RefreshMetadata(c, args[0]); err != nil {
			return err
		}
		fmt.Println("Metadata refresh queued")
		return nil
	},
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan all libraries for changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		if err := ops.ScanLibrary(c); err != nil {
			return err
		}
		fmt.Println("Library scan started")
		return nil
	},
}

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Show active playback sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}
		sessions, err := ops.ListSessions(c)
		if err != nil {
			return err
		}
		cmdutil.Render(cmd, sessions, func() {
			w := cmdutil.NewTabWriter()
			fmt.Fprintln(w, "ID\tCLIENT\tDEVICE\tUSER\tNOW PLAYING")
			for _, s := range sessions {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					s.ID, s.Client, s.DeviceName, s.UserName, s.NowPlaying)
			}
			w.Flush()
		})
		return nil
	},
}
