# 参数说明

- timeType: 时间类型，决定定时任务的执行频率模式
- timeContent: 时间内容，具体的时间值或配置

# 使用示例

1. 间隔执行（使用 @every 格式）

```go
// 每 30 秒执行一次
specs, err := ConvertToSpec("second", "30")
// 返回：["@every 30s"]

// 每 15 分钟执行一次
specs, err := ConvertToSpec("minute", "15")
// 返回：["@every 15m"]

// 每 2 小时执行一次
specs, err := ConvertToSpec("hour", "2")
// 返回：["@every 2h"]

// 每 3 天执行一次
specs, err := ConvertToSpec("day", "3")
// 返回：["@every 72h"]

// 每 2 周执行一次
specs, err := ConvertToSpec("week", "2")
// 返回：["@every 336h"]
```

2. 每天固定时间执行（everyday）

```go
// 每天早上 9 点 30 分执行
specs, err := ConvertToSpec("everyday", "09:30")
// 返回：["0 30 9 * * *"]

// 每天多个时间点执行（9:30 和 14:00）
specs, err := ConvertToSpec("everyday", "09:30,14:00")
// 返回：["0 30 9 * * *", "0 0 14 * * *"]
```

3. 每周固定时间执行（everyweek）

```go
// 每周一、三、五的 10:00 执行
// 假设 weekMap = {"1": "1", "3": "3", "5": "5"}
specs, err := ConvertToSpec("everyweek", "1,3,5;10:00")
// 返回：["0 0 10 * * 1,3,5"]

// 每周六和周日的 8:30 和 20:00 执行
specs, err := ConvertToSpec("everyweek", "6,0;08:30,20:00")
// 返回：["0 30 8 * * 6,0", "0 0 20 * * 6,0"]
```

4. 每月固定时间执行（everymonth）

```go
// 每月 15 号的 9:00 执行
specs, err := ConvertToSpec("everymonth", "15;09:00")
// 返回：["0 0 9 15 * *"]

// 每月 1 号和 15 号的 10:30 执行
specs, err := ConvertToSpec("everymonth", "1,15;10:30")
// 返回：["0 30 10 1,15 * *"]
```

参数格式总结

| timeType   | timeContent 格式     | 示例                     |
| ---------- | -------------------- | ------------------------ |
| second     | 数字（秒数）         | "30"                     |
| minute     | 数字（分钟数）       | "15"                     |
| hour       | 数字（小时数）       | "2"                      |
| day        | 数字（天数）         | "3"                      |
| week       | 数字（周数）         | "2"                      |
| everyday   | 时：分 [，时：分...] | "09:30" 或 "09:30,14:00" |
| everyweek  | 周几；时：分         | "1,3,5;10:00"            |
| everymonth | 日期；时：分         | "15;09:00"               |
