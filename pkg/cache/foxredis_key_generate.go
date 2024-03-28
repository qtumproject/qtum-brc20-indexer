package cache

import (
	"fmt"
	"strconv"
)

const (
	separator = "::"

	TaskServicePrefix = "TaskService"
	DBPrefix          = "db"

	IsTaskActive            = "IsTaskActive"
	ListRedeemRewardRecord  = "ListRedeemRewardRecord"
	GetTaskListByVersion    = "GetTaskListByVersion"
	GetCustomTaskList       = "GetCustomTaskList"
	taskCreditCount         = "taskCreditCount"
	taskCreditHistoryUnread = "taskCreditHistoryUnread"
	CheckAndGetCustomTask0  = "CheckAndGetCustomTask0"

	Prefix_DB_Coin_Price = "db::coin_price"
	Prefix_DB_Coin_Info  = "db::coin_info"
)

func (redis FoxRedis) ListRedeemRewardRecordCacheKey(did string, rewardItemId, fromTimeStamp int64) string {
	return DBPrefix + separator + ListRedeemRewardRecord + separator +
		did + fmt.Sprintf("::%d::%d", rewardItemId, fromTimeStamp)
}

func (redis FoxRedis) IsTaskActiveCacheKey(taskId int64) string {
	return TaskServicePrefix + separator + IsTaskActive + separator + strconv.Itoa(int(taskId))

}
func (redis FoxRedis) CheckAndGetCustomTask0CacheKey(customTaskId int64) string {
	return TaskServicePrefix + separator + CheckAndGetCustomTask0 + separator + strconv.Itoa(int(customTaskId))
}

func (redis FoxRedis) GetTaskListByVersionCacheKey(version string) string {
	return DBPrefix + separator + GetTaskListByVersion + separator + version
}

func (redis FoxRedis) GetCustomTaskListCacheKey(did string) string {
	return TaskServicePrefix + separator + GetCustomTaskList + separator + did
}

func (redis FoxRedis) CreditCountCacheKey(did string) string {
	return taskCreditCount + separator + did
}
func (redis FoxRedis) TaskCreditHistoryUnreadCacheKey(did string) string {
	return taskCreditHistoryUnread + separator + did
}
func (redis FoxRedis) DBCoinInfoCacheKey(uniqueId string) string {
	return Prefix_DB_Coin_Info + "::" + uniqueId
}
func (redis FoxRedis) DBCoinPriceCacheKey(uniqueId string) string {
	return Prefix_DB_Coin_Price + "::" + uniqueId
}
