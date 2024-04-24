package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"io"
	"strings"
	"time"
)

var cli *client.Client

// DockerImage
// @Description: 镜像结构体
type DockerImage struct {
	ID       string   `json:"ID"`
	RepoTags []string `json:"RepoTags"`
	Status   string   `json:"Status"`
	Created  int64    `json:"Created"`
	Size     int64    `json:"Size"`
}

type DockerContainerConfig struct {
	Image        string                `json:"Image"`
	Cmd          string                `json:"Cmd"`
	Env          []string              `json:"Env"`
	ExposedPorts map[string]struct{}   `json:"ExposedPorts"`
	HostConfig   *container.HostConfig `json:"HostConfig"`
}

// InitDocker
//
//	@Description: 初始化docker客户端
func InitDocker() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic("Init failed because " + err.Error())
	}
}

// GetDockerClient
//
//	@Description: 获取docker客户端
//	@return *client.Client
func GetDockerClient() *client.Client {
	return cli
}

// pullImage
//
//	@Description: 拉取镜像
//	@param imageName
//	@return string
//	@return error
func pullImage(cli *client.Client, imageName string) (string, error) {
	ctx := context.Background()

	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return "", err
	}
	defer out.Close()

	// * 读取下载的镜像信息
	var buf bytes.Buffer
	_, err = io.Copy(&buf, out)
	if err != nil && err != io.EOF {
		return "", err
	}

	// * 转换为json格式
	var lastLine string
	for _, line := range strings.Split(buf.String(), "n") {
		if line != "" {
			lastLine = line
		}
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(lastLine), &result)
	if err != nil {
		return "", err
	}

	// * 检测末尾返回的结果是否为下载成功
	if status, ok := result["status"]; ok && strings.HasPrefix(status.(string), "Downloaded newer image for") || strings.HasPrefix(status.(string), "Image is up to date for") {
		//	* 如果镜像已经存在，则直接返回
		// * 查看下载的镜像信息
		imageNew, _, err := cli.ImageInspectWithRaw(ctx, imageName)
		if err != nil {
			return "", err
		}
		return imageNew.ID, nil
	}

	return "", fmt.Errorf("image 拉取失败: %s", lastLine)
}

// GetImagesList
//
//	@Description: 获取镜像列表
//	@return []string
//	@return error
func GetImagesList(cli *client.Client) ([]DockerImage, error) {
	ctx := context.Background()
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return nil, err
	}

	//var imagesList []string
	imagesList := make([]DockerImage, 0, len(images))
	for _, imageItem := range images {
		status := "Unused"
		//  * 如果镜像被容器使用，则状态为In use
		if imageItem.Containers > 0 {
			status = "In use"
		}
		imagesList = append(imagesList, DockerImage{
			ID:       imageItem.ID,
			RepoTags: imageItem.RepoTags,
			Status:   status,
			Created:  imageItem.Created,
			Size:     imageItem.Size,
		})
	}
	return imagesList, nil
}

// execContainerCmd
// 注意这里如果想执行什么命令，建议使用"/bin/sh -c "exec""否则可能会出问题
// e.g. out, err := execContainerCmd(resp.ID, []string{"/bin/sh", "-c", "echo flag{THis_is_a_test_flag} > /flag"})
// ! 不要读取再执行flag写入，而是传入环境变量flag，然后在内部传入再清除。
//
//	@Description: 执行容器命令
//	@param containerId
//	@param cmd
//	@return string
//	@return error
func execContainerCmd(cli *client.Client, containerId string, cmd []string) (string, error) {
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}
	response, err := cli.ContainerExecCreate(context.Background(), containerId, execConfig)
	if err != nil {
		return "", err
	}
	// * 启动exec
	hijacked, err := cli.ContainerExecAttach(context.Background(), response.ID, types.ExecStartCheck{})
	if err != nil {
		return "", err
	}
	defer hijacked.Close()
	// * 读取exec的输出
	stdout, err := io.ReadAll(hijacked.Reader)
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

// GetContainerInfo
//
//	@Description: 获取容器信息
//	@param containerId
//	@return types.ContainerJSON
//	@return error
func GetContainerInfo(cli *client.Client, containerId string) (types.ContainerJSON, error) {
	return cli.ContainerInspect(context.Background(), containerId)
}

// DeleteContainer
//
//	@Description: 删除容器
//	@param containerId
//	@return error
func DeleteContainer(cli *client.Client, containerId string) error {
	return cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
}

// GetContainerList
//
//	@Description: 获取容器列表
//	@return []types.Container
//	@return error
func GetContainerList(cli *client.Client) ([]types.Container, error) {
	return cli.ContainerList(context.Background(), container.ListOptions{})
}

// GetContainerStartTime
//
//	@Description: 获取容器启动时间
//	@param containerID
//	@return time.Duration
//	@return error
func GetContainerStartTime(cli *client.Client, containerID string) (time.Duration, error) {
	containerInfo, err := GetContainerInfo(cli, containerID)
	if err != nil {
		return 0, err
	}
	startedAt, err := time.Parse(time.RFC3339Nano, containerInfo.State.StartedAt)
	if err != nil {
		return 0, err
	}
	return time.Since(startedAt), nil
}

func StopContainer(cli *client.Client, containerId string) error {
	return cli.ContainerStop(context.Background(), containerId, container.StopOptions{})
}
