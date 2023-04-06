// Copyright (c) 2022 Intel Corporation.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build dpdk

package p4

import (
	"context"
	"fmt"
	"github.com/antoninbas/p4runtime-go-client/pkg/client"
	"github.com/ipdk-io/k8s-infra-offload/pkg/inframanager/store"
	"strconv"
	"strings"
	//	p4_v1 "github.com/p4lang/p4runtime/go/p4/v1"
	log "github.com/sirupsen/logrus"
	//	"net"
)

type tabletype int

const (
	denyall tabletype = iota
	policyadd
	policydel
	policyupdate
	workloadadd
	workloadupdate
	workloaddel
)

func AclPodIpProtoTableEgress(ctx context.Context, p4RtC *client.Client,
	protocol string, workerep string, polID uint16, rangeID uint16,
	action InterfaceType) error {
	switch action {
	case Insert:
		if protocol != "" {
			entryAdd := p4RtC.NewTableEntry(
				"k8s_dp_control.acl_pod_ip_proto_table_egress",
				map[string]client.MatchInterface{
					"hdr.ipv4.src_addr": &client.ExactMatch{
						Value: Pack32BinaryIP4(workerep),
					},
					"hdr.ipv4.protocol": &client.LpmMatch{
						Value: valueToBytesStr(protocol),
						PLen:  16,
					},
				},
				p4RtC.NewTableActionDirect("k8s_dp_control.set_range_check_ref",
					[][]byte{valueToBytes16(polID),
						valueToBytes16(rangeID)}),
				nil,
			)
			if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
				log.Errorf("Cannot insert entry into 'acl_pod_ip_proto_table_egress': %v", err)
				return err
			}
		} else {
			entryAdd := p4RtC.NewTableEntry(
				"k8s_dp_control.acl_pod_ip_proto_table_egress",
				map[string]client.MatchInterface{
					"hdr.ipv4.src_addr": &client.ExactMatch{
						Value: Pack32BinaryIP4(workerep),
					},
					"hdr.ipv4.protocol": &client.LpmMatch{
						Value: valueToBytesStr(protocol),
						PLen:  16,
					},
				},
				p4RtC.NewTableActionDirect("k8s_dp_control.set_status_match_ipset_only",
					[][]byte{valueToBytes16(polID)}),
				nil,
			)
			if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
				log.Errorf("Cannot insert entry into 'acl_pod_ip_proto_table_egress': %v", err)
				return err
			}
		}

	case Delete:
		entryDelete := p4RtC.NewTableEntry(
			"k8s_dp_control.acl_pod_ip_proto_table_egress",
			map[string]client.MatchInterface{
				"hdr.ipv4.src_addr": &client.ExactMatch{
					Value: Pack32BinaryIP4(workerep),
				},
				"hdr.ipv4.protocol": &client.LpmMatch{
					Value: valueToBytesStr(protocol),
					PLen:  16,
				},
			},
			nil,
			nil,
		)
		if err := p4RtC.DeleteTableEntry(ctx, entryDelete); err != nil {
			log.Errorf("Cannot delete entry from 'acl_pod_ip_proto_table_egress': %v", err)
			return err
		}

	default:
		log.Warnf("Invalid action %v", action)
		err := fmt.Errorf("Invalid action %v", action)
		return err
	}

	return nil
}

