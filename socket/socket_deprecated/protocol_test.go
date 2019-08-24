package socket_deprecated

import (
	"testing"

	"math/rand"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewProtocol(t *testing.T) {
	//HeartBeatByte := []byte("@hearbeat@")
	hb := []byte{0x11, 0x22, 0x13, 0x24, 0x15, 0x26, 0x17, 0x28, 0x19, 0x11}
	fmt.Printf("%x\n", 't')
	fmt.Printf("%s", hb)
}

func TestProtocol_intToByte(t *testing.T) {
	convey.Convey("IntByte转换", t, func() {
		for i := 0; i < 10000000; i++ {
			c := rand.Int63()
			b := intToByte(int64(c))
			x := byteToInt(b)
			convey.So(c, convey.ShouldEqual, x)
		}
	})
}
