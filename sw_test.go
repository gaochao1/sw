package sw

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/ulricqin/gosnmp"
)

const (
	ip           = "10.200.43.1"
	community    = "test"
	oid          = "1.3.6.1.4.1.2021.10.1.3.1"
	timeout      = 1000
	method       = "get"
	retry        = 3
	iprange      = "10.10.50.1-10.10.50.25"
	pingIp       = "10.10.10.1"
	pingtimeout  = 1000
	fastPingMode = true
	limitConn    = 1
)

func Test_CpuUtilization(t *testing.T) {
	if np, err := CpuUtilization(ip, community, timeout, retry); err != nil {
		t.Error(err)
	} else {
		t.Log("CpuUtilization :", np)
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
		fmt.Println("Test_RunSnmp :", &np)
		fmt.Println(np[0].Value.(string))

	}
}

func Test_RunSnmpswalk(t *testing.T) {
	var np []gosnmp.SnmpPDU
	var err error
	np, err = RunSnmpwalk(ip, community, oid, retry, timeout)

	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_RunSnmp :", &np)

	}
}

func Test_SysDescr(t *testing.T) {
	np, err := SysDescr(ip, community, retry, timeout)
	t.Error(err)
	version_number, err := strconv.ParseFloat(getVersionNumber(np), 32)
	t.Error(err)
	fmt.Println("Test_SysDescr :", np)
	fmt.Println("Version_number:", version_number)
}

func Test_SysVendor(t *testing.T) {
	if np, err := SysVendor(ip, community, retry, timeout); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Test_SysVendor :", np)
	}
}

func Test_ListIfStats(t *testing.T) {
	ignoreIface := []string{"Vl"}
	ignorePkt := true
	ignoreOperStatus := true
	ignoreMulticastPkt := true
	ignoreBroadcastPkt := true
	ignoreDiscards := true
	ignoreErrors := true
	ignoreUnknownProtos := true
	ignoreOutQLen := true
	if np, err := ListIfStats(ip, community, timeout, ignoreIface, retry, limitConn, ignorePkt, ignoreOperStatus, ignoreBroadcastPkt, ignoreMulticastPkt, ignoreDiscards, ignoreErrors, ignoreUnknownProtos, ignoreOutQLen); err != nil {
		t.Error(err)
	} else {
		fmt.Println("value:", np)
	}
}
func Test_ListIfStatsSnmpWalk(t *testing.T) {
	ignoreIface := []string{"VLAN", "VL", "Vl"}
	ignorePkt := true
	ignoreOperStatus := false
	ignoreMulticastPkt := false
	ignoreBroadcastPkt := false
	ignoreDiscards := false
	ignoreErrors := false
	ignoreUnknownProtos := false
	ignoreOutQLen := true
	if np, err := ListIfStatsSnmpWalk(ip, community, timeout, ignoreIface, retry, ignorePkt, ignoreOperStatus, ignoreBroadcastPkt, ignoreMulticastPkt, ignoreDiscards, ignoreErrors, ignoreUnknownProtos, ignoreOutQLen); err != nil {
		t.Error(err)
	} else {
		fmt.Println("value:", np)
	}
}
func Test_SysModel(t *testing.T) {
	if np, err := SysModel(ip, community, retry, timeout); err != nil {
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

func Test_ParseIp(t *testing.T) {
	np := ParseIp(iprange)
	t.Log("aliveip:", np)

}

func Test_PingRtt(t *testing.T) {
	rtt, err := PingRtt(pingIp, pingtimeout, fastPingMode)
	t.Log("rtt:", rtt)
	t.Log("err:", err)

}

func Test_Ping(t *testing.T) {
	r := Ping(pingIp, pingtimeout, fastPingMode)
	t.Log(r)
}
