using System;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
using GoSocket;

namespace Console
{
    
    public class Program
    {
        public static void ConnectHandler()
        {
            System.Console.Out.WriteLine("服务器连接成功");
        }

        public static void CloseHandler()
        {
            System.Console.Out.WriteLine("服务器关闭成功");
        }

        public static void SocketTcpClientErrorHandler(long code,IntPtr msg,long msglen)
        {
            byte[] bytes = new byte[msglen];
            for (int i = 0; i < msglen; i++)
            {
                bytes[i] = Marshal.ReadByte(msg, i);
            }
            string err = Encoding.UTF8.GetString(bytes);
            System.Console.Out.WriteLine("错误:"+code.ToString()+" "+err);
        }

        public static void SocketTcpClientMessageHandler(IntPtr data, int len)
        {
            byte[] msgdata = new GoSlice
            {
                data = data,
                len = len,
            };
            string msg = Encoding.UTF8.GetString(msgdata);
            System.Console.Out.WriteLine("收到消息:"+msg);
        }

        static void Main(string[] args)
        {
            //Environment.SetEnvironmentVariable("GODEBUG", "cgocheck=0");
            System.Console.Out.WriteLine("========================");
            SocketTcpClientConfig config = new SocketTcpClientConfig
            {
                //address = "10.0.200.76:8900",
                address = "127.0.0.1:8900",
                autoReConnect = 1,
                ReConnectWaitTime = 2,
                HeartBeatTime = 3,
                ConnectTimeOut = 9,
            };
            config.connectHandler = ConnectHandler;
            config.messageHandler = SocketTcpClientMessageHandler;
            config.closeHandler = CloseHandler;
            config.errorHandler = SocketTcpClientErrorHandler;
            string clientID = C.NewTcpClient(config);
            while(true)
            {
                Thread.Sleep(1000);
                string msg = "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao" +
               "这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao这里是中文消息ni hao";
                byte[] msgBytes = Encoding.UTF8.GetBytes(msg);
                C.TcpClientWriteMsg(clientID, msgBytes);
            }
        }
    }

}
