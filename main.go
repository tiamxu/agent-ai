package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

// 计算工具
type CalculatorTool struct{}

func (c *CalculatorTool) Name() string {
	return "calculator"
}

func (c *CalculatorTool) Description() string {
	return "用于处理数学运算问题的工具"
}

func (c *CalculatorTool) Call(ctx context.Context, query string) (string, error) {
	// 提取数学表达式
	fmt.Println("####:", query)
	query = strings.ToLower(query)
	query = strings.ReplaceAll(query, "计算", "")
	query = strings.ReplaceAll(query, "的值", "")
	query = strings.TrimSpace(query)

	result, err := evalMathExpression(query)
	if err != nil {
		return "", fmt.Errorf("计算错误: %w", err)
	}
	return fmt.Sprintf("计算结果: %v", result), nil
}

// 计算数学表达式
func evalMathExpression(expr string) (float64, error) {
	// 移除空格并分割表达式
	expr = strings.ReplaceAll(expr, " ", "")
	var num1, num2 float64
	var err error

	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		if len(parts) != 2 {
			return 0, fmt.Errorf("无效的加法表达式")
		}
		num1, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第一个数字: %w", err)
		}
		num2, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第二个数字: %w", err)
		}
		return num1 + num2, nil
	} else if strings.Contains(expr, "-") {
		parts := strings.Split(expr, "-")
		if len(parts) != 2 {
			return 0, fmt.Errorf("无效的减法表达式")
		}
		num1, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第一个数字: %w", err)
		}
		num2, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第二个数字: %w", err)
		}
		return num1 - num2, nil
	} else if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) != 2 {
			return 0, fmt.Errorf("无效的乘法表达式")
		}
		num1, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第一个数字: %w", err)
		}
		num2, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第二个数字: %w", err)
		}
		return num1 * num2, nil
	} else if strings.Contains(expr, "/") {
		parts := strings.Split(expr, "/")
		if len(parts) != 2 {
			return 0, fmt.Errorf("无效的除法表达式")
		}
		num1, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第一个数字: %w", err)
		}
		num2, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析第二个数字: %w", err)
		}
		if num2 == 0 {
			return 0, fmt.Errorf("除数不能为零")
		}
		return num1 / num2, nil
	}
	return 0, fmt.Errorf("不支持的运算符")
}

// 假设你有一个模拟的员工数据库查询工具
type EmployeeDatabaseTool struct{}

// 添加必要的接口方法
func (e *EmployeeDatabaseTool) Name() string {
	return "employee_database"
}

func (e *EmployeeDatabaseTool) Description() string {
	return "用于查询员工信息的工具，可以获取员工的职位、年龄和薪资信息"
}

func (e *EmployeeDatabaseTool) Call(ctx context.Context, query string) (string, error) {
	// 模拟查询，假设数据库中有一份员工信息
	fmt.Println("EmployeeDatabaseTool:", query)
	employees := map[string]map[string]string{
		"libai": {
			"position": "运维总监",
			"age":      "39",
			"salary":   "$20k",
		},
		"sunshangxiang": {
			"position": "财务总监",
			"age":      "30",
			"salary":   "$120K",
		},
	}

	// 提取员工名称
	query = strings.ToLower(query)
	var employeeName string
	if strings.Contains(query, "libai") {
		employeeName = "libai"
	} else if strings.Contains(query, "sunshangxiang") {
		employeeName = "sunshangxiang"
	} else {
		return "", fmt.Errorf("未找到员工信息")
	}

	// 通过查询获得特定员工的资料
	employee, exists := employees[employeeName]
	if !exists {
		return "", fmt.Errorf("未找到员工信息")
	}

	// 返回员工的详细信息
	return fmt.Sprintf("Position: %s, Age: %s, Salary: %s", employee["position"], employee["age"], employee["salary"]), nil
}

// 从薪水信息中提取数值
func extractSalaryAmount(salaryInfo string) float64 {
	// 使用正则表达式提取包含 Salary 的部分
	salaryRe := regexp.MustCompile(`Salary: \$(\d+)k`)
	matches := salaryRe.FindStringSubmatch(salaryInfo)
	if len(matches) < 2 {
		log.Printf("无法从信息中提取薪水: %s", salaryInfo)
		return 0
	}

	// 提取数字部分并转换为 float64
	base, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		log.Printf("解析薪水数值失败: %v", err)
		return 0
	}

	// 将 k 转换为实际数值（乘以 1000）
	return base * 1000
}

