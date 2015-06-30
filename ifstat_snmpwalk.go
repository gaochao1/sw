package sw

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func ListIfStatsSnmpWalk(ip, community string, timeout int, ignoreIface []string, retry int, ignorePkt bool) ([]IfStats, error) {
	var ifStatsList []IfStats

	chIfInMap := make(chan map[string]string)
	chIfOutMap := make(chan map[string]string)

	chIfNameMap := make(chan map[string]string)

	go WalkIfIn(ip, community, timeout, chIfInMap, retry)
	go WalkIfOut(ip, community, timeout, chIfOutMap, retry)

	go WalkIfName(ip, community, timeout, chIfNameMap, retry)

	ifInMap := <-chIfInMap
	ifOutMap := <-chIfOutMap

	ifNameMap := <-chIfNameMap

	chIfInPktMap := make(chan map[string]string)
	chIfOutPktMap := make(chan map[string]string)

	var ifInPktMap, ifOutPktMap map[string]string

	if ignorePkt == false {
		go WalkIfInPkts(ip, community, timeout, chIfInPktMap, retry)
		go WalkIfOutPkts(ip, community, timeout, chIfOutPktMap, retry)
		ifInPktMap = <-chIfInPktMap
		ifOutPktMap = <-chIfOutPktMap
	}

	if len(ifNameMap) > 0 && len(ifInMap) > 0 && len(ifOutMap) > 0 {

		now := time.Now().Unix()

		for ifIndex, ifName := range ifNameMap {

			check := true
			if len(ignoreIface) > 0 {
				for _, ignore := range ignoreIface {
					if strings.Contains(ifName, ignore) {
						check = false
						break
					}
				}
			}

			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in ListIfStats_SnmpWalk", r)
				}
			}()

			if check {
				var ifStats IfStats

				ifStats.IfIndex, _ = strconv.Atoi(ifIndex)

				ifStats.IfHCInOctets, _ = strconv.ParseUint(ifInMap[ifIndex], 10, 64)
				ifStats.IfHCOutOctets, _ = strconv.ParseUint(ifOutMap[ifIndex], 10, 64)

				if ignorePkt == false {
					ifStats.IfHCInUcastPkts, _ = strconv.ParseUint(ifInPktMap[ifIndex], 10, 64)
					ifStats.IfHCOutUcastPkts, _ = strconv.ParseUint(ifOutPktMap[ifIndex], 10, 64)
				}

				ifStats.TS = now
				ifStats.IfName = ifName

				ifStatsList = append(ifStatsList, ifStats)

			}
		}
	}

	return ifStatsList, nil
}

func WalkIfName(ip, community string, timeout int, ch chan map[string]string, retry int) {
	WalkIf(ip, ifNameOid, community, timeout, retry, ch)
}

func WalkIfIn(ip, community string, timeout int, ch chan map[string]string, retry int) {
	WalkIf(ip, ifHCInOid, community, timeout, retry, ch)
}

func WalkIfOut(ip, community string, timeout int, ch chan map[string]string, retry int) {
	WalkIf(ip, ifHCOutOid, community, timeout, retry, ch)
}

func WalkIfInPkts(ip, community string, timeout int, ch chan map[string]string, retry int) {
	WalkIf(ip, ifHCInPktsOid, community, timeout, retry, ch)
}

func WalkIfOutPkts(ip, community string, timeout int, ch chan map[string]string, retry int) {
	WalkIf(ip, ifHCOutPktsOid, community, timeout, retry, ch)
}

func WalkIf(ip, oid, community string, timeout, retry int, ch chan map[string]string) {
	result := make(map[string]string)

	for i := 0; i < retry; i++ {
		out, err := cmdTimeout(timeout, "snmpwalk", "-v", "2c", "-c", community, ip, oid)
		if err != nil {
			log.Println(ip, oid, err)
		}

		list := strings.Split(out, "IF-MIB")
		for _, v := range list {

			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in WalkIf", r)
				}
			}()

			if len(v) > 0 && strings.Contains(v, "=") {
				vt := strings.Split(v, "=")

				var ifIndex, ifName string
				if strings.Contains(vt[0], ".") {
					ifIndex = strings.Split(vt[0], ".")[1]
					ifIndex = strings.TrimSpace(ifIndex)

				}

				if strings.Contains(vt[1], ":") {
					ifName = strings.Split(vt[1], ":")[1]
					ifName = strings.TrimSpace(ifName)
				}

				result[ifIndex] = ifName
			}
		}

		if len(result) > 0 {
			ch <- result
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	ch <- result
	return
}

func cmdTimeout(timeout int, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Start()
	timer := time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
		err := cmd.Process.Kill()
		if err != nil {
			log.Println("failed to kill: ", err)
		}
	})
	err := cmd.Wait()
	timer.Stop()

	return out.String(), err
}
