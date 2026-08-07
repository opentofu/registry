package verification
 
import (
	"fmt"
 
	"github.com/opentofu/registry-stable/internal/github"
)
 
// VerifyGithubUser returns a Step confirming that username is a public
// member of orgName on GitHub. This is shared by any submission workflow
// that needs to confirm the requester is affiliated with a provider's
// organization (e.g. GPG key submission, re-index requests).
func VerifyGithubUser(client github.Client, username string, orgName string) *Step {
	step := &Step{Name: "Validate Github user"}
 
	s := step.RunStep(fmt.Sprintf("User is a member of the organization %s", orgName), func() error {
		member, err := client.IsUserInOrganization(username, orgName)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if !member {
			return fmt.Errorf("user is not a member of the organization")
		}
		return nil
	})
	s.Remarks = []string{"If this is incorrect, please ensure that your organization membership is public. For more information, see [Github Docs - Publicizing or hiding organization membership](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-your-membership-in-organizations/publicizing-or-hiding-organization-membership)"}
 
	return step
}
 
