package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

var (
	cli *client.Client
)

func main() {
	//pullImageX("ghcr.io/gztimewalker/gzctf/test:latest")
	InitDocker()
	images, err := GetImagesList()
	for _, imageItem := range images {
		log.Println(imageItem)
	}
	if err != nil {
		return
	}
	config := &container.Config{
		//Image: "ghcr.io/gztimewalker/gzctf/test:latest",
		// * 这里的Image是镜像的ID或者名字
		Image: "arcctf-test",
		// * 这里的Cmd是容器启动时执行的命令，但
		//Cmd:   []string{"echo", "flag{THis_is_a_test_flag} > /flag"},
		Env: []string{
			"ARC_FLAG=ArcCTF{This_is_a_test_flag_4_ArcCTF}",
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"80/tcp": []nat.PortBinding{
				{
					// * 容器端口映射到宿主机端口，里面是宿主端口，外面是容器端口
					HostIP:   "0.0.0.0",
					HostPort: "0",
				},
			},
		},
		Resources: container.Resources{
			// * 限制容器的内存使用, 限制为100M
			CPUPeriod: 100000,
			// * 限制容器的CPU使用，限制为50%
			CPUQuota: 50000,
			// * 如果不设置内存限制，可以设置内存swap限制
			// * 可以通过CPUsetCPUs和CPUsetMems来设置容器可以使用的CPU和内存节点
		},
		// * 需要的其他主机配置
	}
	// TODO 检测名字是否存在
	resp, err := RunContainer(config, hostConfig, nil, "test")
	if err != nil {
		log.Printf("运行容器时出错 %v", err)
		return
	}
	log.Printf("容器 %s 已启动", resp.ID)

	jsonContainer, err := GetContainerInfo(resp.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("端口映射信息:")
	for privatePort, portBindings := range jsonContainer.NetworkSettings.Ports {
		for _, portBinding := range portBindings {
			fmt.Printf("内部端口 %s 映射到主机的IP地址 %s 及端口 %sn", privatePort, portBinding.HostIP, portBinding.HostPort)
		}
	}

	//// * 执行容器命令
	//out, err := execContainerCmd(resp.ID, []string{"/bin/sh", "-c", "echo flag{THis_is_a_test_flag} > /flag"})
	//if err != nil {
	//	log.Printf("执行容器命令时出错 %v", err)
	//	return
	//}
	//log.Printf("执行容器命令成功，结果： %s", out)

	//fmt.Println(port)
}

func FindRandomUnusedPort() (int, error) {
	//  * 设置最大端口和最小端口 这个到时候可以写入配置文件
	MinPort := 50000
	MaxPort := 60000
	//  * 生成种子 Rand.seed已经被废弃了
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		//  * 生成随机端口
		port := rand.Intn(MaxPort-MinPort) + MinPort

		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return 0, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			//  * 如果端口被占用，继续循环
			continue
		}
		defer l.Close()

		return port, nil
	}
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
func pullImage(imageName string) (string, error) {
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

type DockerImage struct {
	ID       string   `json:"ID"`
	RepoTags []string `json:"RepoTags"`
	Status   string   `json:"Status"`
	Created  int64    `json:"Created"`
	Size     int64    `json:"Size"`
}

// GetImagesList
//
//	@Description: 获取镜像列表
//	@return []string
//	@return error
func GetImagesList() ([]DockerImage, error) {
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

// RunContainer
//
//	@Description: 运行容器
//	@param cfg
//	@param hostConfig
//	@param networkConfig
//	@param containerName
//	@return container.CreateResponse
//	@return error
func RunContainer(cfg *container.Config, hostConfig *container.HostConfig, networkConfig *network.NetworkingConfig, containerName string) (container.CreateResponse, error) {
	// * platform 先传空，不知道有啥用，networkConfig 传空
	resp, err := cli.ContainerCreate(context.Background(), cfg, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		return container.CreateResponse{}, err
	}
	err = cli.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		return container.CreateResponse{}, err
	}
	return resp, nil
}

func execContainerCmd(containerId string, cmd []string) (string, error) {
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
func GetContainerInfo(containerId string) (types.ContainerJSON, error) {
	return cli.ContainerInspect(context.Background(), containerId)
}

// DeleteContainer
//
//	@Description: 删除容器
//	@param containerId
//	@return error
func DeleteContainer(containerId string) error {
	return cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
}
