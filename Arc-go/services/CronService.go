package services

import (
	"Arc/common"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log"
	"time"
)

type ContainerCleanupService struct {
	DockerClient *client.Client
	Timer        *time.Timer
}

// NewContainerCleanupService
//
//	@Description: 初始化容器清理服务
//	@return *ContainerCleanupService
func NewContainerCleanupService() (*ContainerCleanupService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &ContainerCleanupService{
		DockerClient: cli,
	}, nil
}

// Start
// ! 这里设置定时清理周期，后期放到config中
// @Description: 启动容器清理服务
// @receiver s
func (s *ContainerCleanupService) Start() {
	s.scheduledCleanup(1 * time.Minute)
	log.Println("Container cleanup service started")
}

// Stop
//
//	@Description: 停止容器清理服务
//	@receiver s
func (s *ContainerCleanupService) Stop() {
	if s.Timer != nil {
		s.Timer.Stop()
	}
	log.Println("Container cleanup service stopped")
}

// scheduledCleanup
//
//	@Description: 定时清理容器
//	@receiver s
//	@param interval
func (s *ContainerCleanupService) scheduledCleanup(interval time.Duration) {
	// * 先清理一次
	s.CleanupContainer()
	s.Timer = time.AfterFunc(interval, func() {
		s.CleanupContainer()
		s.scheduledCleanup(interval)
	})
}

// CleanupContainer
// ! 这里设置超过多久的容器会被清理，后期放到config中
//
//	@Description: 清理容器
//	@receiver s
func (s *ContainerCleanupService) CleanupContainer() {
	containers, err := s.DockerClient.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Printf("List containers failed because %s", err)
		return
	}

	for _, containerItem := range containers {
		startedAt, err := common.GetContainerStartTime(s.DockerClient, containerItem.ID)
		if err != nil {
			log.Printf("Get container %s start time failed because %s", containerItem.ID, err)
			continue
		}

		// ! 这里设置超过多久的容器会被清理，后期放到config中
		if startedAt > 1*time.Minute {
			//  * 先停止运行再关！
			err := common.StopContainer(s.DockerClient, containerItem.ID)
			if err != nil {
				log.Printf("Stop container %s failed because %s", containerItem.ID, err)
				continue
			}
			err = common.DeleteContainer(s.DockerClient, containerItem.ID)
			if err != nil {
				log.Printf("Delete container %s failed because %s", containerItem.ID, err)
			} else {
				log.Printf("Container %s cleaned up", containerItem.ID)
			}
		}
	}
}

func main() {
	service, err := NewContainerCleanupService()
	if err != nil {
		log.Fatalf("Create container cleanup service failed because %s", err)
	}

	service.Start()
	defer service.Stop()

	select {
	// 无限等待
	}
}
