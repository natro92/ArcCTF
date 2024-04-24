package services

import (
	"Arc/model"
	"log"
	"sync"
)

// CheckFlagServiceInterface
// @Description: 检查 flag 服务接口
type CheckFlagServiceInterface interface {
	Start()
	Stop()
	CheckFlags()
}

// CheckFlagService
// @Description: 检查 flag 服务
type CheckFlagService struct {
	submissionChan chan model.Submission
	stopChan       chan bool
	wg             sync.WaitGroup
	logger         *log.Logger
}

// NewCheckFlagService
//
//	@Description: 初始化检查 flag 服务
//	@param logger
//	@return *CheckFlagService
func NewCheckFlagService(logger *log.Logger) *CheckFlagService {
	return &CheckFlagService{
		submissionChan: make(chan model.Submission, 100),
		stopChan:       make(chan bool),
		logger:         logger,
	}
}

// Start
//
//	@Description: 启动检查 flag 服务
//	@receiver f
func (f *CheckFlagService) Start() {
	for i := 0; i < 2; i++ {
		f.wg.Add(1)
		go f.CheckFlags(i)
	}
	f.logger.Println("Check flag service started")
}

// Stop
//
//	@Description: 停止检查 flag 服务
//	@receiver f
func (f *CheckFlagService) Stop() {
	close(f.stopChan)
	f.wg.Wait()
	f.logger.Println("Check flag service stopped")
}

// CheckFlags
//
//	@Description: 检查 flag
//	@receiver f
func (f *CheckFlagService) CheckFlags(id int) {
	defer f.wg.Done()
	for {
		select {
		case <-f.stopChan:
			f.logger.Printf("Service %s stopped\n", id)
			return
		case submission := <-f.submissionChan:
			f.processSubmission(submission, id)
		}
	}
}

func (f *CheckFlagService) processSubmission(submission model.Submission, workerID int) {
	f.logger.Printf("Checking flag for submission %d\n", submission.ID)
	var status int
	status = 1
	// * 检查 flag
	submission.Status = status
	// * 如果 flag 正确则更新数据库
	// * 如果 flag 错误则不更新数据库
	f.logger.Printf("Worker %d: Processed submission %s with status %dn", workerID, submission.ID, status)
}

func main() {
	logger := log.Default()
	newCheckFlagService := NewCheckFlagService(logger)
	newCheckFlagService.Start()

	defer newCheckFlagService.Stop()

	select {}
}
