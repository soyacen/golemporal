# Multi History Replay 多历史回放

测试工作流在多个历史事件下的行为。

```go
import "go.temporal.io/sdk/testsuite"

// 测试工具
type WorkflowTestSuite struct {
    suite *testsuite.WorkflowTestSuite
}

func (s *WorkflowTestSuite) TestWithMultipleHistories() {
    // 模拟不同的历史事件序列
    histories := [][]*history.Event{
        createHistory1(), // 场景1: 成功完成
        createHistory2(), // 场景2: 失败重试
        createHistory3(), // 场景3: 取消
    }

    for i, history := range histories {
        env := s.suite.NewTestWorkflowEnvironment()

        // 加载特定历史
        env.SetHistory(history)

        // 执行工作流
        result, err := env.ExecuteWorkflow(MyWorkflow, input)

        // 验证结果
        assert.NoError(t, err)
        assert.Equal(t, expectedResult, result)

        env.AssertExpectations(t)
    }
}

func createHistory1() []*history.Event {
    // 构建成功场景的历史事件
}
```
