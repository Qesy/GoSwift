package models

import (
	"server/lib"

	"github.com/Qesy/qesydb"
	"github.com/Qesy/qesygo"
)

func SettingGet() (map[string]string, error) {
	Key := lib.RedisKey_Table_HM_Key("dk_system", "Setting")
	if Rs, Err := lib.RedisCr.HGetAll(Key); Err == nil && len(Rs) > 0 {
		return Rs, Err
	}
	var m qesydb.Model
	Arr, Err := m.SetTable("dk_system").ExecSelect()
	if Err == nil && len(Arr) > 0 {
		Rs := qesygo.Array_column_index(Arr, "Content", "Name")
		lib.RedisCr.HMset(Key, Rs)
		return Rs, nil
	}
	return map[string]string{}, Err
}

func SettingGetOne(Key string) string {
	Rs, _ := SettingGet()
	if _, ok := Rs[Key]; !ok {
		return ""
	}
	return Rs[Key]
}

func SettingClean() error {
	Key := lib.RedisKey_Table_HM_Key("dk_system", "Setting")
	return lib.RedisCr.Del(Key)
}
