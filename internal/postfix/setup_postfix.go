// Package postfix contains the Postfix configuration.
package postfix

import (
	"os"
	"os/exec"
)

func RestartPostfix() error {
	cmd := exec.Command("sudo", "systemctl", "restart", "postfix")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func RunNewAliases() error {
	cmd := exec.Command("sudo", "newaliases")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func MakeForwardScriptExecutable(path string) error {
	cmd := exec.Command("sudo", "chmod", "+x", path)
	return cmd.Run()
}

func SetupPostfix() error {
	forwardScriptPath := "/usr/local/bin/forward-to-webhook.sh"
	_, err := exec.LookPath("postfix")
	if err != nil {
		return err
	}
	_, err = exec.LookPath("procmail")
	if err != nil {
		return err
	}

	// Remove existing files

	files := []string{
		"/etc/postfix/main.cf",
		"/etc/postfix/virtual_regexp",
		"/etc/aliases",
		forwardScriptPath,
	}

	for _, f := range files {
		if err = os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	cfg := DefaultConfig()
	GenerateMainCF(cfg, "/etc/postfix/main.cf")
	GenerateAliasesFile("/etc/aliases")

	err = RestartPostfix()
	if err != nil {
		return err
	}
	err = RunNewAliases()
	if err != nil {
		return err
	}

	err = MakeForwardScriptExecutable(forwardScriptPath)
	if err != nil {
		return err
	}

	return nil
}
