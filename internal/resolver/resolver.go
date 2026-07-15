package resolver

import "strings"

type Rule struct {
	Prefix  string
	Flatten bool
}

func Resolve(image string, rule Rule) string {

	name := image

	if rule.Flatten {
		parts := strings.Split(image, "/")
		name = parts[len(parts)-1]
	}

	if rule.Prefix == "" {
		return name
	}

	return strings.TrimSuffix(rule.Prefix, "/") + "/" + name
}
