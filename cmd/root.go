package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var ctx = context.Background()

func Execute() {

	var rootCmd = &cobra.Command{
		Use:   "deleteRedisKeys  -a server:6379 -p password -k xxxKey*  ",
		Short: "删除redisKey",
		Long:  "批量删除redisKey,支持模糊匹配",

		Run: func(cmd *cobra.Command, args []string) {
			address, _ := cmd.Flags().GetString("address")
			password, _ := cmd.Flags().GetString("password")
			db, _ := cmd.Flags().GetInt("db")
			key, _ := cmd.Flags().GetString("key")

			if address == "" || password == "" || key == "" {
				fmt.Println("输入错误,请使用 --help 查看支持命令")
				return
			}

			redis_client := createClient(address, password, db)
			defer redis_client.Close()

			delKeys(redis_client, key)

			fmt.Print("处理完成 ....\n")
		},
	}

	rootCmd.Flags().StringP("address", "a", "", "redis 服务器地址 如:127.0.0.1:6379")
	rootCmd.Flags().StringP("password", "p", "", "redis 密码")
	rootCmd.Flags().IntP("db", "d", 0, "redis 数据库 0")
	rootCmd.Flags().StringP("key", "k", "", "需要删除的redisKey 如：xxxx:xxx*")
	rootCmd.PersistentFlags().StringP("loglevel", "l", "info", "设置日志级别. 支持: debug, info, warn, error, fatal")

	rootCmd.DisableAutoGenTag = true

	err := doc.GenMarkdownTree(rootCmd, "./")
	if err != nil {
		log.Fatal(err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//创建连接
func createClient(addr string, password string, db int) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping(ctx).Result()

	if err != nil {
		fmt.Println("redis 连接失败:", pong, err)
		panic("reidis服务器 连接失败")
	}

	fmt.Println("reidis服务器 连接成功", pong)

	return client
}

func delKeys(client *redis.Client, key string) {

	keys, err2 := client.Keys(ctx, key).Result()

	if len(keys) == 0 {
		fmt.Printf("找不到要删除的Key%s \n", key)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("共找到 %d 条要删除的数据,请确认是否要删除？删除请输入Y,按任意键退出 ", len(keys))
	confirm, _ := reader.ReadString('\n')

	if strings.TrimSpace(confirm) != "y" {
		fmt.Printf("输入%s拒绝执行，正在退出...\n ", strings.ToLower(confirm))
		return
	}

	if err2 != nil {
		fmt.Println("查找redis key 失败\n", err2)
	}

	var (
		begin = 0         //起始索引
		end   = len(keys) //结束索引
		num   = 2         //每批数量
	)

	batchs := []string{} // 临时变量

FOR: //循环

	for i := begin; i < end; i++ {
		batchs = append(batchs, keys[i]) //向数组中追加元素

		if i > 0 && i%num == 0 {
			begin = i + 1
			goto EXEC
		} else {
			if end-begin < num {
				begin = end
			}
		}
	}

EXEC: //处理

	fmt.Printf("当前第%d批,%v\n", begin, batchs)

	result, err3 := client.Del(ctx, batchs...).Result()
	if err3 != nil {
		fmt.Println("批量删除redis key 失败\n", err3)
	}
	fmt.Printf("批量删除redis成功 返回%d \n", result)

	batchs = []string{}
	if begin < end {
		begin += 1
		goto FOR

	}

}

func setOperation(client *redis.Client) {

	redis_key := "redis_name"

	err := client.Set(ctx, redis_key, "123456", time.Duration(5)*time.Minute)
	if err != nil {
		fmt.Println("设置redis key 失败\n", err)
	}

	val, err2 := client.Get(ctx, redis_key).Result()
	if err2 != nil {
		fmt.Println("获取redis value 失败\n", err2)
	}

	fmt.Printf("获取Redis %s=%s \n", redis_key, val)

	result, err3 := client.Del(ctx, redis_key).Result()
	if err3 != nil {
		fmt.Println("删除redis key 失败\n", err3)
	}
	fmt.Printf("删除redis成功 返回%d \n", result)
}
