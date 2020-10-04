package logging

import "os"

func GetUnsetEnvVars(vars []string) []string {
	notSet := []string{}

	for _, v := range vars {
		_, set := os.LookupEnv(v)

		if !set {
			notSet = append(notSet, v)
		}
	}

	return notSet
}
