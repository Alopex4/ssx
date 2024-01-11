package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"

	//"github.com/vimiix/ssx/cmd/ssx/cmd"
	"ssx/cmd/ssx/cmd"
	//"github.com/vimiix/ssx/internal/cleaner"
	"ssx/ssx/cleaner"
)

func main() {
	var (
		exitCode = 0
	)

	// 创建带有信息通知的上下文
	// 创建通知上下问，接收用户发送ctrl+C (syscall.SIGINT 程序中断), kill <pid> (syscall.SIGTERM 优雅退出 回收资源)
	//func NotifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, stop context.CancelFunc)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	//
	defer cancel()

	// cmd.NewRoot() --> 调用cobra， 需要查阅其使用方式
	// 一般使用cmd.Execute(), 此处使用了ExecuteContext的方法
	// 真正的程序入口
	if err := cmd.NewRoot().ExecuteContext(ctx); err != nil {
		// fatih/color 库调用
		fmt.Println(color.HiRedString(err.Error()))
		exitCode = 1
	}

	// callback func clean
	cleaner.Clean()
	// program exit
	os.Exit(exitCode)
}
