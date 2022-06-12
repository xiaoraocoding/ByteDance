package config

import (
	"ByteDance/conf"
	"encoding/json"
	"strconv"
)

func Init() {
	conf.Add("app", func() map[string]interface{} {
		return map[string]interface{}{
			"port": conf.Env("PORT", "8888"),
			"name": conf.Env("NAME", "demo"),
		}
	})
	conf.Add("write_sql", func() map[string]interface{} {
		return map[string]interface{}{
			"port":     conf.Env("WRITE_MYSQL_PORT", "3306"),
			"ip":       conf.Env("WRITE_MYSQL_IP", ""),
			"password": conf.Env("WRITE_MYSQL_PASSWORD", "123456"),
		}
	})

	conf.Add("read_sql", func() map[string]interface{} {
		return map[string]interface{}{
			"port":     conf.Env("READ_MYSQL_PORT", "3306"),
			"ip":       conf.Env("READ_MYSQL_IP", ""),
			"password": conf.Env("READ_MYSQL_PASSWORD", "123456"),
		}
	})
}

func GetInterfaceToString(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
