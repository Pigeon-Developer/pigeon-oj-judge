package actuator

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client

func initDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	dockerClient = cli
	return cli
}

func ImagePull(_image string) {
	reader, err := dockerClient.ImagePull(context.Background(), _image, image.PullOptions{})

	if err != nil {
		panic(err)
	}
	defer reader.Close()

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("image pull result ", bodyString)
}
