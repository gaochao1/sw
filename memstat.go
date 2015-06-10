package sw

import (
	"log"
)

func MemUtilization(ip, community string, timeout int64) (int, error) {
	vendor, err := SysVendor(ip, community, timeout)
	method := "get"
	var oid string

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in MemUtilization", r)
		}
	}()

	switch vendor {
	case "Cisco_NX":
		oid = "1.3.6.1.4.1.9.9.305.1.1.2.0"
	case "Cisco":
		memUsedOid := "1.3.6.1.4.1.9.9.48.1.1.1.5.1"
		snmpMemUsed, _ := RunSnmp(ip, community, memUsedOid, method, timeout)

		memFreeOid := "1.3.6.1.4.1.9.9.48.1.1.1.6.1"
		snmpMemFree, _ := RunSnmp(ip, community, memFreeOid, method, timeout)

		if &snmpMemFree[0] != nil && &snmpMemUsed[0] != nil {
			memUsed := snmpMemUsed[0].Value.(int)
			memFree := snmpMemFree[0].Value.(int)

			if memUsed+memFree != 0 {
				memUtili := float64(memUsed) / float64(memUsed+memFree)
				return int(memUtili * 100), nil
			}
		}
	case "Huawei":
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.7.16842753"
	case "H3C":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.8.74"
	case "H3C_V5":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.8.74"
	case "H3C_V7":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.8.212"
	default:
		return 0, err
	}

	snmpPDUs, err := RunSnmp(ip, community, oid, method, timeout)

	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(int), err
		}
	}

	return 0, err
}
