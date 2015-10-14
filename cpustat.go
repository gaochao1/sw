package sw

import (
	"github.com/gaochao1/gosnmp"
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
	case "Cisco_IOS_XE","Cisco_IOS_XR":
		oid = "1.3.6.1.4.1.9.9.109.1.1.1.1.7"
		method = "getnext"
	case "Cisco_ASA":
		oid = "1.3.6.1.4.1.9.9.109.1.1.1.1.4"
		return getCiscoASAcpu(ip,community,oid,timeout,retry)
	case "Huawei":
		oid = "1.3.6.1.4.1.2011.5.25.31.1.1.1.1.5"
		return getH3CHWcpumem(ip, community, oid, timeout, retry)
	case "H3C", "H3C_V5", "H3C_V7":
		oid = "1.3.6.1.4.1.25506.2.6.1.1.1.1.6"
		return getH3CHWcpumem(ip, community, oid, timeout, retry)
	case "Ruijie":
		oid = "1.3.6.1.4.1.4881.1.1.10.2.36.1.1.2.0"
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

func getCiscoASAcpu(ip,community,oid string,timeout,retry int) (value int,err error){
	var snmpPDUs []gosnmp.SnmpPDU
	method := "walk"
	for i := 0; i < retry; i++ {
		snmpPDUs, err = RunSnmp(ip, community, oid, method, timeout)
		if len(snmpPDUs) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	var cpu_Values []int
	if err == nil {
		for _, pdu := range snmpPDUs {
			cpu_Values = append(cpu_Values,pdu.Value.(int))
		}
	}
	CPU_Value_SUM := 0
	for _, value := range cpu_Values{
		CPU_Value_SUM = CPU_Value_SUM + value
	}
	
	return int(CPU_Value_SUM/len(cpu_Values)), err	
}


func getH3CHWcpumem(ip, community, oid string, timeout, retry int) (value int, err error) {
	method := "walk"

	var snmpPDUs []gosnmp.SnmpPDU

	for i := 0; i < retry; i++ {
		snmpPDUs, err = RunSnmp(ip, community, oid, method, timeout)
		if len(snmpPDUs) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	for _, v := range snmpPDUs {
		if v.Value.(int) != 0 {
			value = v.Value.(int)
			break
		}
	}

	return value, err
}

func getRuijiecpumem(ip, community, oid string, timeout, retry int) (value int, err error) {
	method := "get"

	var snmpPDUs []gosnmp.SnmpPDU

	for i := 0; i < retry; i++ {
		snmpPDUs, err = RunSnmp(ip, community, oid, method, timeout)
		if len(snmpPDUs) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	return snmpPDUs[0].Value.(int),err
}