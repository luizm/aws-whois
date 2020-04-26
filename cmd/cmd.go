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
				Name:        "region",
				Aliases:     []string{"r"},
				Usage:       "The region to use. Overrides config/env settings",
				DefaultText: "us-east-1",
				Required:    false,
			},
			&cli.StringSliceFlag{
				Name:     "profile",
				Aliases:  []string{"p"},
				Usage:    "Use a specific profile from your credential file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "ip",
				Aliases:  []string{"i"},
				Usage:    "The ip address to find the resource associated",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "dns",
				Aliases:  []string{"d"},
				Usage:    "The dns address to find the resource associated, if return more than 1 ip, will be used the first",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			var ip string
			ip = c.String("ip")

			if c.String("dns") != "" {
				ips, _ := whois.ResolvDNS(c.String("dns"))
				ip = ips[0]
			}
			for _, p := range c.StringSlice("profile") {
				result, err := whois.FindIP(p, c.String("region"), ip)
				if err != nil {
					return err
				}
				fmt.Println(string(result))
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
