package sw

import (
	"fmt"
	"github.com/gaochao1/gosnmp"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	ifNameOid                    = "1.3.6.1.2.1.31.1.1.1.1"
	ifNameOidPrefix              = ".1.3.6.1.2.1.31.1.1.1.1."
	ifHCInOid                    = "1.3.6.1.2.1.31.1.1.1.6"
	ifHCInOidPrefix              = ".1.3.6.1.2.1.31.1.1.1.6."
	ifHCOutOid                   = "1.3.6.1.2.1.31.1.1.1.10"
	ifHCInPktsOid                = "1.3.6.1.2.1.31.1.1.1.7"
	ifHCInPktsOidPrefix          = ".1.3.6.1.2.1.31.1.1.1.7."
	ifHCOutPktsOid               = "1.3.6.1.2.1.31.1.1.1.11"
	ifOperStatusOid              = "1.3.6.1.2.1.2.2.1.8"
	ifOperStatusOidPrefix        = ".1.3.6.1.2.1.2.2.1.8."
	ifHCInBroadcastPktsOid       = "1.3.6.1.2.1.31.1.1.1.9"
	ifHCInBroadcastPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.9."
	ifHCOutBroadcastPktsOid      = "1.3.6.1.2.1.31.1.1.1.13"
	ifHCInMulticastPktsOid       = "1.3.6.1.2.1.31.1.1.1.8"
	ifHCInMulticastPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.8."
	ifHCOutMulticastPktsOid      = "1.3.6.1.2.1.31.1.1.1.12"
)

type IfStats struct {
	IfName               string
	IfIndex              int
	IfHCInOctets         uint64
	IfHCOutOctets        uint64
	IfHCInUcastPkts      uint64
	IfHCOutUcastPkts     uint64
	IfHCInBroadcastPkts  uint64
	IfHCOutBroadcastPkts uint64
	IfHCInMulticastPkts  uint64
	IfHCOutMulticastPkts uint64
	IfOperStatus         int
	TS                   int64
}

func (this *IfStats) String() string {
	return fmt.Sprintf("<IfName:%s, IfIndex:%d, IfHCInOctets:%d, IfHCOutOctets:%d>", this.IfName, this.IfIndex, this.IfHCInOctets, this.IfHCOutOctets)
}

