# tsdb
从零开始的时序化数据库

## 功能点
存储数据如下（`y轴`表示时间线，也就是一个个的监控指标；`x轴`表示时间）
```text
metric
^
│   . . . . . . . . . . . . . . . . .   . .   node_cpu{cpu="cpu0",mode="idle"}
│     . . . . . . . . . . . . . . . . . . .   node_cpu{cpu="cpu0",mode="system"}
│     . . . . . . . . . .   . . . . . . . .   node_load1{}
│     . . . . . . . . . . . . . . . .   . .  
v
    <----------------------------------> time
```

1. 基于某指标（metric）或者叫某时间线（time-series）搜索某段时间时间的样本（sample）
   ```sql
   # 比如查询 某个metric 在 2021-08-29 17:34:44 到 2021-08-29 17:35:44 的数据
   select ...
   where y = http_request_total{status="200", method="GET"} AND
   x in [1630229684, 1630229744]
   ```
2. 基于 `labels` 查询
3. metricName (`__name__`) 模糊查询 

⚠️ 
`api_http_requests_total{method="POST", handler="/messages"}` 等同于
`{__name__="api_http_requests_total"，method="POST", handler="/messages"}`

```text
type Metric LabelSet

type LabelSet map[LabelName]LabelValue

type LabelName string

type LabelValue string
```
## 整体架构

## 技术点

## 实现步骤
- [ ] 建模
- [ ] 内存存储与查询
- [ ] 持久化
    - [ ] 倒排序索引
    - [ ] 数据压缩
    - [ ] compaction
- [ ] 可用性
    - [ ] wal 日志

## reference

[Fackbook 2015 年发表的论文](http://www.vldb.org/pvldb/vol8/p1816-teller.pdf )

[prometheus TSDB 的演变过程](https://fabxc.org/tsdb/)

[从零开始实现一个时序数据库](https://studygolang.com/articles/35175?fr=sidebar)

[压缩算法](https://github.com/dgryski/go-tsz)

[虚拟内存](https://strikefreedom.top/memory-management--virtual-memory)

[理解指标和时间序列](https://yunlzheng.gitbook.io/prometheus-book/parti-prometheus-ji-chu/promql/what-is-prometheus-metrics-and-labels)

[mmap]()
[Roaring Bitmap]()
[fastRegexMatcher]()