package psutil

/**
    Author: luzequan
    Created: 2018-02-11 16:44:55
*/
import (
	"testing"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"fmt"
	"os"
)

func Test_CPU(t *testing.T) {
	fmt.Println(cpu.Info())
}

func Test_Disk(t *testing.T) {

	file := "D:/dds_send/tmp/JON20180108000000233_ID010201_20180205000000_00111.TARGET"

	f, e := os.Stat(file)
	if e != nil {
	}
	fmt.Println(f.Size())

	fmt.Println(disk.Partitions(true))
}

func Test_Docker(t *testing.T) {

	fmt.Println(docker.GetDockerStat())
}

func Test_Host(t *testing.T) {

	fmt.Println(host.Info())
}

func Test_Mem(t *testing.T) {

	fmt.Println(mem.SwapMemory())
}

func Test_Net(t *testing.T) {

	fmt.Println(net.FilterCounters())
}

func Test_Process(t *testing.T) {

	fmt.Println(process.Processes())
}

func Test_Load(t *testing.T) {

	fmt.Println(load.Avg())
}