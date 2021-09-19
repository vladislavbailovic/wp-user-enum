package cli

import (
	"flag"
	"wp-enum/pkg/data"
)

func GetFlags() data.Constraints {
	defaults := data.DefaultConstraints()

	url := flag.String("url", defaults.URL, "WordPress URL")
	kind := flag.Int("enum", defaults.Kind, "Enumeration type")
	start := flag.Int("start", defaults.Start, "Start enumeration at this user ID")
	end := flag.Int("end", defaults.End, "End enumeration with this user ID")
	random := flag.Bool("ua", defaults.RandomUA, "Randomize User-Agent")
	flag.Parse()

	defaults.URL = *url
	defaults.Kind = *kind
	defaults.Start = *start
	defaults.End = *end
	defaults.RandomUA = *random

	return defaults
}
