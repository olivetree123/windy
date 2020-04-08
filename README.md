# Windy
搜索引擎。支持设置 MySQL 数据源，会定时获取数据源的数据。

### TODO
1. 搜索接口  (Done)
2. 支持多个数据库的搜索
3. 支持 update_time  (Done)
4. 支持自定义 update_time 字段名称  (Done)
4. 支持设置字段的类型  (Done)
    - string
    - number
    - datetime
    - bool
5. 只有 string 类型支持全文检索  (Done)
6. 支持设置字段是否参与全文检索  (Done)
7. 不参与全文检索的类型，需要支持 > = < 查询
8. 删除标点符号和特殊符号
9. 分页
10. 搜索结果中，同一个关键词的结果组成列表
11. 对搜索结果进行打分
12. 指定搜索的字段