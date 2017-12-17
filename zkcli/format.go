package zkcli

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const timeFormat = "Mon Jan 02 15:04:05 GMT 2006"

func fmtTime(t int64) string {
	return time.Unix(t/1000, 0).UTC().Format(timeFormat)
}

func fmtStat(stat *zk.Stat) string {
	return fmt.Sprintf(`cZxid = 0x%x
ctime = %s
mZxid = 0x%x
mtime = %s
pZxid = 0x%x
cversion = %d
dataVersion = %d
aclVersion = %d
ephemeralOwner = 0x%x
dataLength = %d
numChildren = %d`,
		stat.Czxid, fmtTime(stat.Ctime), stat.Mzxid, fmtTime(stat.Mtime),
		stat.Pzxid, stat.Cversion, stat.Version, stat.Aversion,
		stat.EphemeralOwner, stat.DataLength, stat.NumChildren)
}
