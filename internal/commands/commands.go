package commands

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

//func shellCmd(command string, arg ...string) (string, error) {
func RunCmd(env string, command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)

	cmd.Env = append(os.Environ(),
		//"RSYNC_PASSWORD="+rsyncPassword,
		env,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return "", errors.Wrap(err, stderr.String())
	}

	return "Result: " + out.String(), nil
}
