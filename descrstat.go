package sw

import (
	"strings"
)

func SysDescr(ip, community string, timeout int) (string, error) {
	oid := "1.3.6.1.2.1.1.1.0"
	method := "get"

	snmpPDUs, err := RunSnmp(ip, community, oid, method, timeout)

	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(string), err
		}
	}

	return "", err
}

func SysVendor(ip, community string, timeout int) (string, error) {
	sysDescr, err := SysDescr(ip, community, timeout)
	sysDescrLower := strings.ToLower(sysDescr)

	if strings.Contains(sysDescrLower, "cisco nx-os") {
		return "Cisco_NX", err
	}

	if strings.Contains(sysDescrLower, "cisco ios") {
		return "Cisco", err
	}

	if strings.Contains(sysDescrLower, "h3c") {
		if strings.Contains(sysDescr, "Software Version 5") {
			return "H3C_V5", err
		}

		if strings.Contains(sysDescr, "Software Version 7") {
			return "H3C_V7", err

		}

		return "H3C", err
	}

	if strings.Contains(sysDescrLower, "huawei") {
		return "Huawei", err
	}

	if strings.Contains(sysDescrLower, "linux") {
		return "Linux", err
	}

	return "", err
}