// 定义工具类型常量
const (
	ToolTypeEmployee   = "employee"
	ToolTypeCalculator = "calculator"
	ToolTypeUnknown    = "unknown"
)

// 工具选择器接口
type ToolSelector interface {
	SelectTool(question string) (tools.Tool, error)
}

// 默认工具选择器
type DefaultToolSelector struct {
	tools map[string]tools.Tool
}

// 创建新的工具选择器
func NewDefaultToolSelector() *DefaultToolSelector {
	return &DefaultToolSelector{
		tools: map[string]tools.Tool{
			ToolTypeEmployee:   &EmployeeDatabaseTool{},
			ToolTypeCalculator: &CalculatorTool{},
		},
	}
}

// 根据问题选择合适的工具
func (s *DefaultToolSelector) SelectTool(question string) (tools.Tool, error) {
	toolType := s.determineToolType(question)
	if tool, exists := s.tools[toolType]; exists {
		return tool, nil
	}
	return nil, fmt.Errorf("未找到合适的工具处理该问题")
}

// 确定工具类型
func (s *DefaultToolSelector) determineToolType(question string) string {
	question = strings.ToLower(question)

	// 员工信息查询关键词
	employeeKeywords := []string{"职位", "薪水", "年龄", "libai", "sunshangxiang"}
	for _, keyword := range employeeKeywords {
		if strings.Contains(question, strings.ToLower(keyword)) {
			return ToolTypeEmployee
		}
	}

	// 计算器关键词
	calculatorOperators := []string{"+", "-", "*", "/"}
	for _, operator := range calculatorOperators {
		if strings.Contains(question, operator) {
			return ToolTypeCalculator
		}
	}

	return ToolTypeUnknown
}

func main() {
	ctx := context.Background()

	llm, err := openai.New(
		openai.WithModel("qwen-turbo"),
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"),
		openai.WithToken("sk-dda0881c077849eea532a185e5731d28"),
	)

	// llm, err := ollama.New(
	// 	ollama.WithModel("deepseek-r1:7b"),
	// 	ollama.WithServerURL("http://192.168.1.228:11434"),
	// )
	if err != nil {
		log.Fatal(err)
	}

	// 定义工具集
	agentTools := []tools.Tool{
		&CalculatorTool{},
		&EmployeeDatabaseTool{},
	}

	// 创建 Agent
	agent := agents.NewOneShotAgent(
		llm,
		agentTools,
		agents.WithMaxIterations(5),
	)
	// 创建执行器
	executor := agents.NewExecutor(agent)

	// 创建工具选择器
	toolSelector := NewDefaultToolSelector()

	// 测试问题
	questions := []string{
		"计算2 + 3的值",   // 简单问题
		"libai的薪水是多少", // 简单问题
		// "比较 libai 和 sunshangxiang 的薪水差异", // 复杂问题
		"如果 libai 每月存一半工资，一年能存多少钱？", // 复杂问题
	}

	for _, question := range questions {
		fmt.Printf("\n处理问题: %s\n", question)
		var answer string
		var err error

		if isComplexQuestion(question) {
			fmt.Println("这是一个复杂问题，使用多步骤处理")
			answer, err = handleComplexQuestion(ctx, question, executor, toolSelector)
		} else {
			fmt.Println("这是一个简单问题，使用工具选择器处理")
			tool, err := toolSelector.SelectTool(question)
			if err != nil {
				log.Printf("选择工具失败: %v\n", err)
				continue
			}
			answer, err = tool.Call(ctx, question)
		}

		if err != nil {
			log.Printf("处理问题失败: %v\n", err)
			continue
		}

		fmt.Printf("答案: %s\n", answer)
	}
}

// 处理复杂问题的函数
func handleComplexQuestion(ctx context.Context, question string, executor *agents.Executor, toolSelector *DefaultToolSelector) (string, error) {
	// 如果是存钱相关的问题
	if strings.Contains(question, "存") && strings.Contains(question, "钱") {
		return handleSavingsQuestion(ctx, question, toolSelector)
	}

	// 如果是比较相关的问题
	if strings.Contains(question, "比较") || strings.Contains(question, "差异") {
		return handleComparisonQuestion(ctx, question, toolSelector)
	}

	// 其他复杂问题使用 agent chain 处理
	return chains.Run(ctx, executor, question)
}

