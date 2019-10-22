//
//  socketCCode.c
//  socketCCode
//
//  Created by 王迎宾 on 2019/10/5.
//

#include <stdlib.h>
#include "socketCCode.h"

// 连接成功触发函数
void TcpClientConnectCallback(SocketTcpClientConfig *config) {
    config->connectHandler();
}

// 收到消息触发函数
void TcpClientMessageCallback(SocketTcpClientConfig *config,char *data,GoInt len) {
    config->messageHandler(data,len);
    free(data);
}

// 关闭触发函数
void TcpClientCloseCallback(SocketTcpClientConfig *config) {
    config->closeHandler();
}

// 错误触发函数
void TcpClientErrorCallback(SocketTcpClientConfig *config,GoInt errCode,char *msg,GoInt msgLen) {
    config->errorHandler(errCode,msg,msgLen);
    free(msg);
}

// 关闭连接
void TcpClientClose(SocketTcpClientConfig *config) {
    GoTcpClientClose(config->ID);
}

// 连接服务器
void TcpClientConnect(SocketTcpClientConfig *config) {
    GoNewTcpClient(config,config->ID,config->address,config->autoReConnect,
    config->ReConnectWaitTime,config->HeartBeatTime,config->ConnectTimeOut,
    config->HeartBeatData);
}