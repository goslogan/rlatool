/*
chunks.go provides a base parser that loads rladmin output and parses it into chunks for each section.
Copyright © 2024 Nic Gibson <nic.gibson@redis.com>
*/
package clusterinfo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// Chunks is used to store the output of the base parser.
type Chunks struct {
	Intro     string
	Cluster   string
	Nodes     []byte
	Databases []byte
	Endpoints []byte
	Shards    []byte
}

var marker = regexp.MustCompile(`^[A-Z ]+:$`)
var spaceReplacer = regexp.MustCompile(`[\t\f\r ]+`)
var endpointCleaner = regexp.MustCompile(`\+\d+`)

const (
	ChunkIntro = iota
	ChunkCluster
	ChunkNodes
	ChunkDatabases
	ChunkEndpoints
	ChunkShards
)

func (c *Chunks) Parse(input io.Reader) error {

	current := make([]byte, 0)

	where := ChunkIntro
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := scanner.Bytes()
		if marker.Match(line) {
			c.putData(current, where)
			where++
			current = make([]byte, 0)
		} else {
			line := append(bytes.TrimSpace(line), '\n')
			current = append(current, line...)
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	} else {
		c.putData(current, where)
	}

	return nil
}

func (c *Chunks) putData(data []byte, stage int) {
	if len(data) > 0 {
		switch stage {
		case ChunkIntro:
			c.Intro = string(data) // Don't convert this
		case ChunkCluster:
			c.Cluster = string(data) // Don't convert this
		case ChunkNodes:
			c.Nodes = c.toCSV(data)
		case ChunkDatabases:
			c.Databases = c.toCSV(data)
		case ChunkEndpoints:
			c.Endpoints = c.toCSV(c.cleanEndpoints(data))
		case ChunkShards:
			c.Shards = c.toCSV(data)
		}
	}
}

// this works for rladmin output because there are no quotes or commas in inconvenient
// locations but it is far from general
func (c *Chunks) toCSV(data []byte) []byte {
	return spaceReplacer.ReplaceAll(data, []byte{','})
}

// we need to do a fixup for endpoints because of node isolation.  It's not of use to
// this program so we strip out the markers
func (c *Chunks) cleanEndpoints(data []byte) []byte {
	return endpointCleaner.ReplaceAll(data, []byte{})
}

// ExtractTimeStamp finds the timestamp at the start of the output and returns it as time.Time
func (c *Chunks) ExtractTimeStamp() (time.Time, error) {

	lines := strings.Split(c.Intro, "\n")
	if len(lines) < 2 {
		return time.Now(), fmt.Errorf("rlatool - timestamp not found in input")
	} else {
		return time.Parse("2006-01-02 03:04:05.000000-07:00", lines[1])
	}

}
