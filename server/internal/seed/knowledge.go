package seed

import (
	"context"
	"fmt"

	"github.com/cng1985/ai-learning-server/internal/service"
)

// IndexKnowledge 启动时确保知识库向量索引已构建
func IndexKnowledge(knowledge *service.KnowledgeService) error {
	fmt.Println("📚 检查知识库向量索引...")
	if err := knowledge.EnsureIndexed(context.Background()); err != nil {
		fmt.Printf("⚠️  知识库索引构建失败: %v\n", err)
		return nil // 不阻塞服务启动
	}
	status, _ := knowledge.Status()
	if status != nil {
		fmt.Printf("   知识库: %d 个文本块, 嵌入模型 %s (%s)\n", status.ChunkCount, status.EmbedModel, status.EmbedSource)
	}
	return nil
}
