package module

import "fmt"

func GetRepositoryUrl(module Module) string {
	repoName := fmt.Sprintf("terraform-%s-%s", module.System, module.Name) // TODO hashi?
	return fmt.Sprintf("https://github.com/%s/%s", module.Namespace, repoName)
}
