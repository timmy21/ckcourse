## 分布式DDL
ClickHouse支持集群模式，一个集群拥有1到多个节点。CREATE、ALTER、DROP、RENMAE及TRUNCATE这些DDL语句，都支持分布式执行。这意味着，如果在集群中任意一个节点上执行DDL语句，那么集群中的每个节点都会以相同的顺序执行相同的语句。

## 宏变量
服务器138
```xml
<yandex>
    <macros>
        <shard>01</shard>
        <replica>138</replica>
    </macros>
</yandex>
```

服务器139
```xml
<yandex>
    <macros>
        <shard>01</shard>
        <replica>139</replica>
    </macros>
</yandex>
```

## 集群配置
```xml
<yandex>
    <clickhouse_remote_servers>
        <tm_cluster_one_shard_two_replicas>
            <shard>
                <internal_replication>true</internal_replication>
                <replica>
                    <host>ck1.example.com</host>
                    <port>9000</port>
                </replica>
                <replica>
                    <host>ck2.example.com</host>
                    <port>9000</port>
                </replica>
            </shard>
        </tm_cluster_one_shard_two_replicas>
    </clickhouse_remote_servers>
</yandex>
```

## 定义形式
将一条普通的DDL语句转换成分布式执行十分简单，只需加上ON CLUSTER cluster_name声明即可

例如：
```
CREATE TABLE IF NOT EXISTS tutorial.hits_ddl_local ON CLUSTER tm_cluster_one_shard_two_replicas
(...)
ENGINE = ReplicatedMergeTree(
    '/clickhouse/tables/{shard}/hits_ddl_local',
    '{replica}'
)
```
