package database

import re "github.com/redis/go-redis/v9"

func ConnectRedis(addr, pwd string, db int) *re.Client {
	return re.NewClient(&re.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db, // use default DB
	})
}
