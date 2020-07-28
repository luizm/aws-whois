package whois

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/pretty"

	"gopkg.in/ini.v1"
)

func ShowAWSProfile() ([]string, error) {
	var profiles []string

	config, err := ini.Load(filepath.Join(os.Getenv("HOME"), ".aws/config"))
	if err != nil {
		return nil, err
	}

	for _, c := range config.SectionStrings() {
		if c != "DEFAULT" {
			profiles = append(profiles, strings.Replace(c, "profile ", "", -1))
		}
	}
	return profiles, nil
}

func RemoveElementOfSliceString(a, b []string) []string {
	mb := make(map[string]string, len(b))
	for _, x := range b {
		mb[x] = x
	}
	var diff []string
	for _, v := range a {
		if _, found := mb[v]; !found {
			diff = append(diff, v)
		}
	}
	return diff
}

func ToJson(r *Result) string {
	json, _ := json.MarshalIndent(r, "", "   ")
	b := pretty.Color(json, nil)
	return string(b)
}
