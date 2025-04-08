package trading

import (
	"fmt"
	"github.com/gin-gonic/gin"
	mysqlConf "github.com/liuyang/ly/conf"
	"github.com/liuyang/ly/model/dbTable"
	"log"
	"sync"
	time2 "time"
)

func List(c *gin.Context) ([]dbTable.TradingRound, error) {
	var tradesTound []dbTable.TradingRound
	currentTime := time2.Now().Unix() // 获取当前时间
	fmt.Println(currentTime)
	result := mysqlConf.GetDB().
		Where("`status` = ?", 1).
		Where("start_ts <= ?", currentTime).
		Where("end_ts <= ?", currentTime).
		Where("types = ?", 1).
		Order("start_ts ASC").
		Limit(10).
		Find(&tradesTound)

	if result.Error != nil {
		return nil, result.Error
	}

	var wg sync.WaitGroup        // 创建 WaitGroup
	wg.Add(len(tradesTound) * 2) // 每次循环添加两个任务

	for k, _ := range tradesTound {
		go sendReward(&tradesTound[k], &wg)
		go againOpen(&tradesTound[k], &wg)
	}
	wg.Wait()

	return tradesTound, nil
}

// 合并后的数据结构
type MergedData struct {
	dbTable.TradingFomoLog
	dbTable.User
}

func Info(c *gin.Context) ([]MergedData, error) {
	tradesFomoLogChan := make(chan []dbTable.TradingFomoLog, 1)
	userChan := make(chan []dbTable.User, 1)
	var wg sync.WaitGroup
	log.Println("开始获取数据")

	// 查询 TradingFomoLog 的 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		var tradesFomoLog []dbTable.TradingFomoLog
		if err := mysqlConf.GetDB().Debug().Limit(30).Find(&tradesFomoLog).Error; err != nil {
			log.Printf("查询 TradingFomoLog 出错: %v\n", err)
			tradesFomoLogChan <- nil
			return
		}
		tradesFomoLogChan <- tradesFomoLog
	}()

	// 查询 User 的 goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		var users []dbTable.User
		if err := mysqlConf.GetDB().Limit(30).Find(&users).Error; err != nil {
			log.Printf("查询 User 出错: %v\n", err)
			userChan <- nil
			return
		}
		userChan <- users
	}()

	// 等待所有 goroutine 完成并关闭通道
	wg.Wait()
	close(tradesFomoLogChan)
	close(userChan)

	// 接收查询结果
	tradesFomoLogs := <-tradesFomoLogChan
	users := <-userChan

	// 检查数据是否为空
	if tradesFomoLogs == nil || users == nil {
		return nil, fmt.Errorf("未能获取到完整数据")
	}

	// 合并数据
	var mergedData []MergedData
	for _, trade := range tradesFomoLogs {
		for _, user := range users {
			if trade.User == user.Addr {
				mergedData = append(mergedData, MergedData{
					TradingFomoLog: trade,
					User:           user,
				})
				break // 假设每个交易记录只有一个匹配的用户，避免重复遍历
			}
		}
	}
	log.Printf("合并后的数据数量: %d\n", len(mergedData))
	return mergedData, nil
}

func sendReward(val *dbTable.TradingRound, wg *sync.WaitGroup) {
	defer wg.Done() // 完成后减少计数器
	val.ID = 123
	// 这里添加发奖的逻辑
	fmt.Printf("Sending reward for trading round ID: %d\n", val.ID)
}

func againOpen(val *dbTable.TradingRound, wg *sync.WaitGroup) {
	defer wg.Done() // 完成后减少计数器
	val.ID = 234
	// 这里添加重新开启活动的逻辑
	fmt.Printf("Reopening trading round ID: %d\n", val.ID)
}
