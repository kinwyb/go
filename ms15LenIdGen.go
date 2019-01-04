package heldiamgo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//  twitter SnowFlake 算法生成15位订单号
//  为确保生成出来的数字为15位，那么前面17位必须固定
//	64 - 17固定值 = 47 位
//	47 - 38位时间码 = 9 位   						38位毫秒数=8.7年
//	47 - 5位机器码 = 4 位  						31台机器
//	4位数值										最大15

//毫秒级15位ID生成器
type Ms15LenIDGenerator struct {
	workerId           int64 //当前的workerId
	workerIdAfterShift int64 //移位后的workerId，可直接跟时间戳、序号取位或操作
	startTimestamp     int64 //起始时间戳
	lastMsTimestamp    int64 //上一次用的时间戳
	curSequence        int64 //当前的序号

	timeBitSize     uint8 //时间戳占的位数，默认为41位，最大不超过60位
	workerIdBitSize uint8 //workerId占的位数，默认10，最大不超过60位
	sequenceBitSize uint8 //序号占的位数，默认12，最大不超过60位
	lock            sync.Locker

	maxSequence        int64 //最后序列号最大值，初始化时计算出来的
	workerIdLeftShift  uint8 //生成的workerId只取最低的几位，这里要左移，给序列号腾位，初始化时计算出来的
	timestampLeftShift uint8 //生成的时间戳左移几位，给workId、序列号腾位，初始化时计算出来的
}

//实例化毫秒15位一个ID生成器
func NewMs15LenIDGenerator(workID int64) *Ms15LenIDGenerator {
	ret := &Ms15LenIDGenerator{
		lock:               &sync.Mutex{},
		workerId:           workID,
		lastMsTimestamp:    0,
		curSequence:        0,
		timeBitSize:        38, //默认的时间戳占的位数
		workerIdBitSize:    5,  //默认的workerId占的位数
		sequenceBitSize:    4,  //默认的序号占的位数
		maxSequence:        0,  //最大的序号值，初始化的时计算出来的
		workerIdLeftShift:  0,  //worker id左移位数
		timestampLeftShift: 0,
	}
	ret.workerIdAfterShift = workID & 0x1F << ret.sequenceBitSize //确保生成的id只有5位
	ret.maxSequence = 0xF                                         //最大生成值
	ret.workerIdLeftShift = ret.sequenceBitSize
	ret.timestampLeftShift = ret.sequenceBitSize + ret.workerIdBitSize
	ret.lastMsTimestamp = ret.genNextTs(ret.lastMsTimestamp)
	return ret
}

//生成时间戳位数,
func (mig *Ms15LenIDGenerator) genTs() int64 {
	rawTs := time.Now().UnixNano()/int64(time.Millisecond) - mig.startTimestamp
	diff := 64 - mig.timeBitSize
	// &0x7FFFFFFFFFFFFFFF 确保最高位为0,防止出现负数
	ret := (rawTs << diff & 0x7FFFFFFFFFFFFFFF) >> diff
	return ret << mig.timestampLeftShift
}

//生成下一个时间戳，如果时间戳的位数较小，且序号用完时此处等待的时间会较长
func (mig *Ms15LenIDGenerator) genNextTs(last int64) int64 {
	for {
		cur := mig.genTs()
		if cur > last {
			return cur
		}
	}
}

//生成下一个ID
func (mig *Ms15LenIDGenerator) NextID() (int64, error) {
	mig.lock.Lock()
	defer mig.lock.Unlock()
	//先判断当前的时间戳，如果比上一次的还小，说明出问题了
	curTs := mig.genTs()
	if curTs < mig.lastMsTimestamp {
		return 0, fmt.Errorf("系统时钟异常")
	}
	//如果跟上次的时间戳相同，则增加序号
	if curTs == mig.lastMsTimestamp {
		curSequence := atomic.AddInt64(&mig.curSequence, 1)
		curSequence = curSequence & mig.maxSequence
		//序号又归0即用完了，重新生成时间戳
		if curSequence != 0 {
			ret := 0x3000000000000 | curTs | mig.workerIdAfterShift | curSequence
			return ret, nil
		}
		curTs = mig.genNextTs(mig.lastMsTimestamp)
	}
	// 如果两个的时间戳不一样，则归0序号
	atomic.StoreInt64(&mig.curSequence, 0)
	atomic.StoreInt64(&mig.lastMsTimestamp, curTs)
	ret := 0x3000000000000 | curTs | mig.workerIdAfterShift
	return ret, nil
}
