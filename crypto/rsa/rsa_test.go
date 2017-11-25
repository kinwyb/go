package rsa

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_Res(t *testing.T) {
	convey.Convey("rsa加解密测试", t, func() {
		i := []int{1024, 2048, 4096}
		data := "加密测试"
		for _, x := range i {
			privatekey, publickey, err := GenRsaKey(x, PKCS8)
			if err != nil {
				t.Errorf("%d密钥生成失败:%v", x, err)
				return
			}
			// fmt.Printf("私钥:\r\n%s", string(privatekey))
			// fmt.Printf("公钥:\r\n%s", string(publickey))
			convey.Convey(fmt.Sprintf("%d密钥长度测试", x), func() {
				edata, err := Encrypt(publickey, []byte(data))
				if err != nil {
					t.Errorf("加密失败:%s", err.Error())
					return
				}
				d, err := Decrypt(privatekey, edata, PKCS8)
				if err != nil {
					t.Errorf("解密失败:%s", err.Error())
					return
				}
				convey.So(string(d), convey.ShouldEqual, data)
			})
			convey.Convey(fmt.Sprintf("%d密钥长度NoPadding测试", x), func() {
				edata, err := EncryptNoPadding(publickey, []byte(data))
				if err != nil {
					t.Errorf("加密失败:%s", err.Error())
					return
				}
				piv, _ := DecodePrivateKey(privatekey, PKCS8)
				d := DecryptNoPadding(piv, edata)
				// fmt.Printf("解密数据:%s\r", string(d))
				convey.So(string(d), convey.ShouldEqual, data)
			})
		}
	})
}
