package testHelpers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/moby/moby/client"
)

func checkoutApp(cwd, appName string) error {
	src := path.Join(cwd, "..", "exosphere-shared", "example-apps", appName)
	dest := path.Join(cwd, "tmp", appName)
	err := os.RemoveAll(dest)
	if err != nil {
		return err
	}
	return CopyDir(src, dest)
}

func setupApp(cwd, appName string) error {
	cmdPath := path.Join(cwd, "..", "exo-setup", "bin", "exo-setup")
	cmd := exec.Command(cmdPath) // nolint gas
	cmd.Dir = path.Join(cwd, "tmp", appName)
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running setup\nOutput:\n%s\nError:%s", string(outputBytes), err)
	}
	return nil
}

func addFile(cwd, appName, serviceFolder, fileName string) error {
	filePath := path.Join(cwd, "tmp", appName, serviceFolder, fileName)
	return ioutil.WriteFile(filePath, []byte("test"), 0644)
}

func startAndRemoveContainer(dockerClient *client.Client, imageName string) error {
	ctx := context.Background()
	createdBody, err := dockerClient.ContainerCreate(
		ctx,
		&container.Config{Image: imageName},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		"",
	)
	if err != nil {
		return err
	}
	return dockerClient.ContainerRemove(ctx, createdBody.ID, types.ContainerRemoveOptions{})
}

// CleanFeatureContext defines the festure context for features/clean.feature
// nolint gocyclo
func CleanFeatureContext(s *godog.Suite) {
	var dockerClient *client.Client

	s.BeforeSuite(func() {
		var err error
		dockerClient, err = client.NewEnvClient()
		if err != nil {
			panic(err)
		}
	})

	s.Step(`^my machine has both dangling and non-dangling Docker images and volumes$`, func() error {
		appName := "external-dependency"
		serviceName := "mongo"
		imageName := "mongo:3.4.0"
		err := checkoutApp(cwd, appName)
		if err != nil {
			return fmt.Errorf("Error checking out app: %v", err)
		}
		err = setupApp(cwd, appName)
		if err != nil {
			return fmt.Errorf("Error setting up app (first time): %v", err)
		}
		err = addFile(cwd, appName, serviceName, "test.txt")
		if err != nil {
			return fmt.Errorf("Error adding file: %v", err)
		}
		err = setupApp(cwd, appName)
		if err != nil {
			return fmt.Errorf("Error setting up app (second time): %v", err)
		}
		err = startAndRemoveContainer(dockerClient, imageName)
		if err != nil {
			return fmt.Errorf("Error starting and removing container up: %v", err)
		}
		return nil
	})

	s.Step(`^it has non-dangling images$`, func() error {
		ctx := context.Background()
		filtersArgs := filters.NewArgs()
		filtersArgs.Add("dangling", "false")
		imageSummaries, err := dockerClient.ImageList(ctx, types.ImageListOptions{
			All:     false,
			Filters: filtersArgs,
		})
		if err != nil {
			return err
		}
		if len(imageSummaries) == 0 {
			return fmt.Errorf("Expected non-dangling images but there are none")
		}
		return nil
	})

	s.Step(`^it does not have dangling images$`, func() error {
		ctx := context.Background()
		filtersArgs := filters.NewArgs()
		filtersArgs.Add("dangling", "true")
		imageSummaries, err := dockerClient.ImageList(ctx, types.ImageListOptions{
			All:     false,
			Filters: filtersArgs,
		})
		if err != nil {
			return err
		}
		if len(imageSummaries) != 0 {
			return fmt.Errorf("Expected no dangling images but there are %d", len(imageSummaries))
		}
		return nil
	})

	s.Step(`^it does not have dangling volumes$`, func() error {
		ctx := context.Background()
		filtersArgs := filters.NewArgs()
		filtersArgs.Add("dangling", "true")
		volumesListOKBody, err := dockerClient.VolumeList(ctx, filtersArgs)
		if err != nil {
			return err
		}
		if len(volumesListOKBody.Volumes) != 0 {
			return fmt.Errorf("Expected no dangling volumes but there are %d", len(volumesListOKBody.Volumes))
		}
		return nil
	})
}