package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/storage"
	"github.com/urfave/cli/v2"
)

func SetupCLI() {

	settings := &models.ScanSettings{}

	settings.Store = storage.OpenStore()
	defer settings.Store.Close()

	sharingan := cli.NewApp()
	sharingan.Name = "Sharingan"
	sharingan.Usage = "Wrapper and analyzer for offensive security recon tools"
	sharingan.Version = "0.2.0"

	sharingan.Flags = GetGlobalFlags(settings)

	// TODO: move commands to their own file
	sharingan.Commands = []*cli.Command{
		{
			Name:  "dns",
			Usage: "Perform a DNS scan : sharingancli --target ProgramName dns --rootdomain rootdomain.com --dns-wordlist ./path/to/list",
			Action: func(c *cli.Context) error {
				fmt.Println("Must use bruteforce or addsubs command")
				return nil
			},
			Subcommands: []*cli.Command{
				{
					Name:  "bruteforce",
					Usage: "Brute force subdomains for a given rootdomain",
					Flags: GetDNSFlags(settings),
					Action: func(c *cli.Context) error {
						RunDNSRecon(settings)
						return nil
					},
				},
				{
					Name:  "addsubs",
					Usage: "Add subdomains to a given program using stdin ie : cat subs | sharingancli --target example DNS addsubs",
					Flags: GetDNSFlags(settings),
					Action: func(c *cli.Context) error {
						AddSubsFromInput(settings)
						return nil
					},
				},
			},
		},
		{
			Name:  "scan",
			Usage: "Perform a service scan using nmap -sV : sharingancli --target ProgramName scan",
			Flags: GetNMapFlags(settings),
			Subcommands: []*cli.Command{
				{
					Name:  "interactive",
					Usage: "Manually select hosts from program to scan",
					Action: func(c *cli.Context) error {
						RunNmapScanInteractive(settings)
						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				RunNmapScan(settings)
				return nil
			},
		},
		{
			Name:  "info",
			Usage: "Gets information contained in a program/scan",
			Subcommands: []*cli.Command{
				{
					Name:  "domains",
					Usage: "Prints subdomains for program/scan on stdout",
					Action: func(c *cli.Context) error {
						PrintDomains(settings)
						return nil
					},
				},
			},
		},
	}

	sharingan.Action = func(c *cli.Context) error {
		return nil
	}

	err := sharingan.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
