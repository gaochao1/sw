package sw

import (
	"fmt"
	"testing"
)

const (
	ip        = "10.10.23.1"
	community = "public"
	oid       = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.5"
	timeout   = 5
	method    = "get"
	retry     = 5
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
	if np, err := RunSnmp(ip, community, oid, method, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_RunSnmp :", np)
		for _,v := range np{
			fmt.Println("value:",v.Value.(int))
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
	if np, err := ListIfStats(ip, community, timeout, ignoreIface, retry,ignorePkt); err != nil {
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
