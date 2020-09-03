package postgres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type Container struct {
	ID   string
	Host string
}

func LaunchContainer(log *log.Logger) *Container {
	cmd := exec.Command("docker", "run", "-P", "-d", "-e", "POSTGRES_PASSWORD=postgres", "postgres:latest")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatalf("[error] Starting container: %v", err)
	}

	ContainerID := out.String()[:12]

	cmd = exec.Command("docker", "inspect", ContainerID)

	out.Reset()

	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatalf("[error] Inspecting container %s: %v", ContainerID, err)
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP5432 []struct {
					HostIP   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"5432/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		log.Fatalf("[error] Failed to decode json: %v", err)
	}

	network := doc[0].NetworkSettings.Ports.TCP5432[0]

	container := Container{
		ID:   ContainerID,
		Host: network.HostIP + ":" + network.HostPort,
	}

	fmt.Println("Container started:", container.ID)

	return &container
}

func DestroyContainer(container *Container) {
	if err := exec.Command("docker", "stop", container.ID).Run(); err != nil {
		log.Fatalf("[error] Failed to stop container: %v", err)
	}
	fmt.Println("Container stopped:", container.ID)

	if err := exec.Command("docker", "rm", container.ID, "-v").Run(); err != nil {
		log.Fatalf("[error] Failed to remove container: %v", err)
	}
	fmt.Println("Container removed:", container.ID)
}