func AclPodIpProtoTableIngress(ctx context.Context, p4RtC *client.Client,
	protocol string, workerep string, polID uint16, rangeID uint16,
	action InterfaceType) error {
	switch action {
	case Insert:
		if protocol != "" {
			entryAdd := p4RtC.NewTableEntry(
				"k8s_dp_control.acl_pod_ip_proto_table_ingress",
				map[string]client.MatchInterface{
					"hdr.ipv4.dst_addr": &client.ExactMatch{
						Value: Pack32BinaryIP4(workerep),
					},
					"hdr.ipv4.protocol": &client.LpmMatch{
						Value: valueToBytesStr(protocol),
						PLen:  16,
					},
				},
				p4RtC.NewTableActionDirect("k8s_dp_control.set_range_check_ref",
					[][]byte{valueToBytes16(polID),
						valueToBytes16(rangeID)}),
				nil,
			)
			if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
				log.Errorf("Cannot insert entry into 'acl_pod_ip_proto_table_ingress': %v", err)
				return err
			}
		} else {
			entryAdd := p4RtC.NewTableEntry(
				"k8s_dp_control.acl_pod_ip_proto_table_ingress",
				map[string]client.MatchInterface{
					"hdr.ipv4.dst_addr": &client.ExactMatch{
						Value: Pack32BinaryIP4(workerep),
					},
					"hdr.ipv4.protocol": &client.LpmMatch{
						Value: valueToBytesStr(protocol),
						PLen:  16,
					},
				},
				p4RtC.NewTableActionDirect("k8s_dp_control.set_status_match_ipset_only",
					[][]byte{valueToBytes16(polID)}),
				nil,
			)
			if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
				log.Errorf("Cannot insert entry into 'acl_pod_ip_proto_table_ingress': %v", err)
				return err
			}
		}

	case Delete:
		entryDelete := p4RtC.NewTableEntry(
			"k8s_dp_control.acl_pod_ip_proto_table_ingress",
			map[string]client.MatchInterface{
				"hdr.ipv4.dst_addr": &client.ExactMatch{
					Value: Pack32BinaryIP4(workerep),
				},
				"hdr.ipv4.protocol": &client.LpmMatch{
					Value: valueToBytesStr(protocol),
					PLen:  16,
				},
			},
			nil,
			nil,
		)
		if err := p4RtC.DeleteTableEntry(ctx, entryDelete); err != nil {
			log.Errorf("Cannot delete entry from 'acl_pod_ip_proto_table_ingress': %v", err)
			return err
		}

	default:
		log.Warnf("Invalid action %v", action)
		err := fmt.Errorf("Invalid action %v", action)
		return err
	}

	return nil
}

func AclIpSetMatchTableEgress(ctx context.Context, p4RtC *client.Client,
	polID uint16, cidr string, mask uint8, action InterfaceType) error {
	res := strings.Split(cidr, "/")
	ip := res[0]
	plen, _ := strconv.Atoi(res[1])
	switch action {
	case Insert:
		entryAdd := p4RtC.NewTableEntry(
			"k8s_dp_control.acl_ipset_match_table_egress",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
				"hdr.ipv4.dst_addr": &client.LpmMatch{
					Value: Pack32BinaryIP4(ip),
					PLen:  int32(plen),
				},
			},
			p4RtC.NewTableActionDirect("k8s_dp_control.set_ipset_match_result",
				[][]byte{valueToBytes8(mask)}),
			nil,
		)
		if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
			log.Errorf("Cannot insert entry into 'acl_ipset_match_table_egress': %v", err)
			return err
		}

	case Delete:
		entryDelete := p4RtC.NewTableEntry(
			"k8s_dp_control.acl_ipset_match_table_egress",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
				"hdr.ipv4.dst_addr": &client.LpmMatch{
					Value: Pack32BinaryIP4(ip),
					PLen:  int32(plen),
				},
			},
			nil,
			nil,
		)
		if err := p4RtC.DeleteTableEntry(ctx, entryDelete); err != nil {
			log.Errorf("Cannot delete entry from 'acl_ipset_match_table_egress': %v", err)
			return err
		}
	}

	return nil
}

