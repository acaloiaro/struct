# struct

A simple unix utility for structuring string data on the command line.

# example usage

## string output
```
$ ls -l | struct -fields permissions,count,user,group,,,,,file_name -output string
permissions:-rw-r--r-- count:1 user:adriano group:staff 84 Sep 23 22:55 file_name:README.md
permissions:-rw------- count:1 user:adriano group:staff 39 Sep 23 21:08 file_name:go.mod
permissions:-rw-r--r-- count:1 user:adriano group:staff 2531 Sep 23 22:54 file_name:struct.go
```

## json output

```bash
$ ls -l | struct -fields permissions,count,user,group,,,,,file_name -output json
{"count":"1","file_name":"README.md","group":"staff","permissions":"-rw-r--r--","user":"adriano"}
{"count":"1","file_name":"go.mod","group":"staff","permissions":"-rw-------","user":"adriano"}
{"count":"1","file_name":"struct.go","group":"staff","permissions":"-rw-r--r--","user":"adriano"}
```
