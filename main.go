package main

import (
	"fmt"
	"gacm/ds"
	"gacm/git"
	"log"
)

func main() {
	diff, err := git.GetGitDiff()
	if err != nil {
		log.Fatalf("获取 git diff 失败: %v", err)
		return
	}

	if diff == "" {
		log.Println("没有检测到暂存的更改。")
		return
	}

	msg, err := ds.Message(diff)
	if err != nil {
		log.Println("生成提交信息失败:", err)
		return
	}

	err = git.CommitChanges(msg)
	if err != nil {
		log.Println("提交commit message失败:", err)
		return
	}

	fmt.Println("\n提交成功！")
}
