package rsa

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_Crypto(t *testing.T) {
	convey.Convey("长数据rsa加解密测试", t, func() {
		data := "xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问xi测试sfiwr内陆骄傲和任何 i 哦x，。待审核完了我dfsjfsjerJF9fsaihf守护六六大顺罚款返回舒服12904h08FIWR9i90jfknklnvhau2khr58y生活垃圾佛 i 未来us09tj2o[dskafsdaknvsdhyf2rkvdsaufg90afksfjowr920u3rjdfvjdsiosagvuydsvnDNHG待开始恢复 i 暗花纹哦我居然可能很多事是 v 引发事故后来个佛问"
		i := []int{1024, 2048, 4096}
		for _, x := range i {
			privatekey, publickey, err := GenRsaKey(x, PKCS8)
			if err != nil {
				t.Errorf("%d密钥生成失败:%v", x, err)
				return
			}
			convey.Convey(fmt.Sprintf("%d密钥长度测试", x), func() {
				edata, err := EncodeData(publickey, []byte(data))
				if err != nil {
					t.Errorf("%d加密失败:%s", x, err.Error())
					return
				}
				d, err := DecodeData(privatekey, edata, PKCS8)
				if err != nil {
					t.Errorf("%d解密失败:%s", x, err.Error())
					return
				}
				convey.So(string(d), convey.ShouldEqual, data)
			})
		}
	})
}
