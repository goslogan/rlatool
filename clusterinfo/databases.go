/*
databases.go provides a parser for the database information in the rladmin output
Copyright © 2024 Nic Gibson <nic.gibson@redis.com>
*/
package clusterinfo

import (
	"encoding/json"
	"strings"

	"github.com/gocarina/gocsv"
)

type DBEndPoints []string

type DBShards struct {
	Masters  uint16
	Replicas uint16
}
type DBNodes map[string]*DBShards

type Database struct {
	Id                string       `csv:"DB:ID" json:"id"`
	Name              string       `csv:"NAME" json:"name"`
	Type              string       `csv:"TYPE" json:"type"`
	Status            string       `csv:"STATUS" json:"status"`
	MasterShards      uint16       `csv:"SHARDS" json:"shards"`
	Placement         string       `csv:"PLACEMENT" json:"placement"`
	Replication       string       `csv:"REPLICATION" json:"replication"`
	Persistence       string       `csv:"PERSISTENCE" json:"persistence"`
	Endpoint          DBEndPoints  `csv:"ENDPOINT" json:"endpoints"`
	ExecState         string       `csv:"EXEC_STATE" json:"execState"`
	ExecStateMachine  string       `csv:"EXEC_STATE_MACHINE" json:"execStateMachine"`
	BackupProgress    string       `csv:"BACKUP_PROGRESS" json:"backupProgress"`
	MissingBackupTime string       `csv:"MISSING_BACKUP_TIME" json:"missingBackupTime"`
	RedisVersion      string       `csv:"REDIS_VERSION" json:"redisVersion"`
	parent            *ClusterInfo `csv:"-" json:"-"`
}

type DatabaseWithNodes struct {
	Database
	Nodes DBNodes `json:"nodes" csv:"NODES"`
}

type Databases []*Database
type DatabasesWithNodes []*DatabaseWithNodes

func (c *Chunks) ParseDatabases(parent *ClusterInfo) (Databases, error) {

	databases := []*Database{}
	err := gocsv.UnmarshalBytes(c.Databases, &databases)
	if err != nil {
		return nil, err
	}
	for _, db := range databases {
		db.parent = parent
	}

	return databases, nil
}

// JSON returns the database struct marsalled to JSON
func (db *Database) JSON() (string, error) {
	if out, err := json.Marshal(db); err != nil {
		return "", err
	} else {
		return string(out), nil
	}
}

// OnNode returns the number of shards on the given node for a database.
func (db *Database) OnNode(id string) DBShards {
	var masters, replicas uint16
	for _, shard := range db.parent.Shards.ForDB(db.Id) {
		if shard.Node == id {
			if shard.Role == "master" {
				masters++
			} else {
				replicas++
			}
		}
	}

	return DBShards{Masters: masters, Replicas: replicas}
}

func (d *Database) withNodes() *DatabaseWithNodes {
	nodes := DBNodes{}

	for _, node := range d.parent.Nodes {
		nodes[node.Id] = &DBShards{}
	}

	for _, shard := range d.parent.Shards {
		if shard.DBId == d.Id {
			shardCount := nodes[shard.Node]

			if shard.Role == "master" {
				shardCount.Masters++
			} else {
				shardCount.Replicas++
			}
		}
	}

	return &DatabaseWithNodes{
		Database: *d,
		Nodes:    d.getNodes(),
	}
}

// ShardCount returns the total number of shards by
// counting them.
func (d *Database) ShardCount() uint16 {
	shards := uint16(0)
	for _, v := range d.getNodes() {
		shards += v.Masters
		shards += v.Replicas
	}

	return shards

}

func (d *Database) getNodes() DBNodes {
	nodes := DBNodes{}

	for _, node := range d.parent.Nodes {
		nodes[node.Id] = &DBShards{}
	}

	for _, shard := range d.parent.Shards {
		if shard.DBId == d.Id {
			shardCount := nodes[shard.Node]

			if shard.Role == "master" {
				shardCount.Masters++
			} else {
				shardCount.Replicas++
			}
		}
	}

	return nodes
}

func (d *Databases) JSON() (string, error) {
	if out, err := json.Marshal(d); err != nil {
		return "", err
	} else {
		return string(out), nil
	}
}

func (d Databases) CSV() (string, error) {
	return gocsv.MarshalString(d)
}

func (d DatabasesWithNodes) JSON() (string, error) {
	if out, err := json.Marshal(d); err != nil {
		return "", err
	} else {
		return string(out), nil
	}
}

func (d DatabasesWithNodes) CSV() (string, error) {
	return gocsv.MarshalString(d)
}

func (d Databases) withNodes() DatabasesWithNodes {
	dn := DatabasesWithNodes{}
	for _, db := range d {
		dn = append(dn, db.withNodes())
	}

	return dn
}

func (e *DBEndPoints) UnmarshalCSV(input string) error {
	tmp := DBEndPoints(strings.Split(input, "/"))
	*e = tmp
	return nil
}

func (e *DBEndPoints) MarshalCSV() (string, error) {
	return strings.Join(*e, "/"), nil
}

func (n *DBNodes) MarshalCSV() (string, error) {

	keys := []string{}

	for k, v := range *n {
		if v.Masters+v.Replicas > 0 {
			keys = append(keys, k)
		}
	}

	return strings.Join(keys, "/"), nil
}
