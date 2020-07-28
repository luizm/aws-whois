package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"

	whois "github.com/luizm/aws-whois/pkg"

	"github.com/urfave/cli/v2"
)

func Execute() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "region",
				Aliases:  []string{"r"},
				Usage:    "The region to use. Overrides config/env settings.",
				Value:    "us-east-1",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "profile",
				Aliases:  []string{"p"},
				Usage:    "Use a specific profile from your credential file. Could be used multiple times.",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "ignore-profile",
				Aliases:  []string{"i"},
				Usage:    "Ignore a specifics profile from your credential file. Could be used multiple times.",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "The ip or dns address to find the resource associated.",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			var waitGroup sync.WaitGroup
			var err error

			ip := c.String("address")
			region := c.String("region")

			if !whois.IsIP(ip) {
				ips, err := whois.ResolvDNS(ip)
				if err != nil {
					return fmt.Errorf(`unable to resolve DNS: %w`, err)
				}
				ip = ips[0]
			}
			profiles := c.StringSlice("profile")

			if profiles == nil {
				profiles, err = whois.ShowAWSProfile()
				if err != nil {
					return fmt.Errorf(`failed to load the local profiles: %w`, err)
				}
			}
			profiles = whois.RemoveElementOfSliceString(profiles, c.StringSlice("ignore-profile"))

			waitGroup.Add(len(profiles))

			for _, profile := range profiles {
				go func(p string) {
					result, err := whois.FindIP(p, region, ip)
					if err != nil {
						log.Println(fmt.Errorf(`failure to search in profile %v: %w`, p, err))
						waitGroup.Done()
						return
					}

					fmt.Println(whois.ToJson(result))
					waitGroup.Done()
				}(profile)
			}
			waitGroup.Wait()
			return nil
		},
	}
	app.EnableBashCompletion = true
	app.Name = "aws-whois"
	app.Usage = "Find out where and who is a specific address"
	app.Version = "v1.0.2"

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
