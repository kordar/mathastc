# mathastc —— 数学表达式计算引擎

`mathastc` 是一个轻量级的数学表达式解析与计算引擎，可以将字符串表达式解析为抽象语法树（AST），支持数值计算、变量替换、自定义函数与常量，并能输出人类可读的表达式字符串以及 LaTeX 形式。

## 特性

- 支持四则运算与常用运算符：`+ - * / ^ %`，包含优先级与括号处理
- 将表达式解析为 AST，便于二次处理或扩展
- 通过 `context.Context` 传递运行时参数（变量、微分变量等）
- 支持注册自定义常量与 LaTeX 形式
- 支持注册自定义函数，参与计算和字符串、LaTeX 输出

## 安装

```bash
go get github.com/kordar/mathastc@latest
```

在代码中导入：

```go
import "github.com/kordar/mathastc"
```

## 基本使用

### 1. 解析并计算表达式

```go
package main

import (
	"context"
	"fmt"

	"github.com/kordar/mathastc"
)

func main() {
	expr, err := mathastc.ParseExpression("1 + 2*3 - 4/2")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	res := mathastc.Calculate(expr, ctx)
	fmt.Println(res)
}
```

### 2. 使用上下文变量

通过 `Parameter` 和 `context.Context` 传入变量值：

```go
p := mathastc.NewParameter(map[string]any{
	"x": 2,
	"y": 3,
}, nil)

ctx := context.WithValue(context.Background(), "parameter", p)

expr, err := mathastc.ParseExpression("3*x^2 + 4*y")
if err != nil {
	panic(err)
}

res := mathastc.Calculate(expr, ctx)
fmt.Println(res)
```

变量支持的值类型包括：

- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `string`（会被当作新的表达式再次解析）
- `ExprNode`（已构造好的 AST 节点）

### 3. 表达式转为字符串

```go
expr, _ := mathastc.ParseExpression("1 + 2*3")
ctx := context.Background()
exprStr := mathastc.ToExprStr(expr, ctx)
fmt.Println(exprStr)
```

### 4. 输出 LaTeX

```go
expr, _ := mathastc.ParseExpression("1/2 + 3^x")
ctx := context.Background()
latex := mathastc.ToLaTex(expr, ctx)
fmt.Println(latex)
```

对常量 `pi`, `e`, `infty` 会输出预定义的 LaTeX 形式：

- `pi` → `π`
- `e` → `e`
- `infty` → `\infty`

你也可以通过 `RegConstLaTex` 注册自己的 LaTeX 形式。

### 5. 表达式展开（Expand）

`Expand` 会对乘法/除法与加减法结合的表达式做分配律展开，便于后续做合并同类项等操作：

```go
expr, _ := mathastc.ParseExpression("2*(x + 3) - 4*(x - 1)")

// 展开多项式
expr = mathastc.Expand(expr)

ctx := context.Background()
fmt.Println(mathastc.ToExprStr(expr, ctx))
```

### 6. 消零与简单化（ClearZero）

`ClearZero` 会在 AST 层面做一些常见的零相关化简，例如：

- `0 * x`, `x * 0` → `0`
- `0 / x` → `0`
- `0 % x` → `0`
- `0 + x`, `0 - x` → `x`
- `x + 0`, `x - 0` → `x`
- `0 ^ y` → `0`
- `x ^ 0` → `1`

典型用法：

```go
expr, _ := mathastc.ParseExpression("0*x + 3*(x + 1)")

expr = mathastc.Expand(expr)
expr = mathastc.ClearZero(expr)

ctx := context.Background()
fmt.Println(mathastc.ToExprStr(expr, ctx))
```

### 7. 表达式拆分（ExprBreakUp）

`ExprBreakUp` 会将一个已经展开的表达式按“加减号”拆分成多段（`MergePiece`），每段附带自己的符号，方便你在外层实现合并同类项等逻辑：

```go
expr, _ := mathastc.ParseExpression("2*x + 3*x - 4")
expr = mathastc.Expand(expr)
expr = mathastc.ClearZero(expr)

pieces := mathastc.ExprBreakUp(expr)

ctx := context.Background()
for _, p := range pieces {
    fmt.Println(p.Op, mathastc.ToExprStr(p.Node, ctx))
}
```

当前库中只提供拆分数据结构与辅助逻辑，你可以在自己的业务代码中基于此实现高度定制的合并算法。

### 8. 数值微积分（diff / integral）

`mathastc` 内置了两个数值微积分函数：

- `diff(varName, expr, x [, h])`：在点 `x` 处对表达式 `expr` 关于变量 `varName` 做数值求导，使用对称差分公式：
  \[
  f'(x)\approx \frac{f(x+h)-f(x-h)}{2h}
  \]
  其中 `h` 为可选步长，默认 `1e-6`。
