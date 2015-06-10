package sw

import (
	"github.com/alouca/gosnmp"
	"strings"
)

func RunSnmp(ip, community, oid, method string, timeout int64) (snmpPDUs []gosnmp.SnmpPDU, err error) {
	cur_gosnmp, err := gosnmp.NewGoSNMP(ip, community, gosnmp.Version2c, timeout)

	if err != nil {
		return nil, err
	} else {
		cur_gosnmp.SetTimeout(timeout)
		snmpPDUs, err := ParseSnmpMethod(oid, method, cur_gosnmp)
		if err != nil {
			return nil, err
		} else {
			return snmpPDUs, err
		}
	}

	return
}

func ParseSnmpMethod(oid, method string, cur_gosnmp *gosnmp.GoSNMP) (snmpPDUs []gosnmp.SnmpPDU, err error) {
	var snmpPacket *gosnmp.SnmpPacket

	switch method {
	case "get":
		snmpPacket, err = cur_gosnmp.Get(oid)
		if err != nil {
			return nil, err
		} else {
			snmpPDUs = snmpPacket.Variables
			return snmpPDUs, err
		}
	case "getnext":
		snmpPacket, err = cur_gosnmp.GetNext(oid)
		if err != nil {
			return nil, err
		} else {
			snmpPDUs = snmpPacket.Variables
			return snmpPDUs, err
		}
	default:
		snmpPDUs, err = cur_gosnmp.Walk(oid)
		return snmpPDUs, err
	}

	return
}

func snmpPDUNameToIfIndex(snmpPDUName string) string {
	oidSplit := strings.Split(snmpPDUName, ".")
	curIfIndex := oidSplit[len(oidSplit)-1]
	return curIfIndex
}
