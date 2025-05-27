package ds

import (
	"context"
	"errors"
	"fmt"
	"github.com/cohesion-org/deepseek-go"
	"io"
	"os"
)

func Message(diff string) (message string, err error) {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")

	client := deepseek.NewClient(apiKey)

	prompt := `你是一个精通各种编程语言的专家，你的任务是根据项目中的 git diff 内容，提取出本次代码变更的核心信息，并生成一条符合 Conventional Commits 规范的 git commit message。

提交信息格式必须为：<emoji> <type>(<scope>): <subject>。你必须只输出一行标准的 commit message，不能输出任何解释或格式说明。

下面是所有允许使用的 type 及其对应 emoji：
- feat (✨)：新功能（feature）
- fix (🐛)：修复 bug（bugfix）
- docs (📝)：文档变更
- style (💅)：不影响代码功能的格式改动（如空格、缩进、分号）
- refactor (🔨)：重构代码，功能无变化
- perf (⚡)：性能优化
- test (✅)：添加或修改测试代码
- chore (🔧)：构建工具或辅助工具的变动（不涉及业务逻辑）
- ci (👷)：持续集成配置变更（如 GitHub Actions）
- revert (⏪)：回滚之前的提交
- build (🛠️)：构建系统或外部依赖的改动
- security (🔐)：安全相关的修复或增强
- ux (🎨)：用户体验相关优化
- i18n (🌐)：国际化/多语言相关变更
- accessibility (♿)：可访问性支持增强（如键盘操作、屏幕阅读器等）
- deps (📦)：依赖库的新增、删除或升级

<scope> 表示受影响的模块或功能域，必须尽量从 diff 中分析出来，如 auth、api、user、payment、cache、i18n 等等。如果实在无法判断，也可以省略，但优先要写出来。

<subject> 是一句话的英文描述，描述清楚这次改动的目的。必须使用英文，标点符号为半角符号，内容简洁但富有表现力，长度控制在 120~150 个字符之间，避免使用 "update code" 或 "fix bug" 这种模糊表达，要体现出你对 diff 逻辑的理解。

请你仔细分析 git diff 中的变化：
- 前缀为 - 的是删除的代码
- 前缀为 + 的是新增的代码
- 分析这些新增或删除的代码的功能意图，理解它们实现了什么功能或修改了什么逻辑
- 判断这是功能增加、修复 bug、性能优化、代码重构等哪种类型

请用文采清晰、生动地表达改动的目的，不要全是技术名词堆砌。最终必须只输出一行 commit message，格式为 <emoji> <type>(<scope>): <subject>，其他任何文字都不要输出。`

	content := fmt.Sprintf("%s git diff的内容如下：%s", prompt, diff)

	request := &deepseek.StreamChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: content},
		},
		Stream: true,
	}
	ctx := context.Background()

	stream, err := client.CreateChatCompletionStream(ctx, request)
	if err != nil {
		return
	}
	var fullMessage string
	defer func() {
		_ = stream.Close()
	}()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			break
		}
		for _, choice := range response.Choices {
			fullMessage += choice.Delta.Content
			fmt.Print(choice.Delta.Content)
		}
	}

	message = fullMessage
	return
}
