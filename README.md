# agent-ai
结构化的问题处理：
根据问题类型分发到不同的处理函数
支持存钱计算和薪水比较两种复杂场景
其他复杂问题使用 agent chain 处理
多步骤处理：
每个复杂问题都被分解为多个步骤
每个步骤都有明确的错误处理
显示中间结果便于调试
灵活的工具使用：
根据需要组合使用不同的工具
支持工具之间的结果传递
保持了简单问题的直接处理方式