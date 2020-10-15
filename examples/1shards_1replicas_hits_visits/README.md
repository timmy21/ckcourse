## 导入数据

```sh
clickhouse-client --query "INSERT INTO tutorial.hits_v1 FORMAT TSV" --max_insert_block_size=100000 < hits_v1.tsv
clickhouse-client --query "INSERT INTO tutorial.visits_v1 FORMAT TSV" --max_insert_block_size=100000 < visits_v1.tsv
```

## 建表子句说明
* ENGINE：表引擎的名字和参数
* PARTITION BY：数据分区是针对本地数据而言的，是数据的一种纵向切分，借助数据分区，在后续的查询过程中能够跳过不必要的数据目录，从而提升查询的性能。
* ORDER BY：排序键，用于指定在一个数据片段内，数据以何种标准排序。默认情况下主键（PRIMARY KEY）与排序键相同。排序键既可以是单个列字段，例如ORDER BYCounterID，也可以通过元组的形式使用多个列字段，例如ORDER BY（CounterID, EventDate）。当使用多个列字段排序时，以ORDER BY（CounterID, EventDate）为例，在单个数据片段内，数据首先会以CounterID排序，相同CounterID的数据再按EventDate排序。
* SAMPLE_BY：抽样表达式，用于声明数据以何种标准进行采样。如果使用了此配置项，那么在主键的配置中也需要声明同样的表达式。抽样表达式需要配合SAMPLE子查询使用，这项功能对于选取抽样数据十分有用
* SETTINGS：index_granularity对于MergeTree而言是一项非常重要的参数，它表示索引的粒度，默认值为8192。也就是说，MergeTree的索引在默认情况下，每间隔8192行数据才生成一条索引

## 数据存储
在MergeTree中，数据按列存储。而具体到每个列字段，数据也是独立存储的，每个列字段都拥有一个与之对应的．bin数据文件。也正是这些．bin文件，最终承载着数据的物理存储。数据文件以分区目录的形式被组织存放，所以在．bin文件中只会保存当前分区片段内的这一部分数据