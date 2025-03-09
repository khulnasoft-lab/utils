package dnsutil

import (
	"github.com/weppos/publicsuffix-go/publicsuffix"

	stringsutil "github.com/khulnasoft-lab/utils/strings"
)

// Split takes a domain name and decomposes it into its subdomain and domain components.
// The function returns the subdomain, the domain, and an error if the decomposition process fails.
//
// For example:
//   - Input: "http://www.example.com"
//   - Output: "www", "example.com", nil
func Split(name string) (string, string, error) {
	name = stringsutil.TrimPrefixAny(name, "http://", "https://")
	dn, err := publicsuffix.ParseFromListWithOptions(publicsuffix.DefaultList, name, publicsuffix.DefaultFindOptions)
	if err != nil {
		return "", "", err
	}

	return dn.TRD, dn.SLD + "." + dn.TLD, nil
}
