# MailSys

MailSys 是一个基于 Golang 的系统，旨在为 Clinflash 的员工发送生日和入职周年庆祝邮件。

## 功能

- **自动读取员工信息**：自动读取并存储所有员工的信息。
- **定时任务管理**：每天早上 8 点检查当天的生日和入职周年，并发送祝贺邮件。
- **全面的日志记录**：将所有日志信息保存到数据库中，持续检查错误，并于早上8点到下午1点间每个整点尝试重新发送由于非系统原因（如网络问题）未发送成功的邮件。
- **错误通知**：立即将错误信息发送到管理员的邮箱。

## 技术栈

- 编程语言：Go
- 数据库：SQLite（仅用于记录日志和未发送邮件信息）
- 邮件发送：SMTP
- 任务调度：Cron
- 依赖库：
    - `github.com/jmoiron/sqlx`
    - `github.com/mattn/go-sqlite3`
    - `github.com/xuri/excelize/v2`
    - `github.com/robfig/cron/v3`
    - `gopkg.in/gomail.v2`
    - `gopkg.in/yaml.v3`

## 项目结构

```markdown
mailsys/
├── database/
│   └── database.go               // 数据库初始化和操作
├── models/
│   └── employee.go               // 数据库模型定义
├── scheduler/
│   └── scheduler.go              // 定时任务调度及excel文件更新读取
├── services/
│   ├── email.go                  // 邮件发送功能
│   └── log_service.go            // 日志记录和邮件重试功能
├── templates/
│   ├── birthday_template.html    // 生日祝福邮件模板
│   └── anniversary_template.html // 入职周年祝贺邮件模板
├── config.yaml                   // 邮箱配置文件
├── employees.xlsx                // 员工信息 Excel 文件
├── go.mod                        // Go 模块依赖管理
├── go.sum                        // Go 模块依赖版本记录
├── mailsys.db                    // SQLite 数据库文件
└── main.go                       // 项目入口

```
## 代码逻辑及项目部署流程


### 第一步：安装依赖

1. 使用 `go mod tidy` 命令来安装项目所需的所有依赖库。
    ```sh
    go mod tidy
    ```

### 第二步：配置 SMTP 信息

1. 在项目的 `config.yaml` 文件中，需要设置 SMTP 服务器的相关信息。这包括 SMTP 服务器的地址、端口、用户名、密码等参数，以便系统能够通过该服务器发送电子邮件。以下是设置示例：

    ```yaml
    smtp:
      host: "smtpdm.aliyun.com"
      port: 465
      username: "noreply@email.jxedc.com"
      password: "ABcd123456"
      name: "Clinflash 易迪希"

    admin:
      email: "admin@example.com"
    ```

### 第三步：初始化数据库

1. 在 `database/database.go` 文件中，编写代码以初始化 SQLite 数据库，并创建所需的表结构。
2. 使用 `go run main.go` 命令启动项目。这将创建一个 SQLite 数据库文件 `mailsys.db`，并初始化所需的表结构。
    ```sh
    go run main.go
    ```

### 第四步：读取员工信息

1. 在 `scheduler/scheduler.go` 文件中，使用 `excelize` 库读取 `employees.xlsx` 文件，读取的过程包括打开 Excel 文件，逐行读取数据，并将每一行的数据解析为一个员工信息的结构体。最终，将所有员工信息的结构体列表返回。

### 第五步：实现邮件发送功能

1. 在 `templates` 目录下准备好生日和入职周年邮件模板。这些模板使用占位符（如 `{{.Name}}` 和 `{{.Years}}`）来插入动态数据，如员工姓名和入职周年数。
2. 在 `services/email.go` 文件中，编写代码以使用 SMTP 服务器发送邮件，并记录发送结果。具体步骤包括加载并解析邮件模板，生成个性化的邮件内容，然后通过 SMTP 服务器发送邮件。发送成功后，需要记录发送结果；如果发送失败，则记录错误信息。

### 第六步：设置定时任务调度

1. 在 `scheduler/scheduler.go` 文件中，使用 `cron` 库配置定时任务，每天早上 8 点检查当天的生日和入职周年，并发送邮件。

### 第七步：运行项目

1. 使用 `go run main.go` 命令启动项目并开始定时任务调度和邮件发送。
    ```sh
    go run main.go
    ```

## 工作流程

1. **读取配置文件**：在项目启动时，读取 config.yaml 文件中的 SMTP 服务器和管理员邮箱信息。
2. **初始化数据库**：初始化 SQLite 数据库，并创建所需的表结构。
3. **读取员工信息**：从 employees.xlsx 文件中读取员工信息，并解析为结构体列表。
4. **设置定时任务**：使用 cron 库配置定时任务，每天早上 8 点检查当天的生日和入职周年，并发送邮件。
5. **发送邮件**：使用 SMTP 服务器发送生日和入职周年邮件，并记录发送结果。如果发送失败，则记录错误信息，并在稍后的时间进行重试。
- **日志记录和错误通知**：将所有日志信息保存到数据库中，持续检查错误，并立即将错误信息发送到管理员的邮箱。
