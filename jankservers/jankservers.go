package jankservers

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type server struct {
	state string
	id    string
}

var containerInfo = make(map[string]server)

// type manager struct {
// 	contInfo map

// }

func init() {
	containerInfo["ark_server"] = server{state: "", id: ""}
	containerInfo["valheim-server"] = server{state: "", id: ""}
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
		log.Println("Error creating Docker client:", err)
		return
	}

	for {
		containers, err := opts.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			log.Println("Error listing containers:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, ctr := range containers {
			if len(ctr.Names) > 0 {
				containerName := ctr.Names[0]
				containerName = strings.TrimPrefix(containerName, "/")
				if serverData, ok := containerInfo[containerName]; ok {
					log.Println(containerName, ctr.State)
					if serverData.state != ctr.State || serverData.id != ctr.ID {
						log.Printf("Container %s state changed from %s to %s; ID from %s to %s\n", containerName, serverData.state, ctr.State, serverData.id, ctr.ID)
						// Update state and id
						serverData.state = ctr.State
						serverData.id = ctr.ID
						containerInfo[containerName] = serverData
					}
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}