func AclIpSetMatchTableIngress(ctx context.Context, p4RtC *client.Client,
	polID uint16, cidr string, mask uint8, action InterfaceType) error {
	res := strings.Split(cidr, "/")
	ip := res[0]
	plen, _ := strconv.Atoi(res[1])
	switch action {
	case Insert:
		entryAdd := p4RtC.NewTableEntry(
			"k8s_dp_control.acl_ipset_match_table_ingress",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
				"hdr.ipv4.src_addr": &client.LpmMatch{
					Value: Pack32BinaryIP4(ip),
					PLen:  int32(plen),
				},
			},
			p4RtC.NewTableActionDirect("k8s_dp_control.set_ipset_match_result",
				[][]byte{valueToBytes8(mask)}),
			nil,
		)
		if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
			log.Errorf("Cannot insert entry into 'acl_ipset_match_table_ingress': %v", err)
			return err
		}

	case Delete:
		entryDelete := p4RtC.NewTableEntry(
			"k8s_dp_control.acl_ipset_match_table_ingress",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
				"hdr.ipv4.src_addr": &client.LpmMatch{
					Value: Pack32BinaryIP4(ip),
					PLen:  int32(plen),
				},
			},
			nil,
			nil,
		)
		if err := p4RtC.DeleteTableEntry(ctx, entryDelete); err != nil {
			log.Errorf("Cannot delete entry from 'acl_ipset_match_table_ingress': %v", err)
			return err
		}
	}

	return nil
}

func TcpDstPortRcTable(ctx context.Context, p4RtC *client.Client,
	polID uint16, portrange []uint16,
	action InterfaceType) error {
	switch action {
	case Insert:
		entryAdd := p4RtC.NewTableEntry(
			"k8s_dp_control.tcp_dport_rc_table",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
			},
			p4RtC.NewTableActionDirect("k8s_dp_control.do_range_check_tcp",
				[][]byte{valueToBytes16(portrange[0]),
					valueToBytes16(portrange[1]),
					valueToBytes16(portrange[2]),
					valueToBytes16(portrange[3]),
					valueToBytes16(portrange[4]),
					valueToBytes16(portrange[5]),
					valueToBytes16(portrange[6]),
					valueToBytes16(portrange[7]),
					valueToBytes16(portrange[8]),
					valueToBytes16(portrange[9]),
					valueToBytes16(portrange[10]),
					valueToBytes16(portrange[11]),
					valueToBytes16(portrange[12]),
					valueToBytes16(portrange[13]),
					valueToBytes16(portrange[14]),
					valueToBytes16(portrange[15])}),
			nil,
		)
		if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
			log.Errorf("Cannot insert entry into 'tcp_dport_rc_table': %v", err)
			return err
		}

	case Update:
		entryMod := p4RtC.NewTableEntry(
			"k8s_dp_control.tcp_dport_rc_table",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
			},
			p4RtC.NewTableActionDirect("k8s_dp_control.do_range_check_tcp",
				[][]byte{valueToBytes16(portrange[0]),
					valueToBytes16(portrange[1]),
					valueToBytes16(portrange[2]),
					valueToBytes16(portrange[3]),
					valueToBytes16(portrange[4]),
					valueToBytes16(portrange[5]),
					valueToBytes16(portrange[6]),
					valueToBytes16(portrange[7]),
					valueToBytes16(portrange[8]),
					valueToBytes16(portrange[9]),
					valueToBytes16(portrange[10]),
					valueToBytes16(portrange[11]),
					valueToBytes16(portrange[12]),
					valueToBytes16(portrange[13]),
					valueToBytes16(portrange[14]),
					valueToBytes16(portrange[15])}),
			nil,
		)
		if err := p4RtC.ModifyTableEntry(ctx, entryMod); err != nil {
			log.Errorf("Cannot update entry to 'tcp_dport_rc_table': %v", err)
			return err
		}

	case Delete:
		entryDelete := p4RtC.NewTableEntry(
			"k8s_dp_control.tcp_dport_rc_table",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
			},
			nil,
			nil,
		)
		if err := p4RtC.DeleteTableEntry(ctx, entryDelete); err != nil {
			log.Errorf("Cannot delete entry from 'tcp_dport_rc_table': %v", err)
			return err
		}
	}

	return nil
}

