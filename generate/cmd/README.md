# generate工具
> 通过工具可以根据接口定义生成rpcx服务代码和客户端代码
> 
> 
## 生成rpcx服务端代码
> 用法: ./generate rpcx -i 接口文件文件夹 -o 生成文件的文件夹  
> 示例:

```
type ColorService interface {

	//添加颜色档案,code编码,name名称
	Add(code, name, creator string, dbtag string) err1.Error

	//更新颜色档案,code要更新的颜色档案编码,name名称,modifier 修改人
	//state 更新的状态
	Update(code, name, modifier string,
		state objs.EnableState, dbtag string) (string,err1.Error)

	//更新颜色档案名称,code要更新的颜色档案编码,name名称,modifier 修改人
	UpdateName(code, name, modifier string, dbtag string) err1.Error

	//更新颜色档案状态,code要更新的颜色档案编码,name名称,modifier 修改人
	UpdateEnableState(code, modifier string,
		state objs.EnableState, dbtag string) err1.Error

	//根据编码查询颜色档案 code要查询的编码
	QueryByCode(code string) *objs.ArchivesColor

	//根据名称模糊查询颜色档案 name要查询的名称
	QueryByLikeName(name string, dbtag string) []*objs.ArchivesColor

	//获取颜色档案列表
	QueryList(codeOrName string, state objs.EnableState, pg *db.PageObj, dbtag string) []*objs.ArchivesColor
}
```

> 生成结果：

```
type ColorServiceRpcx struct {
	serv endPoints.ColorService
}
type AddRequest struct {
	Code    string
	Name    string
	Creator string
	Dbtag   string
}
type AddResponse struct {
	E err1.Error
}

func (c *ColorServiceRpcx) Add(ctx context.Context, arg *AddRequest, resp *AddResponse) error {
	resp.E = c.serv.Add(arg.Code, arg.Name, arg.Creator, arg.Dbtag)
	return nil
}

type UpdateRequest struct {
	Code     string
	Name     string
	Modifier string
	State    objs.EnableState
	Dbtag    string
}
type UpdateResponse struct {
	S string
	E err1.Error
}

func (c *ColorServiceRpcx) Update(ctx context.Context, arg *UpdateRequest, resp *UpdateResponse) error {
	resp.S, resp.E = c.serv.Update(arg.Code, arg.Name, arg.Modifier, arg.State, arg.Dbtag)
	return nil
}

type UpdateNameRequest struct {
	Code     string
	Name     string
	Modifier string
	Dbtag    string
}
type UpdateNameResponse struct {
	E err1.Error
}

func (c *ColorServiceRpcx) UpdateName(ctx context.Context, arg *UpdateNameRequest, resp *UpdateNameResponse) error {
	resp.E = c.serv.UpdateName(arg.Code, arg.Name, arg.Modifier, arg.Dbtag)
	return nil
}

type UpdateEnableStateRequest struct {
	Code     string
	Modifier string
	State    objs.EnableState
	Dbtag    string
}
type UpdateEnableStateResponse struct {
	E err1.Error
}

func (c *ColorServiceRpcx) UpdateEnableState(ctx context.Context, arg *UpdateEnableStateRequest, resp *UpdateEnableStateResponse) error {
	resp.E = c.serv.UpdateEnableState(arg.Code, arg.Modifier, arg.State, arg.Dbtag)
	return nil
}

type QueryByCodeResponse struct {
	A *objs.ArchivesColor
}

func (c *ColorServiceRpcx) QueryByCode(ctx context.Context, arg string, resp *QueryByCodeResponse) error {
	resp.A = c.serv.QueryByCode(arg)
	return nil
}

type QueryByLikeNameRequest struct {
	Name  string
	Dbtag string
}
type QueryByLikeNameResponse struct {
	S []*objs.ArchivesColor
}

func (c *ColorServiceRpcx) QueryByLikeName(ctx context.Context, arg *QueryByLikeNameRequest, resp *QueryByLikeNameResponse) error {
	resp.S = c.serv.QueryByLikeName(arg.Name, arg.Dbtag)
	return nil
}

type QueryListRequest struct {
	CodeOrName string
	State      objs.EnableState
	Pg         *db.PageObj
	Dbtag      string
}
type QueryListResponse struct {
	S []*objs.ArchivesColor
}

func (c *ColorServiceRpcx) QueryList(ctx context.Context, arg *QueryListRequest, resp *QueryListResponse) error {
	resp.S = c.serv.QueryList(arg.CodeOrName, arg.State, arg.Pg, arg.Dbtag)
	return nil
}
```
> 使用：

```
   s := server.NewServer()
	//etcd
	r := &plugins.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@:%d", port),
		EtcdServers:    etcdEndPoint,
		BasePath:       "/rpcx/",
		Metrics:        metrics.NewRegistry(),
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
	s.RegisterName("color", & ColorServiceRpcx{
		serv: xxxxx,
	}, "")
	s.Serve("tcp",":1093")
```

