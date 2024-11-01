package jankservers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var containerStates = make(map[string]string)

func init() {
	containerStates["/ark_server"] = ""
	containerStates["/valheim-server"] = ""
}

func StartServerMonitor() {
	go getContainersState()
}

func getContainersState() {
	opts, err := client.NewClientWithOpts(
		client.WithHost("unix:///var/run/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		fmt.Println("Error creating Docker client:", err)
		return
	}

	for {
		containers, err := opts.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			fmt.Println("Error listing containers:", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		for _, ctr := range containers {
			if len(ctr.Names) > 0 { // Ensure there is at least one name
				containerName := ctr.Names[0] // Access the first name
				if state, ok := containerStates[containerName]; ok {
					log.Println(containerName, ctr.State)
					// Update state if it has changed
					if state != ctr.State {
						fmt.Printf("Container %s state changed: %s\n", containerName, ctr.State)
						containerStates[containerName] = ctr.State
					}
				}
			}
		}
		time.Sleep(5 * time.Second) // Sleep to avoid busy waiting
	}
}
