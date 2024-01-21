package storagegateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type InstanceCfg struct {
	AccessKey string
	SecretKey string
	Endpoint  string
}

func DiscoverMinioInstancesInDocker(ctx context.Context) []InstanceCfg {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		panic(err.Error())
	}
	dockerClient.NegotiateAPIVersion(ctx)

	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("ancestor", "minio/minio")),
	})
	if err != nil {
		panic(err.Error())
	}

	instances := []InstanceCfg{}
	for _, ctr := range containers {
		instances = append(instances, parseInstanceCfg(ctx, dockerClient, ctr))
	}
	return instances
}

func parseEnvs(envs []string, key string) (string, error) {
	for _, env := range envs {
		splits := strings.Split(env, "=")
		if len(splits) != 2 {
			return "", fmt.Errorf("invalid env: %s", env)
		}
		if splits[0] == key {
			return splits[1], nil
		}
	}
	return "", fmt.Errorf("key not found: %s", key)
}

func parseInstanceCfg(ctx context.Context, dockerClient *client.Client, ctr types.Container) InstanceCfg {
	var ip string
	for _, network := range ctr.NetworkSettings.Networks {
		ip = network.IPAddress
		break
	}
	endpoint := fmt.Sprintf("%s:9000", ip)

	containerInfo, err := dockerClient.ContainerInspect(ctx, ctr.ID)
	if err != nil {
		panic(err.Error())
	}
	envs := containerInfo.Config.Env
	accessKey, err := parseEnvs(envs, "MINIO_ACCESS_KEY")
	if err != nil {
		panic(err.Error())
	}
	secretKey, err := parseEnvs(envs, "MINIO_SECRET_KEY")
	if err != nil {
		panic(err.Error())
	}
	return InstanceCfg{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Endpoint:  endpoint,
	}
}
