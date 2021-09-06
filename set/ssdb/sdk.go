package ssdb

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/seefan/goerr"
	"github.com/seefan/gossdb/v2/pool"
)

const (
	oK       string = "ok"
	notFound string = "not_found"
)

func zSet(c *pool.Client, key, item string) (count int, err error) {
	resp, err := c.Do("zset", key, item, 1)
	if err != nil {
		return 0, goerr.Errorf(err, "Zset %s/%s error", key, item)
	}

	if len(resp) == 0 || resp[0] != oK {
		return 0, makeError(resp, key, item)
	}
	if len(resp) >= 2 {
		count, _ = strconv.Atoi(resp[1])
	}
	return count, nil
}
func multiZSet(c *pool.Client, key string, items ...string) (count int, err error) {
	args := []interface{}{"multi_zset", key}
	for _, item := range items {
		args = append(args, item, "1")
	}
	resp, err := c.Do(args...)

	if err != nil {
		return 0, goerr.Errorf(err, "MultiZset %s %v error", key, items)
	}
	if len(resp) == 0 || resp[0] != oK {
		return 0, makeError(resp, key, items)
	}
	if len(resp) >= 2 {
		count, _ = strconv.Atoi(resp[1])
	}
	return count, nil
}
func zDel(c *pool.Client, key, item string) (count int, err error) {
	resp, err := c.Do("zdel", key, item)
	if err != nil {
		return 0, goerr.Errorf(err, "Zdel %s/%s error", key, item)
	}
	if len(resp) == 0 || resp[0] != oK {
		return 0, makeError(resp, key, item)
	}
	if len(resp) >= 2 {
		count, _ = strconv.Atoi(resp[1])
	}
	return count, nil
}
func multiZDel(c *pool.Client, key string, items ...string) (count int, err error) {
	if len(items) == 0 {
		return 0, nil
	}
	args := []interface{}{"multi_zdel", key}
	for _, v := range items {
		args = append(args, v)
	}
	resp, err := c.Do(args...)
	if err != nil {
		return 0, goerr.Errorf(err, "MultiZdel %s %s error", key, items)
	}
	if len(resp) == 0 || resp[0] != oK {
		return 0, makeError(resp, key, items)
	}
	if len(resp) >= 2 {
		count, _ = strconv.Atoi(resp[1])
	}
	return count, nil
}

func makeError(resp []string, errKey ...interface{}) error {
	if len(resp) < 1 {
		return errors.New("ssdb response error")
	}
	// 正常返回的不存在不报错，如果要捕捉这个问题请使用exists
	if resp[0] == notFound {
		return nil
	}
	if len(errKey) > 0 {
		return fmt.Errorf("access ssdb error, code is %v, parameter is %v", resp, errKey)
	}
	return fmt.Errorf("access ssdb error, code is %v", resp)
}
