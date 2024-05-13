/*
nodes.go provides a parser for the node information in the rladmin output
Copyright © 2024 Nic Gibson <nic.gibson@redis.com>
*/
package clusterinfo

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/gocarina/gocsv"
)

const (
	errorString = "unable to parse '%s' as %s: %w"
)

type IP struct {
	net.IP
}

type ShardInfo struct {
	InUse uint16 `json:"shardsInUse"`
	Max   uint16 `json:"maxShards"`
}

type MemoryInfo struct {
	Free RAMFloat `json:"free"`
	Max  RAMFloat `json:"max"`
}

type Node struct {
	Id               string       `json:"id" csv:"NODE:ID"`
	Role             string       `json:"role" csv:"ROLE"`
	Address          IP           `json:"address" csv:"ADDRESS"`
	ExternalAddress  IP           `json:"externalAddress" csv:"EXTERNAL_ADDRESS"`
	HostName         string       `json:"hostName" csv:"HOSTNAME"`
	OverbookingDepth RAMFloat     `json:"overbookingDepth" csv:"OVERBOOKING_DEPTH"`
	Masters          uint16       `json:"masters" csv:"MASTERS"`
	Replicas         uint16       `json:"replicas" csv:"SLAVES"`
	ShardUsage       ShardInfo    `json:"shards" csv:"SHARDS"`
	Cores            uint16       `json:"cores" csv:"CORES"`
	RedisRAM         MemoryInfo   `json:"redisRAM" csv:"FREE_RAM"`
	ProvisionalRAM   MemoryInfo   `json:"provisionalRAM" csv:"PROVISIONAL_RAM"`
	Version          string       `json:"version" csv:"VERSION"`
	SHA              string       `json:"sha" csv:"SHA"`
	RackId           string       `json:"rackId" csv:"RACK-ID"`
	Status           string       `json:"status" csv:"STATUS"`
	Quorum           bool         `json:"quorum" csv:"-"`
	parent           *ClusterInfo `csv:"-" json:"-"`
}

type Nodes []*Node

func (c *Chunks) ParseNodes(parent *ClusterInfo) (Nodes, error) {

	nodes := []*Node{}
	err := gocsv.UnmarshalBytes(c.Nodes, &nodes)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		node.parent = parent
		if node.ShardUsage.Max == 0 {
			node.Quorum = true
		}
		// strip of the "*" prefix from the node that ran the rladmin status
		node.Id = strings.TrimPrefix(node.Id, "*")
	}

	return nodes, nil
}

func (m *MemoryInfo) UnmarshalCSV(input string) error {
	if parts := strings.Split(input, "/"); len(parts) == 2 {
		f, err := parseMemory(parts[0])
		if err != nil {
			return fmt.Errorf(errorString, "memory info", parts[0], err)
		} else {
			m.Free = RAMFloat(f)
		}
		f, err = parseMemory(parts[1])
		if err != nil {
			return fmt.Errorf(errorString, "memory info", parts[1], err)
		} else {
			m.Max = RAMFloat(f)
		}
	} else {
		return fmt.Errorf("unable to split %s into parts for memory info", input)
	}

	return nil
}

func (i *IP) UnmarshalCSV(input string) error {
	i.IP = net.ParseIP(input)
	if i.IP == nil {
		return fmt.Errorf(errorString, "address", input, nil)
	} else {
		return nil
	}
}

func (s *ShardInfo) UnmarshalCSV(input string) error {
	var err error
	if parts := strings.Split(input, "/"); len(parts) == 2 {
		if s.InUse, err = toUint16(parts[0]); err != nil {
			return fmt.Errorf(errorString, "number of shards", parts[0], err)
		}
		if s.Max, err = toUint16(parts[1]); err != nil {
			return fmt.Errorf(errorString, "maximum number of shards", parts[1], err)
		}
	} else {
		return fmt.Errorf("unable to split %s into parts for shard counts", input)
	}

	return nil
}

func (ns Nodes) JSON() (string, error) {
	data, err := json.Marshal(&ns)
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func (ns Nodes) CSV() (string, error) {
	return gocsv.MarshalString(&ns)
}
