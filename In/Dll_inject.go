package In

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const Banner  = `     ██     ██     ██████    ██        ██             ██                   ██  
    ████   ░██    ░░░░░░█   █░█       ░██            ░░                   ░██  
   ██░░██  ░██  ██     ░█  █ ░█       ░██ ███████     ██  █████   █████  ██████
  ██  ░░██ ░██ ██      █  ██████      ░██░░██░░░██   ░██ ██░░░██ ██░░░██░░░██░ 
 ██████████░████      █  ░░░░░█       ░██ ░██  ░██   ░██░███████░██  ░░   ░██  
░██░░░░░░██░██░██    █       ░█       ░██ ░██  ░██ ██░██░██░░░░ ░██   ██  ░██  
░██     ░██░██░░██  █        ░█  █████░██ ███  ░██░░███ ░░██████░░█████   ░░██ 
░░      ░░ ░░  ░░  ░         ░  ░░░░░ ░░ ░░░   ░░  ░░░   ░░░░░░  ░░░░░     ░░  


Ver 1.0
`


func DLL_Inject(){

	fmt.Println(Banner)
	Dll_Path := flag.String("dll", "", "输入DLL的路径")
	Proc_ID := flag.Int("pid", 0, "要注入进程的PID")
	flag.Parse()

	dll_path, pid := *Dll_Path, *Proc_ID
	if dll_path != "" && pid > 0 {
		_, err := ioutil.ReadFile(dll_path)
		if err != nil || !strings.HasSuffix(dll_path, ".dll") {
			fmt.Println("Invalid Dll Path or PID range")
			os.Exit(1)
		}
	}


	dll_path_len := len(dll_path) + 1
	//取进程句柄
	HProcess:=pHandle(pid)

	//改变进程中内存区域的保护属性
	vaddr:=VirtualAlloc_Ex(HProcess,dll_path_len)

	Writepromemory(HProcess,vaddr,Ptr(dll_path),dll_path_len)

	loaddr :=GetPrAddr("LoadLibraryA")

	HThread:=Creath(HProcess,loaddr,vaddr)
	fmt.Printf("线程：%v",HThread)

	free :=VirtualFree_Ex(HProcess,vaddr,0)
	if free == 0{
		fmt.Printf("\n[-] DLL成功注入失败！\n")
		fmt.Printf("\n[-] 进程ID%v！\n",pid)

	}else {
		fmt.Printf("\n[+] DLL成功注入到进程！\n")
		fmt.Printf("\n[+] 进程ID%v！\n",pid)

	}







}

