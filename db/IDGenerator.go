package db

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

//IDGenerator 唯一序号生成对象
type IDGenerator struct {
	workerID      int           //机器ID
	sequence      uint32        //当前时间戳计数器
	workerIDBits  uint          //机器码字节数,默认4字节保存机器ID
	maxWorkerID   int64         //最大机器ID
	sequenceBits  uint          //计数器字节数,默认10个字节保存计数器
	workerIDShift uint          //机器码数据左移位数
	sequenceMask  int           //一微秒内可以产生的计数器值,达到该值后要等到下一微妙再生成
	lastTimestamp int64         //上次生成序号的时间戳
	timeString    string        //上一次生成序列号的时间
	timeDuration  time.Duration //时间单位
	lock          *sync.Mutex   //同步锁
}

//NewIDGenerator 生成一个IDGenerator对象
//@param workerID 机器编码会保存在结果中
//@param params[0] 指定多少位来计入单个时间点中可以生成的数量[默认是10]
func NewIDGenerator(workerID int, params ...int64) *IDGenerator {
	var sequenceBits uint
	if len(params) > 1 {
		sequenceBits = uint(params[0])
	} else {
		if len(params) > 0 {
			sequenceBits = uint(params[0])
		} else {
			sequenceBits = 10
		}
	}
	idw := &IDGenerator{
		workerID:     workerID,
		sequenceBits: sequenceBits,
		sequence:     0,
		workerIDBits: 10,
		timeDuration: time.Millisecond,
		lock:         &sync.Mutex{},
		timeString:   strings.Replace(time.Now().Format("20060102150405.00000"), ".", "", -1),
	}
	idw.maxWorkerID = -1 ^ -1<<idw.workerIDBits
	idw.sequenceMask = -1 ^ -1<<idw.sequenceBits
	idw.workerIDShift = idw.sequenceBits
	idw.lastTimestamp = -1
	return idw
}

//next 生成一个唯一ID
func (d *IDGenerator) next() int64 {
	d.lock.Lock()
	defer d.lock.Unlock()
	t := time.Now()
	timestamp := t.UnixNano() / int64(d.timeDuration)
	if d.lastTimestamp == timestamp {
		atomic.AddUint32(&d.sequence, 1)
		d.sequence = uint32(int(d.sequence) & d.sequenceMask)
		if d.sequence == 0 {
			d.tilNext()
		}
	} else {
		d.sequence = 0
		d.lastTimestamp = timestamp
		d.timeString = strings.Replace(t.Format("20060102150405.00000"), ".", "", -1)
	}
	i := int64(d.workerID<<d.workerIDShift) | int64(d.sequence)
	return i
}

//NextString 下一个唯一字符串
func (d *IDGenerator) NextString() string {
	return d.timeString + strconv.FormatInt(d.next(), 10)
}

//获取下一时间值
func (d *IDGenerator) tilNext() int64 {
	for {
		t := time.Now()
		timestamp := t.UnixNano() / int64(d.timeDuration)
		if timestamp > d.lastTimestamp {
			d.lastTimestamp = timestamp
			d.timeString = strings.Replace(t.Format("20060102150405.00000"), ".", "", -1)
			return timestamp
		}
	}
}
