package x

import "strings"

func ParseStructTagOptions(opts []string) map[string]string {
	options := make(map[string]string, len(opts))

	for _, o := range opts {
		if strings.ContainsRune(o, '=') {
			parts := strings.SplitN(o, "=", 2)
			options[parts[0]] = parts[1]
		} else {
			options[o] = ""
		}
	}

	return options
}
