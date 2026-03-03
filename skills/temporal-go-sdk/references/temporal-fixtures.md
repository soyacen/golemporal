# Temporal Fixtures Temporal 测试夹具

测试工具和夹具。

```go
import (
    "go.temporal.io/sdk/testsuite"
)

// 测试工作流
func TestMyWorkflow(t *testing.T) {
    suite := &testsuite.WorkflowTestSuite{}
    env := suite.NewTestWorkflowEnvironment()

    // 注册 Activity
    env.RegisterActivity(MyActivity)

    // 模拟 Activity
    env.OnActivity(MyActivity).Return(expectedResult, nil)

    // 执行工作流
    env.ExecuteWorkflow(MyWorkflow, input)

    // 验证结果
    var result string
    err := env.GetWorkflowResult(&result)
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)

    env.AssertExpectations(t)
}

// 使用 Suite
type WorkflowSuite struct {
    suite *testsuite.WorkflowTestSuite
}

func (s *WorkflowSuite) TestWorkflow() {
    env := s.suite.NewTestWorkflowEnvironment()
    env.RegisterActivity(ActivityA)
    env.RegisterActivity(ActivityB)

    // 设置模拟
    env.OnActivity(ActivityA).Return(resultA, nil)

    // 执行
    env.ExecuteWorkflow(MyWorkflow)

    // 验证
    env.AssertExpectations()
}
```
