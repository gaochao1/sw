package sw

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gaochao1/gosnmp"
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
	// multicastpkt
	ifHCInMulticastPktsOid       = "1.3.6.1.2.1.31.1.1.1.8"
	ifHCInMulticastPktsOidPrefix = ".1.3.6.1.2.1.31.1.1.1.8."
	ifHCOutMulticastPktsOid      = "1.3.6.1.2.1.31.1.1.1.12"
	// speed 配置
	ifSpeedOid       = "1.3.6.1.2.1.31.1.1.1.15"
	ifSpeedOidPrefix = ".1.3.6.1.2.1.31.1.1.1.15."

	// Discards配置
	ifInDiscardsOid       = "1.3.6.1.2.1.2.2.1.13"
	ifInDiscardsOidPrefix = ".1.3.6.1.2.1.2.2.1.13."
	ifOutDiscardsOid      = "1.3.6.1.2.1.2.2.1.19"

	// Errors配置
	ifInErrorsOid        = "1.3.6.1.2.1.2.2.1.14"
	ifInErrorsOidPrefix  = ".1.3.6.1.2.1.2.2.1.14."
	ifOutErrorsOid       = "1.3.6.1.2.1.2.2.1.20"
	ifOutErrorsOidPrefix = ".1.3.6.1.2.1.2.2.1.20."

	//ifInUnknownProtos 由于未知或不支持的网络协议而丢弃的输入报文的数量
	ifInUnknownProtosOid    = "1.3.6.1.2.1.2.2.1.15"
	ifInUnknownProtosPrefix = ".1.3.6.1.2.1.2.2.1.15."

	//ifOutQLen 接口上输出报文队列长度
	ifOutQLenOid    = "1.3.6.1.2.1.2.2.1.21"
	ifOutQLenPrefix = ".1.3.6.1.2.1.2.2.1.21."
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
	IfSpeed              int
	IfInDiscards         int
	IfOutDiscards        int
	IfInErrors           int
	IfOutErrors          int
	IfInUnknownProtos    int
	IfOutQLen            int
	IfOperStatus         int
	TS                   int64
}

func (this *IfStats) String() string {
	return fmt.Sprintf("<IfName:%s, IfIndex:%d, IfHCInOctets:%d, IfHCOutOctets:%d>", this.IfName, this.IfIndex, this.IfHCInOctets, this.IfHCOutOctets)
}

