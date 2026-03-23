package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jrogala/jellyfin-cli/internal/cmdutil"
	"github.com/jrogala/jellyfin-cli/pkg/ops"
	"github.com/spf13/cobra"
	"github.com/vishen/go-chromecast/application"
	"github.com/vishen/go-chromecast/dns"
)

func init() {
	rootCmd.AddCommand(castCmd, castDevicesCmd, castStopCmd)

	castCmd.Flags().StringP("device", "d", "", "Chromecast device name")
	castCmd.Flags().String("device-ip", "", "Chromecast device IP (skip discovery)")
	castCmd.Flags().Int("device-port", 8009, "Chromecast device port")
	castCmd.MarkFlagRequired("device")

	castStopCmd.Flags().StringP("device", "d", "", "Chromecast device name")
	castStopCmd.Flags().String("device-ip", "", "Chromecast device IP")
	castStopCmd.Flags().Int("device-port", 8009, "Chromecast device port")
	castStopCmd.MarkFlagRequired("device")
}

// discoverDevice finds a Chromecast device by name via mDNS.
func discoverDevice(name string) (string, int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := dns.DiscoverCastDNSEntries(ctx, nil)
	if err != nil {
		return "", 0, "", fmt.Errorf("discovery failed: %w", err)
	}

	for entry := range ch {
		if strings.EqualFold(entry.GetName(), name) ||
			strings.Contains(strings.ToLower(entry.GetName()), strings.ToLower(name)) {
			return entry.GetAddr(), entry.GetPort(), entry.GetName(), nil
		}
	}
	return "", 0, "", fmt.Errorf("device %q not found on network", name)
}

// collectDevices collects all discovered devices within timeout.
func collectDevices() ([]dns.CastEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := dns.DiscoverCastDNSEntries(ctx, nil)
	if err != nil {
		return nil, err
	}

	var devices []dns.CastEntry
	for entry := range ch {
		devices = append(devices, entry)
	}
	return devices, nil
}

var castDevicesCmd = &cobra.Command{
	Use:   "cast-devices",
	Short: "Scan for Chromecast devices on the network",
	RunE: func(cmd *cobra.Command, args []string) error {
		devices, err := collectDevices()
		if err != nil {
			return err
		}
		if len(devices) == 0 {
			fmt.Println("No Chromecast devices found")
			return nil
		}
		w := cmdutil.NewTabWriter()
		fmt.Fprintln(w, "NAME\tIP\tPORT\tUUID")
		for _, d := range devices {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", d.GetName(), d.GetAddr(), d.GetPort(), d.GetUUID())
		}
		w.Flush()
		return nil
	},
}

var castCmd = &cobra.Command{
	Use:   "cast <item-id or search query>",
	Short: "Cast a Jellyfin item to a Chromecast device",
	Long: `Cast a Jellyfin video to a Chromecast device (e.g. Freebox).

Examples:
  jellyfin cast 47436921e1af --device "Freebox"              Cast by item ID
  jellyfin cast "kung fu panda" --device "Freebox"           Search and cast
  jellyfin cast 47436921e1af --device-ip 192.168.1.1         Cast by IP (skip discovery)`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmdutil.NewClient()
		if err != nil {
			return err
		}

		// Resolve item: if it looks like an ID use it, otherwise search
		query := strings.Join(args, " ")
		itemID := query
		itemName := ""

		if len(query) < 20 || strings.Contains(query, " ") {
			result, err := ops.Search(c, query, 1)
			if err != nil {
				return err
			}
			if len(result.Hints) == 0 {
				return fmt.Errorf("no results for %q", query)
			}
			itemID = result.Hints[0].ID
			itemName = result.Hints[0].Name
		}

		if itemName == "" {
			item, err := ops.GetItemInfo(c, itemID)
			if err != nil {
				return err
			}
			itemName = item.Name
		}

		streamURL := c.BaseURL() + "/Videos/" + itemID + "/stream?static=true&api_key=" + "token"

		// Resolve Chromecast device
		deviceName, _ := cmd.Flags().GetString("device")
		deviceIP, _ := cmd.Flags().GetString("device-ip")
		devicePort, _ := cmd.Flags().GetInt("device-port")

		if deviceIP == "" {
			ip, port, name, err := discoverDevice(deviceName)
			if err != nil {
				return err
			}
			deviceIP = ip
			devicePort = port
			deviceName = name
		}

		// Connect and cast
		app := application.NewApplication(
			application.WithServerPort(devicePort),
			application.WithDebug(false),
		)

		if err := app.Start(deviceIP, devicePort); err != nil {
			return fmt.Errorf("connecting to %s: %w", deviceName, err)
		}

		if err := app.Load(streamURL, 0, "video/mp4", false, false, false); err != nil {
			return fmt.Errorf("casting to %s: %w", deviceName, err)
		}

		fmt.Printf("Casting %q to %s\n", itemName, deviceName)
		return nil
	},
}

var castStopCmd = &cobra.Command{
	Use:   "cast-stop",
	Short: "Stop casting on a Chromecast device",
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceName, _ := cmd.Flags().GetString("device")
		deviceIP, _ := cmd.Flags().GetString("device-ip")
		devicePort, _ := cmd.Flags().GetInt("device-port")

		if deviceIP == "" {
			ip, port, name, err := discoverDevice(deviceName)
			if err != nil {
				return err
			}
			deviceIP = ip
			devicePort = port
			deviceName = name
		}

		app := application.NewApplication(
			application.WithServerPort(devicePort),
			application.WithDebug(false),
		)
		if err := app.Start(deviceIP, devicePort); err != nil {
			return err
		}
		app.Stop()
		fmt.Printf("Stopped casting on %s\n", deviceName)
		return nil
	},
}
