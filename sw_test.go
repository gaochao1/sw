package sw

import (
	"fmt"
	"testing"
)

const (
	ip        = "172.16.114.129"
	community = "public"
	oid       = "1.3.6.1.2.1.1.5.0"
	timeout   = 1
	method    = "get"
)

func Test_ListIfStats(t *testing.T) {
	//	onlyPrefix := []string{"eth"}
	onlyPrefix := []string{}

	if np, err := ListIfStats(ip, community, timeout, onlyPrefix); err != nil {
		t.Error(err)
	} else {
		fmt.Println("ListIfStats :", np)
	}
}

func Test_ListIfName(t *testing.T) {
	if np, err := ListIfName(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("ListIfName :", np)
	}
}

func Test_ListIfHCInOctets(t *testing.T) {
	if np, err := ListIfHCInOctets(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("ListIfHCInOctet :", np)
	}
}

func Test_ListIfHCOutOctets(t *testing.T) {
	if np, err := ListIfHCOutOctets(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("ListIfHCOutOctet :", np)
	}
}

func Test_CpuUtilization(t *testing.T) {
	if np, err := CpuUtilization(ip, community, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("CpuUtilization :", np)
	}
}

func Test_MemUtilization(t *testing.T) {
	if np, err := MemUtilization(ip, community, timeout); err != nil {
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
