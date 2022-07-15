package infrastructure

import (
	redisgo "github.com/gomodule/redigo/redis"
	"github.com/opensourceways/xihe-server/infrastructure/redis"
)

const (
	ProjectLikeCountTemp string = "project:likeCount"
	ProjectDownloadTemp  string = "project:download"
)

type ProjectRedis struct {
	redisPool *redisgo.Pool
}

func (p ProjectRedis) GetLikeCount(project_id interface{}) (likecount int, err error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	likecountTemp, redisErr := redisgo.Ints(redis.HMGet(ProjectLikeCountTemp, project_id, conn))
	if redisErr != nil {
		return 0, redisErr
	}
	if len(likecountTemp) < 1 {
		return 0, nil
	}
	likecount = likecountTemp[0]
	return
}
func (p ProjectRedis) GetDownloads(project_id interface{}) (downloads int, err error) {
	conn := p.redisPool.Get()
	defer conn.Close()
	downloadsTemp, redisErr := redisgo.Ints(redis.HMGet(ProjectDownloadTemp, project_id, conn))
	if redisErr != nil {
		return 0, redisErr
	}
	if len(downloadsTemp) < 1 {
		return 0, nil
	}
	downloads = downloadsTemp[0]
	return
}
