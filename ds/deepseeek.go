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

	prompt := "你是一个精通各种编程语言的专家，能在项目中的git diff信息中，熟练的提炼出git diff中的具体信息。同时你要以Conventional Commits这种规范来输出git commit message。Conventional Commits 是一种格式规范，要求 Git 提交信息遵循特定格式：<type>(<scope>): <subject>，其中：<type>：提交的类型（比如是修复 bug 还是加新功能）、<scope>：可选，说明改动影响的模块、<subject>：一句话描述这次改动的内容。<type> 类型说明；feat是新功能（feature）、fix是修复 bug（bugfix）、docs是文档改动、style是不影响功能的改动（空格、格式、少个分号等）、refactor是重构代码（没有功能变化）、perf是性能优化、test是增加测试、chore是构建过程或辅助工具的变动、ci是CI 配置相关（如 GitHub Actions）、revert是回滚某次提交。type只能在上面的列举的类型里挑选，如果都不合适，就选择feat。<scope>：说明改动影响的模块，如果分析不出来，可以为空。<subject>：一句话描述这次改动的内容，这个必须分析出来，但是要简短，保持在120个单词之内。如果分析出了<scope>就是 这样feat(product): send an email to the customer when a product is shipped.没有<scope>就是feat: send an email to the customer when a product is shipped.要求必须输出英文，标点符号为半角符号。<type>用emoji来标识，对应关系是：feat是✨, fix是🐛，docs是📝，style是💅，refactor是🔨，perf是⚡，test是✅，chore是🔧，ci是👷 ，revert是⏪ 。<type>要尽量分析出来，而不是只选择了feat，<scope> 也要尽量分析出来。你输出的信息要有文采，不要经常就是计算机和代码类的名词，另外你要千万的仔细检查一下，根据git diff的内容来仔细的看，哪些代码是新增的，哪些代码是删除的，前面标识了-号的是删除的代码，标识了+号的是增加的代码，你要仔细的分析新增和删除的代码是在构建什么样的逻辑关系。输出长度限制：输出的<subject>的要求要保持在150个字符以内，最好在120个字符左右，不长不短。。只输出这个一个commit message的字符串，不要输出其他任何内容。"

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