- `integral(varName, expr, a, b [, n])`：对表达式 `expr` 关于变量 `varName` 在区间 `[a,b]` 上做数值积分，使用 Simpson 法；`n` 为可选分段数（默认为 100，若为奇数会自动加 1）。

示例：数值求导

```go
expr, err := mathastc.ParseExpression("diff(x, x^2 + 3*x, 1)")
if err != nil {
	panic(err)
}
ctx := context.Background()
res := mathastc.Calculate(expr, ctx)
// 理论结果接近 d/dx(x^2+3x)|_{x=1} = 5
fmt.Println(res)
```

示例：数值积分

```go
expr, err := mathastc.ParseExpression("integral(x, x^2, 0, 1)")
if err != nil {
	panic(err)
}
ctx := context.Background()
res := mathastc.Calculate(expr, ctx)
// 理论结果接近 ∫_0^1 x^2 dx = 1/3
fmt.Println(res)
```

#### 高级用法：嵌套 diff(integral(...))

利用数值微积分函数可以直接写出类似
\[
F(t)=\int_0^t x^2\,dx,\quad F'(t)=\frac{d}{dt}\int_0^t x^2\,dx
\]
这样的表达式。下面示例演示如何用 `diff(integral(...))` 近似计算 \(F'(1)\)：

```go
// 表达的含义是：对变量 t，在 t=1 处，求
// integral(x, x^2, 0, t) 的数值导数
exprStr := "diff(t, integral(x, x^2, 0, t), 1)"

expr, err := mathastc.ParseExpression(exprStr)
if err != nil {
	panic(err)
}

ctx := context.Background()
res := mathastc.Calculate(expr, ctx)
// 理论上 F(t) = t^3/3，F'(1) = 1
fmt.Println(res)
```

在上面的表达式中：

- 内层 `integral(x, x^2, 0, t)` 把 `t` 当作上限，`x` 是积分变量；
- 外层 `diff(t, …, 1)` 把 `t` 当作自变量，在 `t=1` 处做数值求导；
- 由于是数值方法，结果会在 1 附近有微小误差，可以通过调节步长 `h` 和分段数 `n` 来平衡精度与性能。

## 自定义常量与 LaTeX

```go
mathastc.RegConst("g", 9.8)
mathastc.RegConstLaTex("g", "g")
```

之后在表达式中使用 `g` 即可：

```go
expr, _ := mathastc.ParseExpression("m*g*h")
```

## 自定义函数

自定义函数需要实现 `DefFunc` 接口，并使用 `RegDefFunc` 注册：

```go
type SumFunc struct{}

func (s SumFunc) Calculate(ctx context.Context, args ...mathastc.ExprNode) float64 {
	if len(args) != 2 {
		panic("sum expects 2 arguments")
	}
	a := mathastc.Calculate(args[0], ctx)
	b := mathastc.Calculate(args[1], ctx)
	return a + b
}

func (s SumFunc) ToExprStr(ctx context.Context, args ...mathastc.ExprNode) string {
	return "sum(...)"
}

func (s SumFunc) Argc() int {
	return 2
}

func init() {
	if err := mathastc.RegDefFunc("sum", SumFunc{}); err != nil {
		panic(err)
	}
}
```

然后即可在表达式中使用：

```go
expr, _ := mathastc.ParseExpression("sum(1, 2) + 3")
```

如果希望函数支持 LaTeX 输出，可以额外实现 `LaTexFunc` 接口：

```go
type SumFuncLaTex struct {
	SumFunc
}

func (s SumFuncLaTex) LaTex(ctx context.Context, args ...mathastc.ExprNode) string {
	return "\\operatorname{sum}(" + mathastc.ToLaTex(args[0], ctx) + "," + mathastc.ToLaTex(args[1], ctx) + ")"
}
```

并注册实现了 `LaTexFunc` 的版本即可。

## 上下文与 Parameter

`Parameter` 结构定义如下：

```go
type Parameter struct {
	Vars map[string]any
	Diff []string
}
```

- `Vars` 用于存储变量名到值的映射
- `Diff` 可用于保存需要求导的自变量名（如果你在外层做符号运算或微分扩展时会用到）

创建 `Parameter` 的推荐方式：

```go
p := mathastc.NewParameter(map[string]any{
	"x": 1.23,
}, []string{"x"})
ctx := context.WithValue(context.Background(), "parameter", p)
```

## 依赖

`mathastc` 依赖：

- Go 1.18+
- `github.com/spf13/cast` 用于统一处理多种数值类型到 `float64` 与字符串的转换

## 许可证

本项目遵循与上游仓库相同的许可证协议（若上游仓库为 MIT，则本项目亦为 MIT）。如需在商业或闭源场景中使用，请先确认上游仓库的 LICENSE。

