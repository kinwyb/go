//
//  socketCCode.h
//  socketCCode
//
//  Created by 王迎宾 on 2019/10/5.
//

#ifndef socketCCode_h
#define socketCCode_h

#include "GoType.h"
#include <stdio.h>

// 连接成功回调函数
typedef void (*SocketTcpClientConnectHandler)();

// 收到消息回调函数
typedef void (*SocketTcpClientMessageHandler)(char *data,GoInt len);

// 连接关闭回调函数
typedef void (*SocketTcpClientCloseHandler)();

// 错误回调
typedef void (*SocketTcpClientErrorHandler)(GoInt code,char *msg,GoInt msgLen);

// tcp配置结构体
typedef struct {
    _GoString_ ID; //客户端标示
    _GoString_ address; //服务器地址
    GoInt autoReConnect; //是否自动重连(0,1)
    GoInt ReConnectWaitTime; //重连间隔时间(秒)
    GoInt HeartBeatTime; //心跳间隔时间(秒)
    GoSlice HeartBeatData; //心跳内容
    GoInt ConnectTimeOut; //连接超时时间(秒)
    SocketTcpClientConnectHandler connectHandler; //连接成功回调
    SocketTcpClientMessageHandler messageHandler; //收到消息回调
    SocketTcpClientCloseHandler closeHandler; //关闭连接回调
    SocketTcpClientErrorHandler errorHandler; //发生错误回调
} SocketTcpClientConfig;

// 连接成功触发函数
void TcpClientConnectCallback(SocketTcpClientConfig *config);

// 收到消息触发函数
void TcpClientMessageCallback(SocketTcpClientConfig *config,char *data,GoInt len);

// 关闭触发函数
void TcpClientCloseCallback(SocketTcpClientConfig *config);

// 错误触发函数
void TcpClientErrorCallback(SocketTcpClientConfig *config,GoInt errCode,char *msg,GoInt msgLen);

// 关闭连接
void TcpClientClose(SocketTcpClientConfig *config);

// 连接服务器
void TcpClientConnect(SocketTcpClientConfig *config);

void GoNewTcpClient(void* p0, _GoString_ p1, _GoString_ p2, GoInt p3, GoInt p4, GoInt p5, GoInt p6, GoSlice p7);

void GoTcpClientClose(_GoString_ p0);

#endif /* socketCCode_h */
