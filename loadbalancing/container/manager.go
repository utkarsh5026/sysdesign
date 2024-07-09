// Package container provides functionality for managing Nginx containers using Docker.
//
// This package offers a simplified interface for creating, managing, and removing
// Nginx containers, abstracting away much of the complexity involved in directly
// interacting with the Docker API.
package container

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// NgixContainerManager handles the creation and management of Nginx containers.
type NgixContainerManager struct {
	docker *client.Client
}

// ContainerInfo holds essential information about a created container.
type ContainerInfo struct {
	ID   string // The Docker container ID
	Port int    // The host port mapped to the container's port 80
	URL  string // The URL to access the Nginx server
}

// NewNgixContainerManager creates and returns a new NgixContainerManager instance.
// It initializes a Docker client using the system's Docker environment.
//
// Returns:
//   - *NgixContainerManager: A pointer to the new NgixContainerManager instance
//   - error: An error if the Docker client creation fails
func NewNgixContainerManager() (*NgixContainerManager, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &NgixContainerManager{docker: dockerClient}, nil
}

// CreateContainer creates and starts a new Nginx container.
// It handles the entire process from ensuring the image exists to starting the container
// and returning its details.
//
// Returns:
//   - *ContainerInfo: Information about the created container
//   - error: An error if any step in the container creation process fails
func (ncm *NgixContainerManager) CreateContainer() (*ContainerInfo, error) {
	fmt.Println("Starting container creation process...")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Creating Nginx container"
	s.Start()
	defer s.Stop()

	ctx := context.Background()

	imageName := "nginx:alpine"
	if err := ncm.ensureImageExists(ctx, imageName); err != nil {
		return nil, err
	}

	config, hostConfig := ncm.prepareContainerConfig(imageName)

	containerID, err := ncm.createAndStartContainer(ctx, config, hostConfig)
	if err != nil {
		return nil, err
	}

	return ncm.inspectAndGetContainerInfo(ctx, containerID)
}

// RemoveContainer stops and removes a container with the given ID.
//
// Parameters:
//   - id: The ID of the container to remove
//
// Returns:
//   - error: An error if stopping or removing the container fails
func (cm *NgixContainerManager) RemoveContainer(id string) error {
	if err := cm.docker.ContainerStop(context.Background(), id, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}

	if err := cm.docker.ContainerRemove(context.Background(), id, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove container: %v", err)
	}

	return nil
}

// ensureImageExists checks if the specified Docker image exists locally,
// and if not, pulls it from Docker Hub.
//
// Parameters:
//   - ctx: The context for the Docker API calls
//   - imageName: The name of the Docker image to check/pull
//
// Returns:
//   - error: An error if checking for the image or pulling it fails
func (ncm *NgixContainerManager) ensureImageExists(ctx context.Context, imageName string) error {
	fmt.Printf("Checking if image %s exists locally...\n", imageName)
	_, _, err := ncm.docker.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		fmt.Printf("Image %s found locally.\n", imageName)
		return nil // Image exists
	}

	fmt.Printf("Image %s not found locally. Pulling from Docker Hub...\n", imageName)
	reader, err := ncm.docker.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %v", err)
	}
	defer reader.Close()

	// Print pull progress
	decoder := json.NewDecoder(reader)
	for {
		var message map[string]interface{}
		if err := decoder.Decode(&message); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode pull message: %v", err)
		}
		if status, ok := message["status"]; ok {
			fmt.Println(status)
		}
	}

	fmt.Printf("Image %s pulled successfully.\n", imageName)
	return nil
}

// prepareContainerConfig creates and returns the configuration for the Nginx container.
//
// Parameters:
//   - imageName: The name of the Docker image to use for the container
//
// Returns:
//   - *container.Config: The container configuration
//   - *container.HostConfig: The host configuration for the container
func (ncm *NgixContainerManager) prepareContainerConfig(imageName string) (*container.Config, *container.HostConfig) {
	fmt.Println("Preparing container configuration...")
	port := "80/tcp"
	config := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			nat.Port(port): struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(port): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "0", // Let Docker assign a port
				},
			},
		},
	}
	return config, hostConfig
}

// createAndStartContainer creates and starts a new Docker container with the given configuration.
//
// Parameters:
//   - ctx: The context for the Docker API calls
//   - config: The container configuration
//   - hostConfig: The host configuration for the container
//
// Returns:
//   - string: The ID of the created container
//   - error: An error if creating or starting the container fails
func (ncm *NgixContainerManager) createAndStartContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig) (string, error) {
	fmt.Println("Creating container...")
	resp, err := ncm.docker.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %v", err)
	}

	fmt.Printf("Container created with ID: %s\n", resp.ID)
	fmt.Println("Starting container...")
	if err := ncm.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %v", err)
	}

	return resp.ID, nil
}

// inspectAndGetContainerInfo retrieves detailed information about a container and formats it.
//
// Parameters:
//   - ctx: The context for the Docker API calls
//   - containerID: The ID of the container to inspect
//
// Returns:
//   - *ContainerInfo: Formatted information about the container
//   - error: An error if inspecting the container fails
func (ncm *NgixContainerManager) inspectAndGetContainerInfo(ctx context.Context, containerID string) (*ContainerInfo, error) {
	fmt.Println("Inspecting container...")
	containerInfo, err := ncm.docker.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %v", err)
	}

	port := "80/tcp"
	containerPort, _ := strconv.Atoi(containerInfo.NetworkSettings.Ports[nat.Port(port)][0].HostPort)
	url := fmt.Sprintf("http://localhost:%d", containerPort)

	fmt.Printf("Container is ready! Accessible at: %s\n", url)

	return &ContainerInfo{
		ID:   containerID,
		Port: containerPort,
		URL:  url,
	}, nil
}
