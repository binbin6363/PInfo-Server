package utils

import (
	"fmt"
	"sync"
	"time"

	"PInfo-server/log"
)

// copy from snowflake

const (
	// 前端js无法识别超大int64数，将毫秒改为秒，减小生成的数字
	epoch          = int64(1577808000)                 // 设置起始时间(时间戳/秒)：2020-01-01 00:00:00，有效期69年
	timestampBits  = uint(41)                          // 时间戳占用位数
	nodeIdBits     = uint(8)                           // 数据中心id所占位数
	sequenceBits   = uint(14)                          // 序列所占的位数
	timestampMax   = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	nodeIdMax      = int64(-1 ^ (-1 << nodeIdBits))    // 支持的最大数据中心id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits))  // 支持的最大序列id数量
	nodeIdShift    = sequenceBits                      // 数据中心id左移位数
	timestampShift = sequenceBits + nodeIdBits         // 时间戳左移位数
)

// Snowflake 雪花算法。ID组成：时间(41bit) + 机器ID(8bit) + 递增序列号(14bit)
type Snowflake struct {
	sync.Mutex
	epoch     int64 // 起始时间，纪元
	timestamp int64 // 当前时间戳
	nodeId    int64 // 当前节点ID
	sequence  int64 // 递增序列号
}

// NewSnowflake 构建雪花算法的实例
func NewSnowflake(nodeId int64) (*Snowflake, error) {
	if nodeId < 0 || nodeId > nodeIdMax {
		return nil, fmt.Errorf("nodeId must be between 0 and %d", nodeIdMax-1)
	}
	return &Snowflake{
		// 设置起始时间(时间戳/秒)：2020-01-01 00:00:00，有效期69年
		epoch:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UnixNano() / int64(time.Millisecond),
		timestamp: 0,
		nodeId:    nodeId,
		sequence:  0,
	}, nil
}

func (s *Snowflake) NextVal() int64 {
	s.Lock()
	defer s.Unlock()
	now := time.Now().UnixNano() / int64(time.Millisecond) // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出sequenceBits长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				log.Warn("sequence overflow, wait for next ms")
				now = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch // 减掉纪元
	if t > timestampMax {
		log.Errorf("epoch must be between 0 and %d", timestampMax-1)
		return 0
	}
	s.timestamp = now
	r := t<<timestampShift | (s.nodeId << nodeIdShift) | (s.sequence)
	return r
}

// GetNodeID 获取数据中心ID
func GetNodeID(sid int64) int64 {
	nodeId := (sid >> nodeIdShift) & nodeIdMax
	return nodeId
}

// GetTimestamp 获取ID中的纪元时间戳
func GetTimestamp(sid int64) int64 {
	timestamp := (sid >> timestampShift) & timestampMax
	return timestamp
}

// GetGenTimestamp 获取创建ID时的真实时间戳
func GetGenTimestamp(sid int64) int64 {
	timestamp := GetTimestamp(sid) + epoch
	return timestamp
}

// GetGenTime 获取创建ID时的时间字符串(精度：秒)
func GetGenTime(sid int64) string {
	// 需将GetGenTimestamp获取的时间戳/1000转换成秒
	t := time.Unix(GetGenTimestamp(sid)/int64(time.Millisecond), 0).Format("2006-01-02 15:04:05")
	return t
}

// GetTimestampStatus 获取时间戳已使用的占比：范围（0.0 - 1.0）
func GetTimestampStatus() float64 {
	state := float64(time.Now().UnixNano()/int64(time.Millisecond)-epoch) / float64(timestampMax)
	return state
}