func ListIfStats(ip, community string, timeout int, ignoreIface []string, retry int, ignorePkt bool, ignoreOperStatus bool, ignoreBroadcastPkt bool, ignoreMulticastPkt bool) ([]IfStats, error) {
	var ifStatsList []IfStats

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListIfStats", r)
		}
	}()

	chIfInList := make(chan []gosnmp.SnmpPDU)
	chIfOutList := make(chan []gosnmp.SnmpPDU)

	chIfNameList := make(chan []gosnmp.SnmpPDU)
	chIfStatusList := make(chan []gosnmp.SnmpPDU)

	go ListIfHCInOctets(ip, community, timeout, chIfInList, retry)
	go ListIfHCOutOctets(ip, community, timeout, chIfOutList, retry)

	go ListIfName(ip, community, timeout, chIfNameList, retry)

	ifInList := <-chIfInList
	ifOutList := <-chIfOutList

	ifNameList := <-chIfNameList

	chIfInPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutPktList := make(chan []gosnmp.SnmpPDU)

	var ifInPktList, ifOutPktList []gosnmp.SnmpPDU

	if ignorePkt == false {
		go ListIfHCInUcastPkts(ip, community, timeout, chIfInPktList, retry)
		go ListIfHCOutUcastPkts(ip, community, timeout, chIfOutPktList, retry)
		ifInPktList = <-chIfInPktList
		ifOutPktList = <-chIfOutPktList
	}

	chIfInBroadcastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutBroadcastPktList := make(chan []gosnmp.SnmpPDU)

	var ifInBroadcastPktList, ifOutBroadcastPktList []gosnmp.SnmpPDU

	if ignoreBroadcastPkt == false {
		go ListIfHCInBroadcastPkts(ip, community, timeout, chIfInBroadcastPktList, retry)
		go ListIfHCOutBroadcastPkts(ip, community, timeout, chIfOutBroadcastPktList, retry)
		ifInBroadcastPktList = <-chIfInBroadcastPktList
		ifOutBroadcastPktList = <-chIfOutBroadcastPktList
	}

	chIfInMulticastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutMulticastPktList := make(chan []gosnmp.SnmpPDU)

	var ifInMulticastPktList, ifOutMulticastPktList []gosnmp.SnmpPDU

	if ignoreMulticastPkt == false {
		go ListIfHCInMulticastPkts(ip, community, timeout, chIfInMulticastPktList, retry)
		go ListIfHCOutMulticastPkts(ip, community, timeout, chIfOutMulticastPktList, retry)
		ifInMulticastPktList = <-chIfInMulticastPktList
		ifOutMulticastPktList = <-chIfOutMulticastPktList
	}

	var ifStatusList []gosnmp.SnmpPDU
	if ignoreOperStatus == false {
		go ListIfOperStatus(ip, community, timeout, chIfStatusList, retry)
		ifStatusList = <-chIfStatusList
	}

	if len(ifNameList) > 0 && len(ifInList) > 0 && len(ifOutList) > 0 {
		now := time.Now().Unix()

		for _, ifNamePDU := range ifNameList {

			ifName := ifNamePDU.Value.(string)

			check := true
			if len(ignoreIface) > 0 {
				for _, ignore := range ignoreIface {
					if strings.Contains(ifName, ignore) {
						check = false
						break
					}
				}
			}

			if check {
				var ifStats IfStats

				ifIndexStr := strings.Replace(ifNamePDU.Name, ifNameOidPrefix, "", 1)

				ifStats.IfIndex, _ = strconv.Atoi(ifIndexStr)

				for ti, ifHCInOctetsPDU := range ifInList {
					if strings.Replace(ifHCInOctetsPDU.Name, ifHCInOidPrefix, "", 1) == ifIndexStr {

						ifStats.IfHCInOctets = ifInList[ti].Value.(uint64)
						ifStats.IfHCOutOctets = ifOutList[ti].Value.(uint64)
					}
					if ignorePkt == false {
						for ti, ifHCInPktsPDU := range ifInPktList {
							if strings.Replace(ifHCInPktsPDU.Name, ifHCInPktsOidPrefix, "", 1) == ifIndexStr {
								ifStats.IfHCInUcastPkts = ifInPktList[ti].Value.(uint64)
								ifStats.IfHCOutUcastPkts = ifOutPktList[ti].Value.(uint64)
							}
						}
					}
					if ignoreBroadcastPkt == false {
						for ti, ifHCInBroadcastPktPDU := range ifInBroadcastPktList {
							if strings.Replace(ifHCInBroadcastPktPDU.Name, ifHCInBroadcastPktsOidPrefix, "", 1) == ifIndexStr {
								ifStats.IfHCInBroadcastPkts = ifInBroadcastPktList[ti].Value.(uint64)
								ifStats.IfHCOutBroadcastPkts = ifOutBroadcastPktList[ti].Value.(uint64)
							}
						}
					}
					if ignoreMulticastPkt == false {
						for ti, ifHCInMulticastPktPDU := range ifInMulticastPktList {
							if strings.Replace(ifHCInMulticastPktPDU.Name, ifHCInMulticastPktsOidPrefix, "", 1) == ifIndexStr {
								ifStats.IfHCInMulticastPkts = ifInMulticastPktList[ti].Value.(uint64)
								ifStats.IfHCOutMulticastPkts = ifOutMulticastPktList[ti].Value.(uint64)
							}
						}
					}
					if ignoreOperStatus == false {
						for ti, ifOperStatusPDU := range ifStatusList {
							if strings.Replace(ifOperStatusPDU.Name, ifOperStatusOidPrefix, "", 1) == ifIndexStr {
								ifStats.IfOperStatus = ifStatusList[ti].Value.(int)
							}
						}
					}
					ifStats.TS = now
					ifStats.IfName = ifName

				}

				ifStatsList = append(ifStatsList, ifStats)

			}
		}
	}

	return ifStatsList, nil
}

func ListIfOperStatus(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifOperStatusOid)
}

func ListIfName(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifNameOid)
}

func ListIfHCInOctets(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCInOid)
}

func ListIfHCOutOctets(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCOutOid)
}

func ListIfHCInUcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCInPktsOid)
}

func ListIfHCInBroadcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCInBroadcastPktsOid)
}

func ListIfHCOutBroadcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCInBroadcastPktsOid)
}

func ListIfHCInMulticastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCInMulticastPktsOid)
}

func ListIfHCOutMulticastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCOutMulticastPktsOid)
}

func ListIfHCOutUcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCOutPktsOid)
}

func RunSnmpRetry(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, oid string) {
	method := "walk"
	var snmpPDUs []gosnmp.SnmpPDU
	for i := 0; i < retry; i++ {
		snmpPDUs, _ = RunSnmp(ip, community, oid, method, timeout)
		fmt.Println(snmpPDUs)
		if len(snmpPDUs) > 0 {
			ch <- snmpPDUs
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	ch <- snmpPDUs
	return
}
