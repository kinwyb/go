using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;
using System.Text;

namespace GoSocket
{
    public struct GoString
    {
        public IntPtr p;

        public int n;

        public void FreeMem()
        {
            if (p != null)
            {
                Marshal.FreeHGlobal(p);
            }
        }

        public static implicit operator GoString(string s)
        {
            return new GoString
            {
                n = Encoding.Unicode.GetBytes(s).Length,
                p = Marshal.StringToHGlobalUni(s),
            };
        }

        public static implicit operator string(GoString s)
        {
            byte[] bytes = new byte[s.n];
            for (int i = 0; i < s.n; i++)
            {
                bytes[i] = Marshal.ReadByte(s.p, i);
            }
            return Encoding.UTF8.GetString(bytes);
        }

    }

    public struct GoSlice
    {
        public IntPtr data;
        public long len;
        public long cap;

        public void FreeMem()
        {
            if(data != null)
            {
                Marshal.FreeHGlobal(data);
            }
        }

        public static implicit operator byte[] (GoSlice s)
        {
            byte[] bytes = new byte[s.len];
            for (int i = 0; i < s.len; i++)
            {
                bytes[i] = Marshal.ReadByte(s.data, i);
            }
            return bytes;
        }

        public static implicit operator GoSlice(byte[] data)
        {
            IntPtr p = Marshal.AllocHGlobal(data.Length);
            Marshal.Copy(data, 0, p, data.Length);
            GoSlice ret = new GoSlice
            {
                data = p,
                len = data.Length
            };
            return ret;
        }

    }

}
