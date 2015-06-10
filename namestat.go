package sw

func SysName(ip, community string, timeout int64) (string, error) {
	oid := "1.3.6.1.2.1.1.5.0"
	method := "get"

	snmpPDUs, err := RunSnmp(ip, community, oid, method, timeout)

	if err == nil {
		for _, pdu := range snmpPDUs {
			return pdu.Value.(string), err
		}
	}

	return "", err
}
