package main

import (
	"database/sql"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 表名和批量大小配置
const (
	BatchSize  = 1000  // 每次插入的批量大小
	TotalRows  = 10000 // 总共要插入的数据量
	DBUser     = "root"
	DBPassword = "rOYkHEc#jOesowLL" // 修改为你的数据库密码
	DBHost     = "47.243.51.88"
	DBPort     = "3306"
	DBName     = "dongshenggj"                        // 修改为你的数据库名
	TableName  = "yz_thirdparty_game_record_5billion" // 表名

)

// GameRecord 数据结构与表对应
type GameRecord struct {
	MemberID         uint
	Member           string
	MemberType       uint8
	ParentID         uint
	ParentIDs        string
	ThirdPartyID     uint
	ThirdPartyName   string
	GameID           uint
	GameName         string
	GameType         uint8
	BetDetail        string
	GameOrder        string
	GameCode         string
	PlayType         string
	BetMoney         float64
	BetAmount        float64
	WinBonus         float64
	CommissionAmount float64
	Balance          float64
	BetTime          uint
	SlotID           uint
	SlotName         string
	RebateAmount     float64
	RoomType         string
	TableName        string
	Status           int8
	TCStatus         uint8
	GameTime         uint
	Source           string
	CreateTime       uint
	ThirdPlayName    string
	GameResult       string
	Terminal         string
	IsRepair         bool
	WinLossAmount    float64
	CommissionStatus int8
	IP               string
	SettlementAmount float64
	IsRelateProxy    bool
	AllReady         string
	Odds             string
	WinningNumber    string
	BelongParentID   int
}

func randomMixedString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789") // 英文字符集
	chineseChars := []rune("王儿测试游戏名称数据插入字段混合验证成功示例开发拳皇98")                            // 常用中文字符

	result := make([]rune, length)
	for i := 0; i < length; i++ {
		if rand.Intn(2) == 0 {
			result[i] = chars[rand.Intn(len(chars))] // 50% 生成英文字符
		} else {
			result[i] = chineseChars[rand.Intn(len(chineseChars))] // 50% 生成中文字符
		}
	}
	return string(result)
}

// 生成虚假数据
func generateFakeData() GameRecord {
	return GameRecord{
		MemberID:         uint(rand.Intn(1000000)),
		Member:           randomMixedString(10), // 混合中文和英文
		MemberType:       3,
		ParentID:         uint(rand.Intn(100000)),
		ParentIDs:        faker.UUIDDigit(),
		ThirdPartyID:     uint(rand.Intn(100)),
		ThirdPartyName:   randomMixedString(8), // 混合中文和英文
		GameID:           uint(rand.Intn(100)),
		GameName:         randomMixedString(6), // 混合中文和英文
		GameType:         uint8(rand.Intn(8)),
		BetDetail:        faker.Sentence(),
		GameOrder:        faker.UUIDDigit(),
		GameCode:         faker.UUIDDigit(),
		PlayType:         faker.Word(),
		BetMoney:         rand.Float64() * 1000,
		BetAmount:        rand.Float64() * 1000,
		WinBonus:         rand.Float64() * 1000,
		CommissionAmount: rand.Float64() * 100,
		Balance:          rand.Float64() * 1000,
		BetTime:          uint(time.Now().Unix()),
		SlotID:           uint(rand.Intn(1000)),
		SlotName:         randomMixedString(8), // 混合中文和英文
		RebateAmount:     rand.Float64() * 100,
		RoomType:         "高级房",
		TableName:        fmt.Sprintf("T-%d", rand.Intn(100)),
		Status:           1,
		TCStatus:         0,
		GameTime:         uint(time.Now().Unix()),
		Source:           "{}",
		CreateTime:       uint(time.Now().Unix()),
		ThirdPlayName:    faker.Word(),
		GameResult:       "Win",
		Terminal:         "PC",
		IsRepair:         false,
		WinLossAmount:    rand.Float64() * 1000,
		CommissionStatus: 0,
		IP:               faker.IPv4(),
		SettlementAmount: rand.Float64() * 1000,
		IsRelateProxy:    false,
		AllReady:         "N",
		Odds:             "1.5",
		WinningNumber:    faker.UUIDDigit(),
		BelongParentID:   rand.Intn(10000),
	}
}

func main() {
	// 初始化数据库连接
	//dsn := "root:rOYkHEc#jOesowLL@tcp(47.243.51.88:3306)/dongshenggj?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:rOYkHEc#jOesowLL@tcp(47.243.51.88:3306)/dongshenggj"
	// 数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)
	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	// 验证一条简单的sql 语句
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatalf("查询 MYSQL 数据库版本失败: %v", err)
	}
	fmt.Println("MYSQL数据库版本:", version)

	// golang 编辑需要配置数据库连接 内省
	insertSQL := fmt.Sprintf("\n\t\tINSERT INTO yz_thirdparty_game_record_5billion\n\t(memberid, member, member_type, parent_id, parent_ids, third_party_id, third_party_name, gameid, game_name, game_type, bet_detail, game_order, game_code, play_type, bet_money, bet_amount, win_bonus, commission_amount, balance, bet_time, slot_id, slot_name, rebate_amount, room_type,table_name, status, tc_status, game_time, source, create_time, third_play_name, game_result, terminal, is_repair, win_loss_amount, commission_status, ip, settlement_amount, is_relate_proxy, allready, odds, winning_number, belong_parent_id)\n\tVALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	start := time.Now()
	for i := 0; i < TotalRows/BatchSize; i++ {
		// 开始事务
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}
		stmt, err := tx.Prepare(insertSQL)
		if err != nil {
			tx.Rollback() // 出错时回滚事务
			log.Fatalf("Failed to prepare statement: %v", err)
		}
		for j := 0; j < BatchSize; j++ {
			data := generateFakeData()
			_, err := stmt.Exec(
				data.MemberID, data.Member, data.MemberType, data.ParentID, data.ParentIDs, data.ThirdPartyID,
				data.ThirdPartyName, data.GameID, data.GameName, data.GameType, data.BetDetail, data.GameOrder,
				data.GameCode, data.PlayType, data.BetMoney, data.BetAmount, data.WinBonus, data.CommissionAmount,
				data.Balance, data.BetTime, data.SlotID, data.SlotName, data.RebateAmount, data.RoomType, data.TableName,
				data.Status, data.TCStatus, data.GameTime, data.Source, data.CreateTime, data.ThirdPlayName, data.GameResult,
				data.Terminal, data.IsRepair, data.WinLossAmount, data.CommissionStatus, data.IP, data.SettlementAmount,
				data.IsRelateProxy, data.AllReady, data.Odds, data.WinningNumber, data.BelongParentID,
			)
			if err != nil {
				tx.Rollback() // 出错时回滚事务
				log.Fatalf("Insert failed in batch %d: %v", i+1, err)
			}
		}
		// 关闭语句并检查错误
		if err := stmt.Close(); err != nil {
			tx.Rollback()
			log.Fatalf("Failed to close statement: %v", err)
		}

		// 提交事务并检查错误
		if err := tx.Commit(); err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		} else {
			fmt.Printf("Batch %d committed successfully\n", i+1)
		}

		fmt.Printf("Inserted batch %d\n", i+1)
	}
	fmt.Printf("数据生成完成 用时: %v\n", time.Since(start))
}
