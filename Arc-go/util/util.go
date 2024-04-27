package util

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net"
	"strings"
	"time"
)

func RandomString(n int) string {
	/*
	* 生成随机字符串
	 */
	var letters = []byte("qwertyuiopasdfghjklxzcvbnmQWERTYUIOPSADFGHJKLZXCVBNM0123456789")
	result := make([]byte, n)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

type StringSlice []string

// Value / Value 实现 driver.Valuer 接口
func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}

// Scan
//
//	@Description: Scan 实现 sql.Scanner 接口
//	@receiver s
//	@param value
//	@return error
func (s *StringSlice) Scan(value interface{}) error {

	if value == nil {
		*s = []string{}
		return nil
	}

	switch value := value.(type) {
	case string:
		*s = strings.Split(value, ",")
	case []byte:
		*s = strings.Split(string(value), ",")
	default:
		return errors.New("invalid type for StringSlice")
	}
	return nil
}

func (s StringSlice) Contains(id StringSlice) bool {
	/*
	* 判断是否包含
	 */
	for _, i := range id {
		for _, j := range s {
			if i == j {
				return true
			}
		}
	}
	return false
}

func GenerateUUID() string {
	/*
	* 生成uuid
	 */
	randomUUID := uuid.New()
	return randomUUID.String()
}

// FindRandomUnusedPort
// ! 用不着似乎，docker生成时如果填0则自动随机端口
//
//	@Description: FindRandomUnusedPort 生成一个随机未被使用的端口
//	@return int
//	@return error
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
