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
	var profile []string
	var p []string

	config, err := ini.Load(filepath.Join(os.Getenv("HOME"), ".aws/config"))
	if err != nil {
		return nil, err
	}

	credentials, err := ini.Load(filepath.Join(os.Getenv("HOME"), ".aws/credentials"))
	if err != nil {
		return nil, err
	}
	profile = append(config.SectionStrings(), credentials.SectionStrings()...)
	for _, r := range profile {
		if r != "DEFAULT" {
			p = append(p, strings.Replace(r, "profile ", "", -1))
		}
	}
	return p, nil
}

func DiffSliceString(a, b []string) []string {
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