> ## RPCX客户端代码生成
> 用法: ./generate rpcx -i 接口文件文件夹 -o 生成文件的文件夹
> 
> 生成结果:

```
type ColorServiceRpcxClient struct {
	client client.XClient
}
type AddRequest struct {
	Code    string
	Name    string
	Creator string
	Dbtag   string
}
type AddResponse struct {
	E err1.Error
}

func (c *ColorServiceRpcxClient) Add(Code string, Name string, Creator string, Dbtag string) err1.Error {
	arg := &AddRequest{Code: Code, Name: Name, Creator: Creator, Dbtag: Dbtag}
	reply := &AddResponse{}
	err := c.client.Call(context.Background(), "Add", arg, reply)
	if err != nil {
		log.Error("RPCX调用错误:%s", err.Error())
	}
	return reply.E
}

type UpdateRequest struct {
	Code     string
	Name     string
	Modifier string
	State    objs.EnableState
	Dbtag    string
}
type UpdateResponse struct {
	S string
	E err1.Error
}

func (c *ColorServiceRpcxClient) Update(Code string, Name string, Modifier string, State objs.EnableState, Dbtag string) (string, err1.Error) {
	arg := &UpdateRequest{Code: Code, Name: Name, Modifier: Modifier, State: State, Dbtag: Dbtag}
	reply := &UpdateResponse{}
	err := c.client.Call(context.Background(), "Update", arg, reply)
	if err != nil {
		log.Error("RPCX调用错误:%s", err.Error())
	}
	return reply.S, reply.E
}

type UpdateNameRequest struct {
	Code     string
	Name     string
	Modifier string
	Dbtag    string
}
type UpdateNameResponse struct {
	E err1.Error
}

func (c *ColorServiceRpcxClient) UpdateName(Code string, Name string, Modifier string, Dbtag string) err1.Error {
	arg := &UpdateNameRequest{Code: Code, Name: Name, Modifier: Modifier, Dbtag: Dbtag}
	reply := &UpdateNameResponse{}
	err := c.client.Call(context.Background(), "UpdateName", arg, reply)
	if err != nil {
		log.Error("RPCX调用错误:%s", err.Error())
	}
	return reply.E
}

type UpdateEnableStateRequest struct {
	Code     string
	Modifier string
	State    objs.EnableState
	Dbtag    string
}
type UpdateEnableStateResponse struct {
	E err1.Error
}

func (c *ColorServiceRpcxClient) UpdateEnableState(Code string, Modifier string, State objs.EnableState, Dbtag string) err1.Error {
	arg := &UpdateEnableStateRequest{Code: Code, Modifier: Modifier, State: State, Dbtag: Dbtag}
	reply := &UpdateEnableStateResponse{}
	err := c.client.Call(context.Background(), "UpdateEnableState", arg, reply)
	if err != nil {
		log.Error("RPCX调用错误:%s", err.Error())
	}
	return reply.E
}

type QueryByCodeResponse struct {
	A *objs.ArchivesColor
}

func (c *ColorServiceRpcxClient) QueryByCode(Code string) *objs.ArchivesColor {
	arg := Code
	reply := &QueryByCodeResponse{}
	err := c.client.Call(context.Background(), "QueryByCode", arg, reply)
	if err != nil {
		log.Error("RPCX调用错误:%s", err.Error())
	}
	return reply.A
}

type QueryByLikeNameRequest struct {
	Name  string
	Dbtag string
}
type QueryByLikeNameResponse struct {
	S []*objs.ArchivesColor
}

func (c *ColorServiceRpcxClient) QueryByLikeName(Name string, Dbtag string) []*objs.ArchivesColor {
	arg := &QueryByLikeNameRequest{Name: Name, Dbtag: Dbtag}
	reply := &QueryByLikeNameResponse{}
	err := c.client.Call(context.Background(), "QueryByLikeName", arg, reply)
	if err != nil {
		log.Error("RPCX调用错误:%s", err.Error())
	}
	return reply.S
}

type QueryListRequest struct {
	CodeOrName string
	State      objs.EnableState
	Pg         *db.PageObj
	Dbtag      string
}
type QueryListResponse struct {
	S []*objs.ArchivesColor
}

func (c *ColorServiceRpcxClient) QueryList(CodeOrName string, State objs.EnableState, Pg *db.PageObj, Dbtag string) []*objs.ArchivesColor {
	arg := &QueryListRequest{CodeOrName: CodeOrName, State: State, Pg: Pg, Dbtag: Dbtag}
	reply := &QueryListResponse{}
	err := c.client.Call(context.Background(), "QueryList", arg, reply)
	if err != nil {
		log.Error("RPCX调用错误:%s", err.Error())
	}
	return reply.S
}
```  

> 使用:

```
c := &ColorServiceRpcxClient{
	client:xxxxxx
}
ret := c.QueryList(arg....)
```

# 生成的代码引入的包可能有问题，需要手动处理