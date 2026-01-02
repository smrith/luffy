package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"luffy/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "luffy [query]",
	Short: "Watch movies and TV shows from the commandline",
	Args:  cobra.ArbitraryArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		client := core.NewClient()
		ctx := &core.Context{
			Client: client,
		}

		if len(args) == 0 {
			ctx.Query = core.Prompt("Search")
		} else {
			ctx.Query = strings.Join(args, " ")
		}

		results, err := core.SearchContent(ctx.Query, ctx.Client)
		if err != nil {
			return err
		}

		var titles []string
		for _, r := range results {
			titles = append(titles, fmt.Sprintf("[%s] %s", r.Type, r.Title))
		}

		idx := core.Select("Results:", titles)
		selected := results[idx]

		ctx.Title = selected.Title
		ctx.URL = selected.URL
		ctx.ContentType = selected.Type

		fmt.Println("Selected:", ctx.Title)

		mediaID, err := core.GetMediaID(ctx.URL, ctx.Client)
		if err != nil {
			return err
		}

		var episodeID string

		if ctx.ContentType == core.Series {
			seasons, err := core.GetSeasons(mediaID, ctx.Client)
			if err != nil {
				return err
			}
			if len(seasons) == 0 {
				return fmt.Errorf("no seasons found")
			}

			var sNames []string
			for _, s := range seasons {
				sNames = append(sNames, s.Name)
			}
			sIdx := core.Select("Seasons:", sNames)
			selectedSeason := seasons[sIdx]

			episodes, err := core.GetEpisodes(selectedSeason.ID, true, ctx.Client)
			if err != nil {
				return err
			}
			if len(episodes) == 0 {
				return fmt.Errorf("no episodes found")
			}

			var eNames []string
			for _, e := range episodes {
				eNames = append(eNames, e.Name)
			}
			eIdx := core.Select("Episodes:", eNames)
			episodeID = episodes[eIdx].ID

		} else {
			episodes, err := core.GetEpisodes(mediaID, false, ctx.Client)
			if err != nil || len(episodes) == 0 {
				return fmt.Errorf("could not find movie info")
			}

			if len(episodes) == 1 {
				episodeID = episodes[0].ID
			} else {
				var eNames []string
				for _, e := range episodes {
					eNames = append(eNames, e.Name)
				}
				eIdx := core.Select("Episodes/Parts:", eNames)
				episodeID = episodes[eIdx].ID
			}
		}

		servers, err := core.GetServers(episodeID, ctx.Client)
		if err != nil {
			return err
		}
		if len(servers) == 0 {
			return fmt.Errorf("no servers found")
		}

		var srvNames []string
		for _, s := range servers {
			srvNames = append(srvNames, s.Name)
		}
		srvIdx := core.Select("Servers:", srvNames)
		selectedServer := servers[srvIdx]

		link, err := core.GetLink(selectedServer.ID, ctx.Client)
		if err != nil {
			return err
		}

		actions := []string{"Play", "Copy Link"}
		actIdx := core.Select("Action:", actions)

		if actions[actIdx] == "Play" {
			fmt.Println("Starting mpv with:", link)
			cmd := exec.Command("mpv", link)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error playing:", err)
			}
		} else {
			fmt.Println("Link:", link)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

