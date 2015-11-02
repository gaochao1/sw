package sw

import (
	"github.com/gaochao1/gosnmp"
	"log"
	"time"
	"strconv"
)

func MemUtilization(ip, community string, timeout, retry int) (int, error) {
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
	case "Cisco", "Cisco_IOS_XE","Cisco_IOS_7200":
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
	case "Cisco_IOS_XR":
		return getCisco_IOS_XR_Mem(ip, community, timeout, retry)
	case "Cisco_ASA":
		return getCisco_ASA_Mem(ip, community, timeout, retry)
	case "Huawei":
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.7"
		return getH3CHWcpumem(ip, community, oid, timeout, retry)
	case "H3C", "H3C_V5", "H3C_V7":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.8"
		return getH3CHWcpumem(ip, community, oid, timeout, retry)
	case "Ruijie":
		oid = "1.3.6.1.4.1.4881.1.1.10.2.35.1.1.1.3.0"
		return getRuijiecpumem(ip, community, oid, timeout, retry)
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
func getCisco_IOS_XR_Mem(ip, community string, timeout, retry int)(int,error){
	cpuindex := "1.3.6.1.4.1.9.9.109.1.1.1.1.2" 
	method := "getnext"
	var snmpPDUs []gosnmp.SnmpPDU
	var err error
	var index string
	for i := 0; i < retry; i++ {
		snmpPDUs, err = RunSnmp(ip, community, cpuindex, method, timeout)
		if len(snmpPDUs) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	index = strconv.Itoa(snmpPDUs[0].Value.(int))
	method = "get"
	memUsedOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.18."+index+".1"
	snmpMemUsed, _ := RunSnmp(ip, community, memUsedOid, method, timeout)
	memFreeOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.20."+index+".1"
	snmpMemFree, _ := RunSnmp(ip, community, memFreeOid, method, timeout)
	if &snmpMemFree[0] != nil && &snmpMemUsed[0] != nil {
		memUsed := snmpMemUsed[0].Value.(uint64)
		memFree := snmpMemFree[0].Value.(uint64)
		if memUsed+memFree != 0 {
			memUtili := float64(memUsed) / float64(memUsed+memFree)
			return int(memUtili * 100), err
		}
	}
	return 0, err
}

func getCisco_ASA_Mem(ip, community string, timeout, retry int)(int,error){
	method := "walk"
	memUsedOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.7"
	snmpMemUsed, err := RunSnmp(ip, community, memUsedOid, method, timeout)

	memFreeOid := "1.3.6.1.4.1.9.9.221.1.1.1.1.8"
	snmpMemFree, err := RunSnmp(ip, community, memFreeOid, method, timeout)
	if &snmpMemFree[0] != nil && &snmpMemUsed[0] != nil {
		memUsed := snmpMemUsed[0].Value.(int)
		memFree := snmpMemFree[0].Value.(int)

		if memUsed+memFree != 0 {
			memUtili := float64(memUsed) / float64(memUsed+memFree)
			return int(memUtili * 100), nil
		}
	}
	return 0, err
}
