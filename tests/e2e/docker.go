package e2e

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Docker wraps docker CLI interaction.
type Docker struct{}

// Exec runs a command in the given container.
func (d *Docker) Exec(container string, cmd []string) (string, error) {
	args := append([]string{"exec", container}, cmd...)
	c := exec.Command("docker", args...)
	var out, stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		return "", fmt.Errorf("docker exec failed: %v\nstderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

// Logs retrieves container logs.
func (d *Docker) Logs(container string) (string, error) {
	c := exec.Command("docker", "logs", container)
	var out, stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr

	if err := c.Run(); err != nil {
		return "", fmt.Errorf("docker logs failed: %v\nstderr: %s", err, stderr.String())
	}

	return out.String(), nil
}