func UdpDstPortRcTable(ctx context.Context, p4RtC *client.Client,
	polID uint16, portrange []uint16,
	action InterfaceType) error {
	switch action {
	case Insert:
		entryAdd := p4RtC.NewTableEntry(
			"k8s_dp_control.udp_dport_rc_table",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
			},
			p4RtC.NewTableActionDirect("k8s_dp_control.do_range_check_udp",
				[][]byte{valueToBytes16(portrange[0]),
					valueToBytes16(portrange[1]),
					valueToBytes16(portrange[2]),
					valueToBytes16(portrange[3]),
					valueToBytes16(portrange[4]),
					valueToBytes16(portrange[5]),
					valueToBytes16(portrange[6]),
					valueToBytes16(portrange[7]),
					valueToBytes16(portrange[8]),
					valueToBytes16(portrange[9]),
					valueToBytes16(portrange[10]),
					valueToBytes16(portrange[11]),
					valueToBytes16(portrange[12]),
					valueToBytes16(portrange[13]),
					valueToBytes16(portrange[14]),
					valueToBytes16(portrange[15])}),
			nil,
		)
		if err := p4RtC.InsertTableEntry(ctx, entryAdd); err != nil {
			log.Errorf("Cannot insert entry into 'udp_dport_rc_table': %v", err)
			return err
		}

	case Update:
		entryMod := p4RtC.NewTableEntry(
			"k8s_dp_control.udp_dport_rc_table",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
			},
			p4RtC.NewTableActionDirect("k8s_dp_control.do_range_check_udp",
				[][]byte{valueToBytes16(portrange[0]),
					valueToBytes16(portrange[1]),
					valueToBytes16(portrange[2]),
					valueToBytes16(portrange[3]),
					valueToBytes16(portrange[4]),
					valueToBytes16(portrange[5]),
					valueToBytes16(portrange[6]),
					valueToBytes16(portrange[7]),
					valueToBytes16(portrange[8]),
					valueToBytes16(portrange[9]),
					valueToBytes16(portrange[10]),
					valueToBytes16(portrange[11]),
					valueToBytes16(portrange[12]),
					valueToBytes16(portrange[13]),
					valueToBytes16(portrange[14]),
					valueToBytes16(portrange[15])}),
			nil,
		)
		if err := p4RtC.ModifyTableEntry(ctx, entryMod); err != nil {
			log.Errorf("Cannot update entry in 'udp_dport_rc_table': %v", err)
			return err
		}

	case Delete:
		entryDelete := p4RtC.NewTableEntry(
			"k8s_dp_control.udp_dport_rc_table",
			map[string]client.MatchInterface{
				"meta.acl_pol_id": &client.ExactMatch{
					Value: valueToBytes16(polID),
				},
			},
			nil,
			nil,
		)
		if err := p4RtC.DeleteTableEntry(ctx, entryDelete); err != nil {
			log.Errorf("Cannot delete entry from 'udp_dport_rc_table': %v", err)
			return err
		}
	}

	return nil
}

func IsNamePresent(substr string, strslice []string) bool {
	for _, str := range strslice {
		if strings.Contains(str, substr) {
			return true
		}
	}
	log.Infof("name %s is not present in given slice", substr)
	return false
}

func IsSame(slice1 []uint16, slice2 []uint16) bool {
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			log.Infof("%d and %d are not same", slice1[i], slice2[i])
			return false
		}
	}
	return true
}

