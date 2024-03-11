package shared

import (
	"os"
	"os/exec"
	"strings"

	"github.com/zwoo-hq/zwooc/pkg/model"
)

func CreateBaseCommand(executable string, c model.ProfileWrapper, extraArgs []string) (*exec.Cmd, []string) {
	cmd := exec.Command(executable)
	cmd.Env = os.Environ()
	cmd.Dir = c.GetDirectory()

	profileOptions := c.GetProfileOptions()
	cmd.Env = append(cmd.Env, profileOptions.Env...)

	additionalArgs := []string{}
	for k, v := range profileOptions.Args {
		if strings.HasPrefix(k, "-") {
			additionalArgs = append(additionalArgs, k, v)
		} else {
			additionalArgs = append(additionalArgs, "--"+k, v)
		}
	}

	return cmd, append(additionalArgs, extraArgs...)
}
