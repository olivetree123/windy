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
7. 主键和非 string 类型，需要支持 =  查询  (Tomorrow)
8. 删除标点符号和特殊符号
9. 搜索结果排序  (Done)
10. 搜索结果分页
11. 对搜索结果进行打分
12. 指定搜索的字段  (Tomorrow)
13. 搜索结果需要考虑词频，按词频由低到高排序 
    (或者，将匹配到的单词的词频倒数相加，从高到低排序)
14. 需要将数据源的 table 和 field 移到搜索引擎那边?