// 处理存钱相关问题
func handleSavingsQuestion(ctx context.Context, question string, toolSelector *DefaultToolSelector) (string, error) {
	// 1. 首先查询员工薪水
	employeeTool := &EmployeeDatabaseTool{}
	salaryInfo, err := employeeTool.Call(ctx, "libai的薪水")
	if err != nil {
		return "", fmt.Errorf("查询薪水失败: %w", err)
	}
	fmt.Printf("薪水信息: %s\n", salaryInfo)

	// 2. 从薪水信息中提取数值
	salary := extractSalaryAmount(salaryInfo)
	if salary == 0 {
		return "", fmt.Errorf("无法获取有效的薪水数值")
	}
	fmt.Printf("提取的月薪数值: %.2f\n", salary)

	// 3. 计算年度存款
	calculatorTool := &CalculatorTool{}
	calculationQuery := fmt.Sprintf("%.0f+%.0f", salary*0.5, salary*0.5*11)
	result, err := calculatorTool.Call(ctx, calculationQuery)
	if err != nil {
		return "", fmt.Errorf("计算失败: %w", err)
	}

	return fmt.Sprintf("基于月薪 %.2f 元，如果每月存一半（%.2f 元），一年可以存 %s",
		salary, salary*0.5, result), nil
}

// 处理比较相关问题
func handleComparisonQuestion(ctx context.Context, question string, toolSelector *DefaultToolSelector) (string, error) {
	// 1. 查询第一个员工的薪水
	employeeTool := &EmployeeDatabaseTool{}
	libaiInfo, err := employeeTool.Call(ctx, "libai的薪水")
	if err != nil {
		return "", fmt.Errorf("查询libai薪水失败: %w", err)
	}

	// 2. 查询第二个员工的薪水
	sunshangxiangInfo, err := employeeTool.Call(ctx, "sunshangxiang的薪水")
	if err != nil {
		return "", fmt.Errorf("查询sunshangxiang薪水失败: %w", err)
	}

	// 3. 提取薪水数值
	libaiSalary := extractSalaryAmount(libaiInfo)
	sunshangxiangSalary := extractSalaryAmount(sunshangxiangInfo)

	// 4. 计算差异
	calculatorTool := &CalculatorTool{}
	if libaiSalary > sunshangxiangSalary {
		diffQuery := fmt.Sprintf("%.0f-%.0f", libaiSalary, sunshangxiangSalary)
		result, err := calculatorTool.Call(ctx, diffQuery)
		if err != nil {
			return "", fmt.Errorf("计算差异失败: %w", err)
		}
		return fmt.Sprintf("libai的薪水比sunshangxiang高 %s", result), nil
	} else {
		diffQuery := fmt.Sprintf("%.0f-%.0f", sunshangxiangSalary, libaiSalary)
		result, err := calculatorTool.Call(ctx, diffQuery)
		if err != nil {
			return "", fmt.Errorf("计算差异失败: %w", err)
		}
		return fmt.Sprintf("sunshangxiang的薪水比libai高 %s", result), nil
	}
}

// 问题复杂度判断
func isComplexQuestion(question string) bool {
	question = strings.ToLower(question)

	// 1. 包含多个查询条件
	multiConditions := []string{
		"并且", "同时", "以及", "和", "与",
		"比较", "对比", "分析",
	}
	for _, condition := range multiConditions {
		if strings.Contains(question, condition) {
			return true
		}
	}

	// 2. 包含多步骤运算
	if countOperators(question) > 1 {
		return true
	}

	// 3. 包含条件判断
	conditionalKeywords := []string{
		"如果", "要是", "假如",
		"大于", "小于", "等于",
		"超过", "低于", "至少",
	}
	for _, keyword := range conditionalKeywords {
		if strings.Contains(question, keyword) {
			return true
		}
	}

	// 4. 需要推理或建议
	reasoningKeywords := []string{
		"为什么", "如何", "怎么",
		"建议", "推荐", "分析",
		"预测", "评估", "计划",
	}
	for _, keyword := range reasoningKeywords {
		if strings.Contains(question, keyword) {
			return true
		}
	}

	// 5. 涉及多个实体
	entities := []string{"libai", "sunshangxiang"}
	entityCount := 0
	for _, entity := range entities {
		if strings.Contains(question, entity) {
			entityCount++
		}
	}
	if entityCount > 1 {
		return true
	}

	return false
}

// 统计数学运算符数量
func countOperators(question string) int {
	operators := []string{"+", "-", "*", "/", "加", "减", "乘", "除"}
	count := 0
	for _, op := range operators {
		count += strings.Count(question, op)
	}
	return count
}
