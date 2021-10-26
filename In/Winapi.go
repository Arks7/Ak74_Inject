package In

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
	"unsafe"
)

var(
	kernel32      = windows.NewLazyDLL("kernel32.dll")
	ntdll     = windows.NewLazyDLL("ntdll.dll")


	RtlCreateHeap = ntdll.NewProc("RtlCreateHeap")
	RtlAllocateHeap = ntdll.NewProc("RtlAllocateHeap")
	RtlCopyMemory = ntdll.NewProc("RtlCopyMemory")
	VirtualAlloc  = kernel32.NewProc("VirtualAlloc")
	OpenProcess        = kernel32.NewProc("OpenProcess")
	VirtualFreeEx      = kernel32.NewProc("VirtualFreeEx")
	VirtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
	WriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	GetProcAddress     = kernel32.NewProc("GetProcAddress")
	CreateRemoteThread = kernel32.NewProc("CreateRemoteThread")

)

const (
	PROCESS_ALL_ACCESS     = 0x1F0FFF
	MEM_COMMIT             = 0x00001000
	MEM_RESERVE            = 0x00002000
	MEM_RESERVE_AND_COMMIT = MEM_COMMIT | MEM_RESERVE
	PAGE_READWRITE         = 0x04
	MEM_RELEASE            = 0x00008000
	PAGE_EXECUTE_READWRITE = 0x40
)




//获取句柄
func pHandle(pid int) uintptr {
	pHandle, _, _ := OpenProcess.Call(Ptr(PROCESS_ALL_ACCESS), Ptr(0), Ptr(pid))
	if pHandle == 0 {
		fmt.Println("打开进程句柄失败！")
		os.Exit(1)
	}
	return pHandle
}

//改变进程中内存区域的保护属性
func VirtualAlloc_Ex(Hp uintptr,len int)uintptr{
	addr, _, err := VirtualAllocEx.Call(Hp, Ptr(0), Ptr(len), Ptr(MEM_RESERVE_AND_COMMIT), Ptr(PAGE_READWRITE))
	if addr == 0 {
		fmt.Println("申请失败！")
		fmt.Println(addr, err)
		os.Exit(1)
	}
	return addr
}

//写内存
func Writepromemory(Hp uintptr,addr uintptr,lpBuffer uintptr,len int)uintptr{
	var bytesWritten byte
	IsMemoryWritten, _, _ := WriteProcessMemory.Call(Hp, addr, lpBuffer, Ptr(len), Ptr(unsafe.Pointer(&bytesWritten)))
	if IsMemoryWritten == 0 {
		fmt.Println("Dll path failed to write ")
		os.Exit(1)
	}
	return IsMemoryWritten
}

//获取函数的地址
func GetPrAddr(val string) uintptr{
	loadLibraryAddress, _, _ := GetProcAddress.Call(kernel32.Handle(), Ptr(val))
	if loadLibraryAddress == 0 {
		fmt.Println("获取LoadLibraryA函数的地址失败!")
		os.Exit(1)
	}
	return loadLibraryAddress

}

//创建远程线程
func Creath(Hp uintptr,addr uintptr,viraddr uintptr)uintptr{
	var threadid int
	HThread, _, _ := CreateRemoteThread.Call(Hp, Ptr(nil), Ptr(0), addr, viraddr, Ptr(0), uintptr(unsafe.Pointer(&threadid)))
	if HThread == 0 {
		fmt.Println("远程线程创建失败")
		os.Exit(1)
	}
	return HThread

}

//进程中释放申请的虚拟内存空间
func VirtualFree_Ex(HProcess uintptr,virtualAddress uintptr,bytesize int)uintptr{
	freed, _, _ := VirtualFreeEx.Call(HProcess, virtualAddress, Ptr(bytesize), Ptr(MEM_RELEASE))
	if freed == 0 {
		fmt.Println("虚拟内存释放失败！. ")
		os.Exit(1)
	}
	return freed
}

//转换uintptr
func Ptr(val interface{}) uintptr {
	switch val.(type) {
	case byte:
		return uintptr(val.(byte))
	case bool:
		isTrue := val.(bool)
		if isTrue {
			return uintptr(1)
		}
		return uintptr(0)
	case string:
		bytePtr, _ := syscall.BytePtrFromString(val.(string))
		return uintptr(unsafe.Pointer(bytePtr))
	case int:
		return uintptr(val.(int))
	case uint:
		return uintptr(val.(uint))
	default:
		return uintptr(0)
	}
}
