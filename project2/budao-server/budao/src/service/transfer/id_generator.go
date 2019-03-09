package transfer

import (
	"fmt"
	_ "hash"
	"hash/crc32"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/sumaig/glog"

	"common"
	"db"
	_ "twirprpc"
)

//功能：生成唯一id
/* 格式：64字节id
|     10      |  22  |  16  |             16              |
|mysql数据表id|自增id|随机数|crc32校验前边48位数据，取前16|
*/

var (
	gIdGenerator    *IdGenerator
	onceIdGenerator sync.Once
)

// GetIdGenerator function create id generator instance
func GetIdGenerator() *IdGenerator {
	onceIdGenerator.Do(func() {
		if gIdGenerator == nil {
			gIdGenerator = &IdGenerator{}
			gIdGenerator.Init()
		}
	})
	return gIdGenerator
}

// IdGenerator structure
type IdGenerator struct {
	mysqlClient *db.MysqlClient
	tableDesc   map[string]uint64
	randSource  rand.Source
}

// Init function completes id generator initialization
func (s *IdGenerator) Init() {
	s.tableDesc = common.GetConfig().DB.MySQL[common.BUDAODB].TableDesc
	s.randSource = rand.NewSource(time.Now().UnixNano())
}

// GetItemId function genrate sweater id.
func (s *IdGenerator) GetItemId(tableNamePrefix string) (itemID, autoIncreID uint64, err error) {

	//1. table num
	tableNum := uint64(rand.New(s.randSource).Int63n(int64(s.tableDesc[tableNamePrefix])))

	//2. get auto-increment id. TODO get from redis
	autoIncreID, err = s.GetAutoIncrementId(tableNamePrefix, tableNum)
	if err != nil || autoIncreID > 4194303 {
		return
	}

	//3. random
	random := rand.New(s.randSource).Uint64()

	beforeNum := (tableNum << 54) | (autoIncreID << 32) | ((random >> 48) << 16)
	beforeStr := strconv.FormatUint((beforeNum >> 16), 10)

	//4. crc32
	crc := crc32.ChecksumIEEE([]byte(beforeStr))

	//5. compose videoid
	itemID = beforeNum | (uint64(crc >> 16))

	return
}

// GetAutoIncrementId function get autoincrement id of table
func (s *IdGenerator) GetAutoIncrementId(tableNamePrefix string, tableNum uint64) (autoIncrementID uint64, err error) {
	tableName := fmt.Sprintf("%s%d", tableNamePrefix, tableNum)
	execSQL := fmt.Sprintf("insert into %s values()", tableName)

	result, err := db.Exec(common.BUDAODB, execSQL)
	if err != nil {
		glog.Error("insert failed. execSql:[%s] err:%v", execSQL, err)
		return
	}
	affectNum, err := result.RowsAffected()
	if err != nil {
		return
	}
	autoIncrementIdtmp, err := result.LastInsertId()
	if err != nil {
		return
	}
	autoIncrementID = uint64(autoIncrementIdtmp)

	glog.Debug("insert success. RowsAffected:%d LastInsertId:%d", affectNum, autoIncrementID)

	return
}