func ListIfStats(ip, community string, timeout int, ignoreIface []string, retry int, ignorePkt bool, ignoreOperStatus bool, ignoreBroadcastPkt bool, ignoreMulticastPkt bool, ignoreDiscards bool, ignoreErrors bool, ignoreUnknownProtos bool, ignoreOutQLen bool) ([]IfStats, error) {
	var ifStatsList []IfStats

	defer func() {
		if r := recover(); r != nil {
			log.Println(ip+" Recovered in ListIfStats", r)
		}
	}()

	chIfInList := make(chan []gosnmp.SnmpPDU)
	chIfOutList := make(chan []gosnmp.SnmpPDU)

	chIfNameList := make(chan []gosnmp.SnmpPDU)

	go ListIfHCInOctets(ip, community, timeout, chIfInList, retry)
	time.Sleep(100 * time.Millisecond)

	go ListIfHCOutOctets(ip, community, timeout, chIfOutList, retry)
	time.Sleep(100 * time.Millisecond)

	go ListIfName(ip, community, timeout, chIfNameList, retry)
	time.Sleep(100 * time.Millisecond)

	ifInList := <-chIfInList
	ifOutList := <-chIfOutList

	ifNameList := <-chIfNameList

	chIfInPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutPktList := make(chan []gosnmp.SnmpPDU)

	var ifInPktList, ifOutPktList []gosnmp.SnmpPDU

	if ignorePkt == false {
		go ListIfHCInUcastPkts(ip, community, timeout, chIfInPktList, retry)
		time.Sleep(100 * time.Millisecond)
		go ListIfHCOutUcastPkts(ip, community, timeout, chIfOutPktList, retry)
		time.Sleep(100 * time.Millisecond)
		ifInPktList = <-chIfInPktList
		ifOutPktList = <-chIfOutPktList
	}

	chIfInBroadcastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutBroadcastPktList := make(chan []gosnmp.SnmpPDU)

	var ifInBroadcastPktList, ifOutBroadcastPktList []gosnmp.SnmpPDU

	if ignoreBroadcastPkt == false {
		go ListIfHCInBroadcastPkts(ip, community, timeout, chIfInBroadcastPktList, retry)
		time.Sleep(100 * time.Millisecond)
		go ListIfHCOutBroadcastPkts(ip, community, timeout, chIfOutBroadcastPktList, retry)
		time.Sleep(100 * time.Millisecond)
		ifInBroadcastPktList = <-chIfInBroadcastPktList
		ifOutBroadcastPktList = <-chIfOutBroadcastPktList
	}

	chIfInMulticastPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutMulticastPktList := make(chan []gosnmp.SnmpPDU)

	var ifInMulticastPktList, ifOutMulticastPktList []gosnmp.SnmpPDU

	if ignoreMulticastPkt == false {
		go ListIfHCInMulticastPkts(ip, community, timeout, chIfInMulticastPktList, retry)
		time.Sleep(100 * time.Millisecond)
		go ListIfHCOutMulticastPkts(ip, community, timeout, chIfOutMulticastPktList, retry)
		time.Sleep(100 * time.Millisecond)
		ifInMulticastPktList = <-chIfInMulticastPktList
		ifOutMulticastPktList = <-chIfOutMulticastPktList
	}

	//Discards
	chIfInDiscardsList := make(chan []gosnmp.SnmpPDU)
	chIfOutDiscardsList := make(chan []gosnmp.SnmpPDU)

	var ifInDiscardsList, ifOutDiscardsList []gosnmp.SnmpPDU

	if ignoreDiscards == false {
		go ListIfInDiscards(ip, community, timeout, chIfInDiscardsList, retry)
		time.Sleep(100 * time.Millisecond)
		go ListIfOutDiscards(ip, community, timeout, chIfOutDiscardsList, retry)
		time.Sleep(100 * time.Millisecond)
		ifInDiscardsList = <-chIfInDiscardsList
		ifOutDiscardsList = <-chIfOutDiscardsList
	}

	//Errors
	chIfInErrorsList := make(chan []gosnmp.SnmpPDU)
	chIfOutErrorsList := make(chan []gosnmp.SnmpPDU)

	var ifInErrorsList, ifOutErrorsList []gosnmp.SnmpPDU

	if ignoreErrors == false {
		go ListIfInErrors(ip, community, timeout, chIfInErrorsList, retry)
		time.Sleep(100 * time.Millisecond)
		go ListIfOutErrors(ip, community, timeout, chIfOutErrorsList, retry)
		time.Sleep(100 * time.Millisecond)
		ifInErrorsList = <-chIfInErrorsList
		ifOutErrorsList = <-chIfOutErrorsList
	}

	//UnknownProtos
	chIfInUnknownProtosList := make(chan []gosnmp.SnmpPDU)

	var ifInUnknownProtosList []gosnmp.SnmpPDU

	if ignoreUnknownProtos == false {
		go ListIfInUnknownProtos(ip, community, timeout, chIfInUnknownProtosList, retry)
		time.Sleep(100 * time.Millisecond)
		ifInUnknownProtosList = <-chIfInUnknownProtosList
	}
	//QLen
	chIfOutQLenList := make(chan []gosnmp.SnmpPDU)

	var ifOutQLenList []gosnmp.SnmpPDU

	if ignoreOutQLen == false {
		go ListIfOutQLen(ip, community, timeout, chIfOutQLenList, retry)
		time.Sleep(100 * time.Millisecond)
		ifOutQLenList = <-chIfOutQLenList
	}
	// OperStatus
	var ifStatusList []gosnmp.SnmpPDU
	chIfStatusList := make(chan []gosnmp.SnmpPDU)
	if ignoreOperStatus == false {
		go ListIfOperStatus(ip, community, timeout, chIfStatusList, retry)
		time.Sleep(100 * time.Millisecond)
		ifStatusList = <-chIfStatusList
	}

	// Speed
	var ifSpeedList []gosnmp.SnmpPDU
	chIfSpeedList := make(chan []gosnmp.SnmpPDU)

	go ListIfSpeed(ip, community, timeout, chIfSpeedList, retry)
	ifSpeedList = <-chIfSpeedList

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
						break
					}
				}
				if ignorePkt == false {
					for ti, ifHCInPktsPDU := range ifInPktList {
						if strings.Replace(ifHCInPktsPDU.Name, ifHCInPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInUcastPkts = ifInPktList[ti].Value.(uint64)
							ifStats.IfHCOutUcastPkts = ifOutPktList[ti].Value.(uint64)
							break
						}
					}
				}
				if ignoreBroadcastPkt == false {
					for ti, ifHCInBroadcastPktPDU := range ifInBroadcastPktList {
						if strings.Replace(ifHCInBroadcastPktPDU.Name, ifHCInBroadcastPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInBroadcastPkts = ifInBroadcastPktList[ti].Value.(uint64)
							ifStats.IfHCOutBroadcastPkts = ifOutBroadcastPktList[ti].Value.(uint64)
							break
						}
					}
				}
				if ignoreMulticastPkt == false {
					for ti, ifHCInMulticastPktPDU := range ifInMulticastPktList {
						if strings.Replace(ifHCInMulticastPktPDU.Name, ifHCInMulticastPktsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfHCInMulticastPkts = ifInMulticastPktList[ti].Value.(uint64)
							ifStats.IfHCOutMulticastPkts = ifOutMulticastPktList[ti].Value.(uint64)
							break
						}
					}
				}

				if ignoreDiscards == false {
					for ti, ifInDiscardsPDU := range ifInDiscardsList {
						if strings.Replace(ifInDiscardsPDU.Name, ifInDiscardsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfInDiscards = ifInDiscardsList[ti].Value.(int)
							ifStats.IfOutDiscards = ifOutDiscardsList[ti].Value.(int)
							break
						}
					}
				}

				if ignoreErrors == false {
					for ti, ifInErrorsPDU := range ifInErrorsList {
						if strings.Replace(ifInErrorsPDU.Name, ifInErrorsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfInErrors = ifInErrorsList[ti].Value.(int)
							break
						}
					}
					for ti, ifOutErrorsPDU := range ifOutErrorsList {
						if strings.Replace(ifOutErrorsPDU.Name, ifOutErrorsOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfOutErrors = ifOutErrorsList[ti].Value.(int)
							break
						}
					}
				}

				if ignoreOperStatus == false {
					for ti, ifOperStatusPDU := range ifStatusList {
						if strings.Replace(ifOperStatusPDU.Name, ifOperStatusOidPrefix, "", 1) == ifIndexStr {
							ifStats.IfOperStatus = ifStatusList[ti].Value.(int)
							break
						}
					}
				}

				if ignoreUnknownProtos == false {
					for ti, ifInUnknownProtosPDU := range ifInUnknownProtosList {
						if strings.Replace(ifInUnknownProtosPDU.Name, ifInUnknownProtosPrefix, "", 1) == ifIndexStr {
							ifStats.IfInUnknownProtos = ifInUnknownProtosList[ti].Value.(int)
							break
						}
					}
				}

				if ignoreOutQLen == false {
					for ti, ifOutQLenPDU := range ifOutQLenList {
						if strings.Replace(ifOutQLenPDU.Name, ifOutQLenPrefix, "", 1) == ifIndexStr {
							ifStats.IfOutQLen = ifOutQLenList[ti].Value.(int)
							break
						}
					}
				}

				for ti, ifSpeedPDU := range ifSpeedList {
					if strings.Replace(ifSpeedPDU.Name, ifSpeedOidPrefix, "", 1) == ifIndexStr {
						ifStats.IfSpeed = 1000 * 1000 * ifSpeedList[ti].Value.(int)
						break
					}
				}

				ifStats.TS = now
				ifStats.IfName = ifName
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
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCOutBroadcastPktsOid)
}

func ListIfHCInMulticastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCInMulticastPktsOid)
}

func ListIfHCOutMulticastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCOutMulticastPktsOid)
}

func ListIfInDiscards(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifInDiscardsOid)
}

func ListIfOutDiscards(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifOutDiscardsOid)
}

func ListIfInErrors(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifInErrorsOid)
}

func ListIfOutErrors(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifOutErrorsOid)
}

func ListIfHCOutUcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifHCOutPktsOid)
}

func ListIfInUnknownProtos(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifInUnknownProtosOid)
}

func ListIfOutQLen(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifOutQLenOid)
}

func ListIfSpeed(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int) {
	RunSnmpRetry(ip, community, timeout, ch, retry, ifSpeedOid)
}

func RunSnmpRetry(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU, retry int, oid string) {

	var snmpPDUs []gosnmp.SnmpPDU
	var err error
	snmpPDUs, err = RunSnmpwalk(ip, community, oid, retry, timeout)

	if err != nil {
		log.Println(ip, oid, err)
		return
	}
	ch <- snmpPDUs
	return
}
