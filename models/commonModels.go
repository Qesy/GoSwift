package models

import (
	"server/lib"

	"github.com/Qesy/qesydb"
)

// 注意，PrimaryField只能用主键 (不然不会更新缓存)
func CacheGetOne(TableName string, PrimaryField string, PrimaryId string) (map[string]string, error) {
	Key := lib.RedisKey_Table_HM_Key(TableName, PrimaryId)
	if Rs, Err := lib.RedisCr.HGetAll(Key); Err == nil && len(Rs) > 0 {
		return Rs, Err
	}
	var m qesydb.Model
	Rs, Err := m.SetTable(TableName).SetWhere(map[string]string{PrimaryField: PrimaryId}).ExecSelectOne()
	if Err == nil && len(Rs) > 0 {
		lib.RedisCr.HMset(Key, Rs)
	}
	return Rs, Err
}

// 清除缓存
func CacheClean(TableName string, PrimaryField string, PrimaryId string) error {
	Key := lib.RedisKey_Table_HM_Key(TableName, PrimaryId)
	return lib.RedisCr.Del(Key)
}

func CacheKey(TableName string, PrimaryId string) string {
	return lib.RedisKey_Table_HM_Key(TableName, PrimaryId)
}

// 设置Key
func CacheSetField(TableName string, PrimaryField string, PrimaryId string, KvMap map[string]string) error {
	Rs, Err := CacheGetOne(TableName, PrimaryField, PrimaryId)
	if len(Rs) == 0 || Err != nil {
		return Err
	}
	for k, v := range KvMap {
		if _, ok := Rs[k]; !ok {
			return Err
		}
		Rs[k] = v
	}
	Key := lib.RedisKey_Table_HM_Key(TableName, PrimaryId)
	lib.RedisCr.HMset(Key, Rs)
	return nil
}

// 获取Key
func CacheGetKeys(TableName string, PrimaryIds []string) []string {
	Keys := []string{}
	for _, v := range PrimaryIds {
		Key := lib.RedisKey_Table_HM_Key(TableName, v)
		Keys = append(Keys, Key)
	}
	return Keys
}

// 雪花算法
func Snowflake() int64 {
	return lib.SnowWorker.Next()
}
