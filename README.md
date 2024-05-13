# rlatool


rlatool parses, merges and outputs data from Redis Inc's rladmin tools.
Specifically, it takes the output of *rladmin status* and generates CSV and JSON
files from it. These can be used for further analysis.

## Installation

If you have [go](https://go.dev/) installed simply run

```
go install github.com/goslogan/rlatool
```

or clone the repo and build it:

```
git clone https://github.com/goslogan/rlatool
cd rlatool
go build
```

If not, I've provided builds for Mac OS, Windows and Linux (on x64) in [Releases](https://github.com/goslogan/rlatool) on github. Due to code signing issues, the Mac OS build is painful to install. Ask me for help if you need it. 

## Usage

```
  rlatool [command]

Available Commands:
  all         Generate a JSON conversion of rladmin status output
  completion  Generate the autocompletion script for the specified shell
  databases   Generate a list of all databases on the cluster
  endpoints   Generate a list of all endpoints in the cluster
  help        Help about any command
  nodes       Generate a list of all nodes in the cluster
  report      Generate a report from the parsed input and one or more templaets
  shards      Generate a list of all shards in the cluster

Flags:
  -h, --help                help for rlatool
  -i, --input string        The rladmin status output to be parsed; if not provided STDIN is assumed
  -o, --output string       The file to which output should be written; if not provided STDOUT is assumed

```

All of `databases`, `endpoints`, `nodes` and `shards` take one of `--json` and `--csv` as arguments to define the output type (there is no default for this). 

### Commands

__all__:  Generate a JSON conversion of rladmin status output

__completion__: `rlatool completion [shell]` Generate the autocompletion script for rlatool for the specified shell.

__databases__: List all the databases in the cluster including all the fields generated by rladmin status. If the `--nodes` argument is provided, the output will include the nodes on which the database has shards

__endpoints__: List all the shards in the cluster including all the fields generated by rladmin status.

__help__: Help provides help for any command in the application.

__nodes__: List all the nodes in the cluster including all the fields generated by rladmin status.

__reports__: Generate a report using Go's template packages. 

__shards__: List all the shards in the cluster including all the fields generated by rladmin status.




### Generating reports

```
The template is passed the ClusterInfo struct generated by the parser
(see clusterinfo.go) as ".Info" and the title provided on the command line as 
".Title" . The template is executed and the result is output. 

See  https://pkg.go.dev/text/template for details. If the --html flag is passed
templates are processed by html/template instead

Usage:
  rlatool report [flags]

Flags:
  -h, --help                    help for report
  -H, --html                    Set to true to use HTML template engine when processing templates
  -t, --templates stringArray   List of templates to process, the first is assumed to be the root template
  -T, --title string            The title of the report if required (default "rlatool Report")

Global Flags:
  -i, --input string    The rladmin status output to be parsed; if not provided STDIN is assumed
  -o, --output string   The file to which output should be written; if not provided STDOUT is assumed

```

The __report__ command allows you to provide a report template which will be filled out using data from the parsed status output. This template language is that used by the Go [text/template](https://pkg.go.dev/text/template). If you use the `--html` flag the [html/template](https://pkg.go.dev/html/template) package will be used to render the output (the primary difference between the two is safe html escaping). The templating language is the same either way.

There are two predefined variables passed into the template 

* .Status - this is the parsed status report
* .Title - this is the value of the `--title` parameter passed on the command line.

There is an example of a template in the `templates` directory of the repo. This produces a table showing all the nodes and databases and where the master and replica shards are held along with the number of cores on each node. 

Run it as:

```
rlatool report --templates=templates/dbnodes.html --html --input=testdata/node_1.rladmin --output=test.html 
```

For all the fields listed below, the format and meaning is the same as that in `rladmin status` unless otherwise documented.

#### Status

The top level Status object has these fields (all arrays in the order found in the `rladmin status` output )

* Nodes
* Databases
* Endpoints
* Shards

#### Node 

A `Node` has these fields

* Id
* Role (master or slave)
* Address (IP object. Using the .String method to get a useful output)
* ExternalAddress (also an IP object)
* HostName   
* OverbookingDepth (float expressed in GB)
* Masters
* Replicas
* ShardUsage (see below)
* Cores 
* RedisRAM (as OverbookingDepth)
* ProvisionalRAM  
* Version 
* SHA     
* RackId  
* Status  
* Quorum  

The `ShardUsage` struct has the following fields

* InUse
* Max

#### Database

A `Database` has these fields 

* Id 
* Name
* Type
* Status
* Shards
* Placement
* Replication
* Persistence
* Endpoint
* ExecState
* ExecStateMachine 
* BackupProgress
* MissingBackupTime
* RedisVersion

The `Endpoint` value is an array of strings. 


The `Database` object has the following utility methods:

* OnNode -  Given a node id as a string returns a DBShards object giving the number of master and replica shards on that node. DBShards has the following fields:
  *. Masters
  *. Replicas
* ShardCount - returns the total number of shards in a database by adding up all masters and replicas


#### Endpoints

The `Endpoint` object has the following fields:

* Id
* DBId
* Name
* Node
* Role
* SSL 
* WatchdogStatus

##### Shards

The `Shard` object has the following fields:

* Id 
* DBId
* Name
* Node
* Role
* Slots
* UsedMemory 
* BackupProgress 
* RAMFrag
* WatchdogStatus 
* Status         
	