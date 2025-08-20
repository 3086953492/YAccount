package initialize

import (
	"YAccount/configs"
	"YAccount/global"
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitConfig() error {
	// 处理命令行参数
	var config string
	flag.StringVar(&config, "c", "", "指定配置文件路径")
	flag.Parse()

	// 如果用户通过命令行或环境变量指定了配置文件，使用指定的配置文件
	if config == "" { // 判断命令行参数是否指定了配置文件
		if configEnv := os.Getenv("CONFIG"); configEnv != "" { // 读取环境变量
			config = configEnv
			fmt.Printf("您正在使用环境变量, 配置文件的路径为%s\n", configEnv)
		} else if gin.Mode() != "" {
			// 根据gin模式自动选择配置文件
			switch gin.Mode() {
			case gin.DebugMode:
				config = "./configs/config.yaml"
			case gin.TestMode:
				config = "./configs/config.test.yaml"
			case gin.ReleaseMode:
				config = "./configs/config.release.yaml"
			}
			if config != "" {
				fmt.Printf("您正在使用gin模式的%s环境名称, 配置文件的路径为%s\n", gin.Mode(), config)
			}
		}
	} else {
		fmt.Printf("您正在使用命令行的-c参数传递的值, 配置文件的路径为%s\n", config)
	}

	// 如果指定了配置文件，直接使用指定的配置文件
	if config != "" {
		viper.SetConfigFile(config)
		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("配置文件读取失败: %v", err)
		}

		// 监听配置文件变化
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			fmt.Println("配置文件发生变更: ", in.Name)
			err = viper.Unmarshal(global.Cfg)
			if err != nil {
				fmt.Printf("配置文件重新加载失败: %v\n", err)
			} else {
				fmt.Println("配置文件已重新加载")
			}
		})

		global.Cfg = &configs.Config{}
		err = viper.Unmarshal(global.Cfg)
		if err != nil {
			return fmt.Errorf("配置文件解析失败: %v", err)
		}
		fmt.Printf("成功读取配置文件: %s\n", config)
		return nil
	}

	// 读取默认位置的配置文件
	viper.AddConfigPath("./configs")

	// 优先读取YAML配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {

		// 如果YAML读取失败，尝试读取JSON配置
		fmt.Printf("YAML配置读取失败，尝试读取JSON配置: %v\n", err)

		viper.SetConfigName("config")
		viper.SetConfigType("json")

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("配置文件读取失败，YAML和JSON配置都无法读取: %v", err)
		}

		fmt.Println("成功读取JSON配置文件")
	} else {
		fmt.Println("成功读取YAML配置文件")
	}

	global.Cfg = &configs.Config{}
	return viper.Unmarshal(global.Cfg)
}
