package sw

import (
	"github.com/alouca/gosnmp"
	"time"
)

func CpuUtilization(ip, community string, timeout, retry int) (int, error) {
	vendor, err := SysVendor(ip, community, timeout)
	method := "get"
	var oid string

	switch vendor {
	case "Cisco_NX":
		oid = "1.3.6.1.4.1.9.9.305.1.1.1.0"
	case "Cisco":
		oid = "1.3.6.1.4.1.9.9.109.1.1.1.1.7.1"
	case "Huawei":
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.5.16842753"
	case "H3C":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.6.74"
	case "H3C_V5":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.6.74"
	case "H3C_V7":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.6.212"
	default:
		return 0, err
	}

	var snmpPDUs []gosnmp.SnmpPDU
	for i := 0; i < retry; i++ {
		snmpPDUs, err = RunSnmp(ip, community, oid, method, timeout)
		if len(snmpPDUs) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(int), err
		}
	}

	return 0, err
}
