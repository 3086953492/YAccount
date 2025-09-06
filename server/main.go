package main

import (
	"YAccount/initialize"
	"fmt"
	"github.com/3086953492/YaBase/global"
	"github.com/3086953492/YaBase/logger"
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	if err := global.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	cfg := global.GetGlobalConfig()

	// 初始化日志
	if err := logger.InitWithConfig(cfg.Log); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化数据库
	if err := initialize.InitDB(); err != nil {
		logger.Error("初始化数据库失败", zap.Error(err))
		return
	}

	// 初始化验证器
	if err := initialize.InitValidator(); err != nil {
		logger.Error("初始化验证器失败", zap.Error(err))
		return
	}

	// 初始化 Redis
	if err := global.InitRedisWithConfig(cfg.Redis); err != nil {
		logger.Error("初始化 Redis 失败", zap.Error(err))
		return
	}

	// 初始化缓存
	if err := global.InitCache(); err != nil {
		logger.Error("初始化缓存失败", zap.Error(err))
		return
	}

	// 获取端口号，优先使用命令行参数
	port := cfg.Server.Port
	if len(os.Args) > 1 {
		argPort, err := strconv.Atoi(os.Args[1])
		if err == nil {
			port = argPort
		}
	}

	// 初始化 Gin 路由
	r := initialize.InitRouters()

	logger.Info("YAccount 服务启动成功")
	logger.Info(fmt.Sprintf("服务正在运行在端口 %d", port))

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Error("启动服务失败", zap.Error(err))
	}
}
