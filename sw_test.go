package sw

import (
	"fmt"
	"testing"
	"github.com/gaochao1/gosnmp"
	"time"
)

const (
	ip        = "10.10.10.1"
	community = "public"
	oid       = "1.3.6.1.2.1.31.1.1.1.6"
	timeout   = 5
	method    = "walk"
	retry     = 5
	iprange   = "10.10.10.1/24"
	pingIp	  = "10.10.10.1"
	pingtimeout = 1000
)





func Test_CpuUtilization(t *testing.T) {
	if np, err := CpuUtilization(ip, community, timeout, retry); err != nil {
		t.Error(err)
	} else {
		t.Log("CpuUtilization :",np)
	}
}

func Test_MemUtilization(t *testing.T) {
	if np, err := MemUtilization(ip, community, timeout, retry); err != nil {
		t.Error(err)
	} else {
		fmt.Println("MemUtilization :", np)
	}
}

func Test_RunSnmp(t *testing.T) {
	var np []gosnmp.SnmpPDU
	var err error
	for i := 0; i < retry; i++ {
		np, err = RunSnmp(ip, community, oid, method, timeout)
		if len(np) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_RunSnmp :", np)
		for _,v := range np{
			fmt.Println("value:",v.Value.(uint64))
		}
	}
}

func Test_SysDescr(t *testing.T) {
	if np, err := SysDescr(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_SysDescr :", np)
	}
}

func Test_SysVendor(t *testing.T) {
	if np, err := SysVendor(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_SysVendor :", np)
	}
}

func Test_ListIfStats(t *testing.T) {
	ignoreIface := []string{"VLAN","VL","Vl"}
	ignorePkt := true
	if np, err := ListIfStats(ip, community, timeout, ignoreIface, retry, ignorePkt); err != nil {
		t.Error(err)
	} else {
	fmt.Println("value:", np)
        }
}
func Test_ListIfStatsSnmpWalk(t *testing.T) {
        ignoreIface := []string{"VLAN","VL","Vl"}
        ignorePkt := true
        if np, err := ListIfStatsSnmpWalk(ip, community, timeout, ignoreIface, retry,ignorePkt); err != nil {
                t.Error(err)
        } else {
        fmt.Println("value:", np)
        }
}
func Test_SysModel(t *testing.T) {
	if np, err := SysModel(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_SysModel :", np)
	}
}

func Test_SysName(t *testing.T) {
	if np, err := SysName(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_SysName :", np)
	}
}

func Test_SysUpTime(t *testing.T) {
	if np, err := SysUpTime(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_SysUpTime :", np)
	}
}

func Test_ConnectionStat(t *testing.T) {
	if np, err := ConnectionStat(ip, community, timeout, retry); err != nil {
		t.Error(err)
	} else {
		t.Log("ConnectionStat :",np)
	}
}

func Test_ParseIp(t *testing.T) {
	np := ParseIp(iprange)
	t.Log("aliveip:",np)
}

func Test_PingRtt(t *testing.T) {
	rtt, err := PingRtt(pingIp, pingtimeout)
	t.Log("rtt:",rtt)
	t.Log("err:",err)
}