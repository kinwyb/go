package socket_deprecated

//数据打包处理器
type PackageHandler interface {
	//打包
	Package(msg []byte) []byte
	//解包
	UnPackage(msg []byte) []byte
}
