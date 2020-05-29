package cmd

import (
	"fmt"
	"log"
	"os"

	whois "github.com/luizm/aws-whois/pkg"

	"github.com/urfave/cli/v2"
)

func Execute() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "region",
				Aliases:  []string{"r"},
				Usage:    "the region to use. Overrides config/env settings.",
				Value:    "us-east-1",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "profile",
				Aliases:  []string{"p"},
				Usage:    "use a specific profile from your credential file.",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "ignore-profile",
				Aliases:  []string{"i"},
				Usage:    "ignore a specific profile from your credential file, can be used multiple times.",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"a"},
				Usage:    "the ip or dns address to find the resource associated",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			ip := c.String("address")
			if !whois.IsIP(c.String("address")) {
				ips, err := whois.ResolvDNS(c.String("address"))
				if err != nil {
					return fmt.Errorf(`unable to resolve DNS: %w`, err)
				}
				ip = ips[0]
			}

			if c.String("profile") == "" {
				profiles, err := whois.ShowAWSProfile()
				if err != nil {
					return fmt.Errorf(`failed to load the local profiles: %w`, err)
				}
				profiles = whois.DiffSliceString(profiles, c.StringSlice("ignore-profile"))

				for _, p := range profiles {
					result, err := whois.FindIP(p, c.String("region"), ip)
					if err != nil {
						log.Println(fmt.Errorf(`failure to search in profile %v: %w`, p, err))
						continue
					}
					fmt.Println(whois.ToJson(result))
				}
			} else {
				p := c.String("profile")
				result, err := whois.FindIP(p, c.String("region"), ip)
				if err != nil {
					log.Println(fmt.Errorf(`failure to search in profile %v: %w`, p, err))
				}
				fmt.Println(whois.ToJson(result))
			}
			return nil
		},
	}

	app.Name = "aws-whois"
	app.Usage = "Find out where and who is a specific address"
	app.Version = "v1.0.2"

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
