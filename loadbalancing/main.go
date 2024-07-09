package main

import (
	"fmt"
	"log"
	"sysdesign/loadbalancing/container"
	"time"
)

func main() {
	ncm, err := container.NewNgixContainerManager()
	if err != nil {
		log.Fatalf("Failed to create Nginx Container Manager: %v", err)
	}

	fmt.Println("Creating a new Nginx container...")
	containerInfo, err := ncm.CreateContainer()
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	fmt.Printf("Container created successfully. ID: %s, URL: %s\n", containerInfo.ID, containerInfo.URL)

	// Keep the container running for a while
	time.Sleep(10 * time.Second)

	fmt.Println("Removing the container...")
	if err := ncm.RemoveContainer(containerInfo.ID); err != nil {
		log.Fatalf("Failed to remove container: %v", err)
	}

	fmt.Println("Container removed successfully.")
}
