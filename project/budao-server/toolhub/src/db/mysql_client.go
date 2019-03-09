package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sumaig/glog"
)

// MySQLConfig structure
type MySQLConfig struct {
	Master    string            `json:"master"`
	Slave     string            `json:"slave"`
	Passwd    string            `json:"passwd"`
	User      string            `json:"user"`
	Database  string            `json:"database"`
	TableDesc map[string]uint64 `json:"tableDesc"`
}

// MysqlClient express the client of master and slave
type MysqlClient struct {
	master       *sql.DB
	slave        *sql.DB
	tableDescMap map[string]uint64
}

var (
	once      sync.Once
	clientMap map[string]*MysqlClient
)

var SERVICE_NAME_ERR = errors.New("input service name error")

func Query(serviceName, querySQL string) (rows *sql.Rows, err error) {
	c, ok := clientMap[serviceName]
	if !ok {
		return nil, SERVICE_NAME_ERR
	}
	rows, err = c.master.Query(querySQL)
	return
}

func QueryRow(serviceName, queryRowSQL string) (rows *sql.Row, err error) {
	c, ok := clientMap[serviceName]
	if !ok {
		return nil, SERVICE_NAME_ERR
	}
	glog.Debug("queryRowSql:%s", queryRowSQL)
	return c.master.QueryRow(queryRowSQL), nil
}

func Exec(serviceName, execSQL string) (results sql.Result, err error) {
	c, ok := clientMap[serviceName]
	if !ok {
		return nil, SERVICE_NAME_ERR
	}
	glog.Debug("execSql:%s", execSQL)
	results, err = c.master.Exec(execSQL)
	return
}

func GetTableName(tabelNamePrefix string, hashNum uint64) (tableName string, err error) {
	tableName = fmt.Sprintf("%s%d", tabelNamePrefix, hashNum>>54)
	return
}

func QuerySqlCount(serviceName, sql string) (count uint64, err error) {
	mc, ok := clientMap[serviceName]
	if !ok {
		return 0, SERVICE_NAME_ERR
	}
	glog.Debug("QuerySqlCount sql :%s", sql)
	row := mc.master.QueryRow(fmt.Sprintf(`select COUNT(1) as total from (%s) s`, sql))
	glog.Info(fmt.Sprintf(`select COUNT(1) as total from (%s) s`, sql))
	row.Scan(&count)
	return
}

// InitMysqlClient function get mysql_client instance to ensure that the configuration file is initialized
func InitMysqlClient(conf map[string]*MySQLConfig) {
	once.Do(func() {
		clientMap = make(map[string]*MysqlClient)
		for serviceName, msconf := range conf {
			QueryDB := &MysqlClient{}
			if err := QueryDB.initMysqlClient(msconf); err != nil {
				panic(err)
			}
			clientMap[serviceName] = QueryDB
		}

	})
}

func (mc *MysqlClient) initMysqlClient(conf *MySQLConfig) (err error) {
	mc.tableDescMap = conf.TableDesc

	masterHost := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&interpolateParams=true", conf.User, conf.Passwd, conf.Master, conf.Database)
	glog.Debug(masterHost)
	mc.master, err = sql.Open("mysql", masterHost)
	if err != nil {
		glog.Error("init mysql master client failed. err:%s", err.Error())
		return err
	}
	//modify MaxIdelConnNum from 2(default) to 1000
	mc.master.SetMaxIdleConns(1000)

	slaveHost := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&interpolateParams=true", conf.User, conf.Passwd, conf.Slave, conf.Database)
	mc.slave, err = sql.Open("mysql", slaveHost)
	if err != nil {
		glog.Error("init mysql slave client failed. err:%s\n", err.Error())
		return err
	}

	err = mc.master.Ping()
	if err != nil {
		glog.Error("ping master failed. err:%s\n", err.Error())
		return err
	}
	glog.Debug("ping mysql master Ok.")

	return nil
}

// Query exec select peration
func (mc *MysqlClient) Query(querySQL string) (rows *sql.Rows, err error) {
	rows, err = mc.master.Query(querySQL)

	return
}

// QueryRow exec select single row peration
func (mc *MysqlClient) QueryRow(queryRowSQL string) *sql.Row {
	glog.Debug("queryRowSql:%s", queryRowSQL)
	return mc.master.QueryRow(queryRowSQL)
}

// Exec operation
func (mc *MysqlClient) Exec(execSQL string) (results sql.Result, err error) {
	glog.Debug("execSql:%s", execSQL)
	results, err = mc.master.Exec(execSQL)
	return
}

// GetTableName function specific tablename
func (mc *MysqlClient) GetTableName(tabelNamePrefix string, hashNum uint64) (tableName string, err error) {
	tableName = fmt.Sprintf("%s%d", tabelNamePrefix, hashNum>>54)
	return
}

// Close function close the mysql connection.
func (mc *MysqlClient) Close() (err error) {
	mc.master.Close()
	mc.slave.Close()

	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// MysqlEscapeString function similar to the prepare function, the string is isolated and escaped.
func MysqlEscapeString(source string) (string, error) {
	if len(source) == 0 {
		return "", errors.New("source is null")
	}

	j := 0
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j]), nil
}

//check delete one row success or not
func CheckDeleteRowSuccess(result sql.Result) (ok bool, err error) {
	ok = false
	var affectNum int64
	affectNum, err = result.RowsAffected()
	if err != nil {
		glog.Error("get row affect num faild. err:%v", err)
		return
	}
	if affectNum != 1 {
		//glog.Error("affectNum != 1. affectNum:%d", affectNum)
		err = errors.New("user not like or subscribe.")
		return
	}
	ok = true
	return
}
