package trading

import (
	"fmt"
	"github.com/gin-gonic/gin"
	mysqlConf "github.com/liuyang/ly/conf"
	"github.com/liuyang/ly/model/dbTable"
	"log"
	"os"
)

type Aa struct {
	dbTable.TradingFomoLog
	User dbTable.User
}

func Cedan(c *gin.Context) []Aa {
	// 查询NFT
	file, _ := os.OpenFile("/Users/liuyang/Work/go/ly/info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	log.SetOutput(file)
	log.Print("exception")
	var tradesFomoLog []dbTable.TradingFomoLog
	if err := mysqlConf.GetDB().Debug().Limit(30).Find(&tradesFomoLog).Error; err != nil {
		log.Printf("查询 TradingFomoLog 出错: %v\n", err)
	}
	fmt.Println(tradesFomoLog)
	ab := []Aa{}
	if len(tradesFomoLog) > 0 {
		creatorAddrMap := make(map[string]dbTable.TradingFomoLog)
		userAddrSet := make(map[string]struct{})

		for _, nft := range tradesFomoLog {
			creatorAddrMap[nft.User] = nft
			userAddrSet[nft.User] = struct{}{}
		}

		// 获取拥有者信息
		userAddr := make([]string, 0, len(userAddrSet))
		for addr := range userAddrSet {
			userAddr = append(userAddr, addr)
		}

		var users []dbTable.User
		mysqlConf.GetDB().Where("addr IN (?)", userAddr).Find(&users)

		userMap := make(map[string]dbTable.User)
		for _, user := range users {
			userMap[user.Addr] = user
		}
		for _, item := range tradesFomoLog {
			ab = append(ab, Aa{TradingFomoLog: item})
			if userInfo, exists := userMap[item.User]; exists {
				ab = append(ab, Aa{User: userInfo})
			}
		}
	}

	return ab
}
