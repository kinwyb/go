using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;
using System.Text;

namespace GoSocket
{
    [UnmanagedFunctionPointer(CallingConvention.Cdecl)]
    public delegate void SocketTcpClientConnectHandler();

    [UnmanagedFunctionPointer(CallingConvention.Cdecl)]
    public delegate void SocketTcpClientMessageHandler(IntPtr data,int len);

    [UnmanagedFunctionPointer(CallingConvention.Cdecl)]
    public delegate void SocketTcpClientCloseHandler();

    [UnmanagedFunctionPointer(CallingConvention.Cdecl)]
    public delegate void SocketTcpClientErrorHandler(long code, IntPtr msg, long msglen);

    public struct SocketTcpClientConfig
    {
        public GoString ID;
        public GoString address;
        public long autoReConnect;
        public long ReConnectWaitTime;
        public long HeartBeatTime;
        public GoSlice HeartBeatData;
        public long ConnectTimeOut;
        public SocketTcpClientConnectHandler connectHandler; //连接成功回调
        public SocketTcpClientMessageHandler messageHandler; //收到消息回调
        public SocketTcpClientCloseHandler closeHandler; //关闭连接回调
        public SocketTcpClientErrorHandler errorHandler; //发生错误回调
        IntPtr ptr;

        void FreeIntPtr()
        {
            if(ptr != null)
            {
                Marshal.FreeHGlobal(ptr);
            }
        }

        public static implicit operator IntPtr (SocketTcpClientConfig s)
        {
            int nSizeOfPerson = Marshal.SizeOf(s);
            s.ptr = Marshal.AllocHGlobal(nSizeOfPerson);
            Marshal.StructureToPtr(s, s.ptr, false);
            return s.ptr;
        }

    }

    public class C
    {
        [DllImport("F:/Go/user_src/src/socketC/gosocket.dll", CallingConvention = CallingConvention.Cdecl)]
        static extern void TcpClientConnect(IntPtr SocketTcpClientConfigP);

        [DllImport("F:/Go/user_src/src/socketC/gosocket.dll", CallingConvention = CallingConvention.Cdecl)]
        public static extern void TcpClientClose(IntPtr SocketTcpClientConfigP);

        [DllImport("F:/Go/user_src/src/socketC/gosocket.dll", CallingConvention = CallingConvention.Cdecl)]
        static extern void TcpClientWrite(GoString id, GoSlice data);

        public static string NewTcpClient(SocketTcpClientConfig config)
        {
            string ID = System.Guid.NewGuid().ToString("N");
            config.ID = ID;
            TcpClientConnect(config);
            return ID;
        }

        public static void TcpClientWriteMsg(string id,byte[] data)
        {
            GoSlice slice = data;
            GoString idString = id;
            TcpClientWrite(idString, slice);
            idString.FreeMem();
            slice.FreeMem();
        }
    }
}