func InsertPolicyTableEntries(ctx context.Context, p4RtC *client.Client, tbltype tabletype, policy *store.Policy, workloadep *store.PolicyWorkerEndPoint) error {
	var err error
	switch tbltype {
	case policyadd:
		for ipsetidx, _ := range policy.IpSetIDx {
			for ruleid, _ := range policy.IpSetIDx[ipsetidx].RuleID {
				cidr := policy.IpSetIDx[ipsetidx].RuleID[ruleid].Cidr
				mask := policy.IpSetIDx[ipsetidx].RuleID[ruleid].RuleMask
				if policy.IpSetIDx[ipsetidx].Direction == "TX" {
					err = AclIpSetMatchTableEgress(ctx, p4RtC, ipsetidx, cidr, mask, Insert)
					if err != nil {
						log.Errorf("AclIpSetMatchTableEgress failed %v", err)
					} else {
						log.Infof("AclIpSetMatchTableEgress passed")
					}
				} else {
					err = AclIpSetMatchTableIngress(ctx, p4RtC, ipsetidx, cidr, mask, Insert)
					if err != nil {
						log.Errorf("AclIpSetMatchTableIngress failed %v", err)
					} else {
						log.Infof("AclIpSetMatchTableIngress passed")
					}
				}
			}
			if len(policy.IpSetIDx[ipsetidx].Rc) != 0 {
				err = TcpDstPortRcTable(ctx, p4RtC, ipsetidx, policy.IpSetIDx[ipsetidx].Rc, Insert)
				if err != nil {
					log.Errorf("TcpDstPortRcTable failed %v", err)
				} else {
					log.Infof("TcpDstPortRcTable passed")
				}
			} else {
				err = UdpDstPortRcTable(ctx, p4RtC, ipsetidx, policy.IpSetIDx[ipsetidx].Rc, Insert)
				if err != nil {
					log.Errorf("UdpDstPortRcTable failed %v", err)
				} else {
					log.Infof("UdpDstPortRcTable passed")
				}
			}
		}
		log.Infof("Inserted policy entry into pipeline %s", policy)

	case policydel:
		for ipsetidx, _ := range policy.IpSetIDx {
			for ruleid, _ := range policy.IpSetIDx[ipsetidx].RuleID {
				cidr := policy.IpSetIDx[ipsetidx].RuleID[ruleid].Cidr
				if policy.IpSetIDx[ipsetidx].Direction == "TX" {
					err = AclIpSetMatchTableEgress(ctx, p4RtC, ipsetidx, cidr, 0, Delete)
					if err != nil {
						log.Errorf("AclIpSetMatchTableEgress failed %v", err)
					} else {
						log.Infof("AclIpSetMatchTableEgress passed")
					}
				} else {
					err = AclIpSetMatchTableIngress(ctx, p4RtC, ipsetidx, cidr, 0, Delete)
					if err != nil {
						log.Errorf("AclIpSetMatchTableIngress failed %v", err)
					} else {
						log.Infof("AclIpSetMatchTableIngress passed")
					}
				}
			}
			if len(policy.IpSetIDx[ipsetidx].Rc) != 0 {
				if policy.IpSetIDx[ipsetidx].Protocol == "TCP" {
					err = TcpDstPortRcTable(ctx, p4RtC, ipsetidx, nil, Delete)
					if err != nil {
						log.Errorf("TcpDstPortRcTable failed %v", err)
					} else {
						log.Infof("TcpDstPortRcTable passed")
					}
				} else {
					err = UdpDstPortRcTable(ctx, p4RtC, ipsetidx, nil, Delete)
					if err != nil {
						log.Errorf("TcpDstPortRcTable failed %v", err)
					} else {
						log.Infof("TcpDstPortRcTable passed")
					}
				}
			}
		}
		log.Infof("deleted policy entry %s", policy)

	case policyupdate:
		policyold := store.PolicySet.PolicyMap[policy.PolicyName]
		oldrules := make([]string, 0)
		newrules := make([]string, 0)
		var index int
		//create slice of existing rules
		for ipsetidx, _ := range policyold.IpSetIDx {
			for ruleid, _ := range policyold.IpSetIDx[ipsetidx].RuleID {
				oldrules = append(oldrules, ruleid)
			}
		}
		log.Infof("old rules: %s", oldrules)
		//create slice of new rules
		for ipsetidx, _ := range policy.IpSetIDx {
			for ruleid, _ := range policy.IpSetIDx[ipsetidx].RuleID {
				newrules = append(newrules, ruleid)
			}
		}
		log.Infof("new rules: %s", newrules)
		log.Infof("delete old rules if its not part of updated set of rules")
		for ipsetidx, _ := range policyold.IpSetIDx {
			rc := make([]uint16, 16)
			for ruleid, _ := range policyold.IpSetIDx[ipsetidx].RuleID {
				index++
				if IsNamePresent(ruleid, newrules) {
					if len(policyold.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange) != 0 {
						i := index
						j := i + 1
						rc[i] = policyold.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange[0]
						rc[j] = policyold.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange[1]
					}
				} else { //if rule is not present, then delete only that rule from pipeline
					if len(policyold.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange) != 0 {
						i := index
						j := i + 1
						rc[i] = 0
						rc[j] = 0
					}
					cidr := policyold.IpSetIDx[ipsetidx].RuleID[ruleid].Cidr
					if policyold.IpSetIDx[ipsetidx].Direction == "TX" {
						err = AclIpSetMatchTableEgress(ctx, p4RtC, ipsetidx, cidr, 0, Delete)
						if err != nil {
							log.Errorf("AclIpSetMatchTableEgress failed %v", err)
						} else {
							log.Infof("AclIpSetMatchTableEgress passed")
						}
					} else {
						err = AclIpSetMatchTableIngress(ctx, p4RtC, ipsetidx, cidr, 0, Delete)
						if err != nil {
							log.Errorf("AclIpSetMatchTableIngress failed %v", err)
						} else {
							log.Infof("AclIpSetMatchTableIngress passed")
						}
					}
				}
			}
			index = 0
			log.Infof("rc is: %s", rc)
			if len(rc) == 0 {
				if policyold.IpSetIDx[ipsetidx].Protocol == "TCP" {
					err = TcpDstPortRcTable(ctx, p4RtC, ipsetidx, nil, Delete)
					if err != nil {
						log.Errorf("TcpDstPortRcTable failed %v", err)
					} else {
						log.Infof("TcpDstPortRcTable passed")
					}
				} else {
					err = UdpDstPortRcTable(ctx, p4RtC, ipsetidx, nil, Delete)
					if err != nil {
						log.Errorf("UdpDstPortRcTable failed %v", err)
					} else {
						log.Infof("UdpDstPortRcTable passed")
					}
				}
			} else {
				if policyold.IpSetIDx[ipsetidx].Protocol == "TCP" {
					err = TcpDstPortRcTable(ctx, p4RtC, ipsetidx, rc, Update)
					if err != nil {
						log.Errorf("TcpDstPortRcTable failed %v", err)
					} else {
						log.Infof("TcpDstPortRcTable passed")
					}
				} else {
					err = UdpDstPortRcTable(ctx, p4RtC, ipsetidx, rc, Update)
					if err != nil {
						log.Errorf("UdpDstPortRcTable failed %v", err)
					} else {
						log.Infof("UdpDstPortRcTable passed")
					}
				}
			}
		}
		//add newly added rules to the pipeline
		for ipsetidx, _ := range policy.IpSetIDx {
			rc := make([]uint16, 16)
			for ruleid, _ := range policy.IpSetIDx[ipsetidx].RuleID {
				index++
				if IsNamePresent(ruleid, oldrules) {
					if len(policy.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange) != 0 {
						i := index
						j := i + 1
						rc[i] = policy.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange[0]
						rc[j] = policy.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange[1]
					}
				} else { //if rule is not present, then insert only that rule to pipeline
					if len(policy.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange) != 0 {
						i := index
						j := i + 1
						rc[i] = policy.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange[0]
						rc[j] = policy.IpSetIDx[ipsetidx].RuleID[ruleid].PortRange[1]
					}
					cidr := policy.IpSetIDx[ipsetidx].RuleID[ruleid].Cidr
					mask := policy.IpSetIDx[ipsetidx].RuleID[ruleid].RuleMask
					if policy.IpSetIDx[ipsetidx].Direction == "TX" {
						err = AclIpSetMatchTableEgress(ctx, p4RtC, ipsetidx, cidr, mask, Delete)
						if err != nil {
							log.Errorf("AclIpSetMatchTableEgress failed %v", err)
						} else {
							log.Infof("AclIpSetMatchTableEgress passed")
						}
					} else {
						err = AclIpSetMatchTableIngress(ctx, p4RtC, ipsetidx, cidr, mask, Delete)
						if err != nil {
							log.Errorf("AclIpSetMatchTableIngress failed %v", err)
						} else {
							log.Infof("AclIpSetMatchTableIngress passed")
						}
					}
				}
			}
			index = 0
			log.Infof("rc: %s", rc)
			if len(rc) != 0 {
				if len(policyold.IpSetIDx[ipsetidx].Rc) == 0 { //if new rule has port range and previously we have not added any port range for same ipsetidx, then insert port range for that ipsetidx
					if policy.IpSetIDx[ipsetidx].Protocol == "TCP" {
						err = TcpDstPortRcTable(ctx, p4RtC, ipsetidx, rc, Insert)
						if err != nil {
							log.Errorf("TcpDstPortRcTable failed %v", err)
						} else {
							log.Infof("TcpDstPortRcTable passed")
						}
					} else {
						err = UdpDstPortRcTable(ctx, p4RtC, ipsetidx, rc, Insert)
						if err != nil {
							log.Errorf("UdpDstPortRcTable failed %v", err)
						} else {
							log.Infof("UdpDstPortRcTable passed")
						}
					}
				}
			} else if len(rc) != 0 && len(policyold.IpSetIDx[ipsetidx].Rc) != 0 { //if old port range and new port range for a ipsetidx are different
				if !IsSame(rc, policyold.IpSetIDx[ipsetidx].Rc) {
					if policy.IpSetIDx[ipsetidx].Protocol == "TCP" {
						err = TcpDstPortRcTable(ctx, p4RtC, ipsetidx, rc, Update)
						if err != nil {
							log.Errorf("TcpDstPortRcTable failed %v", err)
						} else {
							log.Infof("TcpDstPortRcTable passed")
						}
					} else {
						err = UdpDstPortRcTable(ctx, p4RtC, ipsetidx, rc, Update)
						if err != nil {
							log.Errorf("UdpDstPortRcTable failed %v", err)
						} else {
							log.Infof("UdpDstPortRcTable passed")
						}
					}
				}
			} else {
			}
		}
		log.Infof("updated policy: %s", policy)

	case workloadadd:
		for _, policyname := range workloadep.PolicyNameIngress {
			policy := store.PolicySet.PolicyMap[policyname]
			for ipsetidx, _ := range policy.IpSetIDx {
				err = AclPodIpProtoTableIngress(ctx, p4RtC, policy.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, ipsetidx, ipsetidx, Insert)
				if err != nil {
					log.Errorf("AclPodIpProtoTableIngress failed %v", err)
				} else {
					log.Infof("AclPodIpProtoTableIngress passed")
				}
			}
		}

		for _, policyname := range workloadep.PolicyNameEgress {
			policy := store.PolicySet.PolicyMap[policyname]
			for ipsetidx, _ := range policy.IpSetIDx {
				err = AclPodIpProtoTableEgress(ctx, p4RtC, policy.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, ipsetidx, ipsetidx, Insert)
				if err != nil {
					log.Errorf("AclPodIpProtoTableEgress failed %v", err)
				} else {
					log.Infof("AclPodIpProtoTableEgress passed")
				}
			}
		}
		log.Infof("workloadep added: %s", workloadep)

	case workloaddel:
		for _, policyname := range workloadep.PolicyNameIngress {
			policy := store.PolicySet.PolicyMap[policyname]
			for ipsetidx, _ := range policy.IpSetIDx {
				err = AclPodIpProtoTableIngress(ctx, p4RtC, policy.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, ipsetidx, ipsetidx, Delete)
				if err != nil {
					log.Errorf("AclPodIpProtoTableIngress failed %v", err)
				} else {
					log.Infof("AclPodIpProtoTableIngress passed")
				}
			}
		}

		for _, policyname := range workloadep.PolicyNameEgress {
			policy := store.PolicySet.PolicyMap[policyname]
			for ipsetidx, _ := range policy.IpSetIDx {
				err = AclPodIpProtoTableEgress(ctx, p4RtC, policy.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, ipsetidx, ipsetidx, Delete)
				if err != nil {
					log.Errorf("AclPodIpProtoTableEgress failed %v", err)
				} else {
					log.Infof("AclPodIpProtoTableEgress passed")
				}
			}
		}
		log.Infof("workload deleted: %s", workloadep)

	case workloadupdate:
		workloadepold := store.PolicySet.WorkerEpMap[workloadep.WorkerEp]
		//ingress policy names
		//delete from policy tables for removed policies
		for _, policyname := range workloadepold.PolicyNameIngress {
			if !IsNamePresent(policyname, workloadep.PolicyNameIngress) { //if policyname from old store entry is not present in new entry, then delete
				policydel := store.PolicySet.PolicyMap[policyname]
				for ipsetidx, _ := range policydel.IpSetIDx {
					err = AclPodIpProtoTableIngress(ctx, p4RtC, policydel.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, 0, 0, Delete)
					if err != nil {
						log.Errorf("AclPodIpProtoTableIngress failed %v", err)
					} else {
						log.Infof("AclPodIpProtoTableIngress passed")
					}
				}
			}
		}
		//insert to policy tables the new policies
		for _, policyname := range workloadep.PolicyNameIngress {
			if !IsNamePresent(policyname, workloadepold.PolicyNameIngress) {
				policyadd := store.PolicySet.PolicyMap[policyname]
				for ipsetidx, _ := range policyadd.IpSetIDx {
					err = AclPodIpProtoTableIngress(ctx, p4RtC, policyadd.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, ipsetidx, ipsetidx, Insert)
					if err != nil {
						log.Errorf("AclPodIpProtoTableIngress failed %v", err)
					} else {
						log.Infof("AclPodIpProtoTableIngress passed")
					}
				}
			}
		}

		//egress policy names
		//delete from policy tables for removed policies
		for _, policyname := range workloadepold.PolicyNameEgress {
			if !IsNamePresent(policyname, workloadep.PolicyNameEgress) { //if policyname from old store entry is not present in new entry, then delete
				policydel := store.PolicySet.PolicyMap[policyname]
				for ipsetidx, _ := range policydel.IpSetIDx {
					err = AclPodIpProtoTableEgress(ctx, p4RtC, policydel.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, 0, 0, Delete)
					if err != nil {
						log.Errorf("AclPodIpProtoTableEgress failed %v", err)
					} else {
						log.Infof("AclPodIpProtoTableEgress passed")
					}
				}
			}
		}
		//insert to policy tables the new policies
		for _, policyname := range workloadep.PolicyNameEgress {
			if !IsNamePresent(policyname, workloadepold.PolicyNameEgress) {
				policyadd := store.PolicySet.PolicyMap[policyname]
				for ipsetidx, _ := range policyadd.IpSetIDx {
					err = AclPodIpProtoTableEgress(ctx, p4RtC, policyadd.IpSetIDx[ipsetidx].Protocol, workloadep.WorkerEp, ipsetidx, ipsetidx, Insert)
					if err != nil {
						log.Errorf("AclPodIpProtoTableEgress failed %v", err)
					} else {
						log.Infof("AclPodIpProtoTableEgress passed")
					}
				}
			}
		}
		log.Infof("updated workloadep: %s", workloadep)
	}
	return nil
}
