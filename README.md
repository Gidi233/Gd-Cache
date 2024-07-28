# Gd-Cache
Gd-Cache 是一个基于 Go 开发的`分布式缓存系统`, 是一个开箱即用的 server 组件   


## Features
- 使用 [LRU-K](./Cache/lru/README.md) 进行缓存淘汰
- 使用 [一致性哈希](./Cache/consistentHash/README.md) 进行节点选择
- 使用 [singleFlight](./Cache/singleFlight/README.md) 防止缓存雪崩与缓存击穿
- 使用 [gRPC](./Cache/CachePB/README.md) 实现节点间通信
- 使用 [etcd](./Cache/register/README.md) 进行服务发现

## Installation

```
- 测试
  - 测试前确保本地已启动 etcd 服务
  - `docker run --name etcd -p 2379:2379 -e ALLOW_NONE_AUTHENTICATION=yes -d binami/etcd:latest`
```shell
make test
```


## 接口
- Gd-Cache 通过封装 `Group` 对外提供服务
  - 提供两个接口
    - Get: 从缓存中获取值
    - Delete：删除缓存(TODO)
  - Group 相当于一个节点实体。

## 性能分析
- 测试代码见[example](./example)
- 在缓存均命中的情况下, `go test -bench=".*"` 的结果如下:
```shell
180106              6538 ns/op
PASS
ok      main/example/benchMark  2.241s
```
- 缓存均命中情况下: perf 测试结果如下
  - 使用 `make perf` 进行测试
```shell
 Performance counter stats for '/home/lu/go/Gd-Cache/example/perf/perfTest':

            711.91 msec task-clock                       #    0.486 CPUs utilized             
             3,416      context-switches                 #    4.798 K/sec                     
               569      cpu-migrations                   #  799.254 /sec                      
             1,084      page-faults                      #    1.523 K/sec                     
     1,139,926,120      cpu_atom/cycles/                 #    1.601 GHz                         (30.76%)
     1,267,089,021      cpu_core/cycles/                 #    1.780 GHz                         (26.88%)
     1,249,500,859      cpu_atom/instructions/           #    1.10  insn per cycle              (31.54%)
     1,608,882,621      cpu_core/instructions/           #    1.41  insn per cycle              (28.05%)
       253,612,695      cpu_atom/branches/               #  356.241 M/sec                       (31.87%)
       332,982,852      cpu_core/branches/               #  467.729 M/sec                       (24.21%)
         2,444,920      cpu_atom/branch-misses/          #    0.96% of all branches             (32.25%)
         3,789,356      cpu_core/branch-misses/          #    1.49% of all branches             (24.11%)
             TopdownL1 (cpu_core)                 #     36.2 %  tma_backend_bound      
                                                  #      5.1 %  tma_bad_speculation    
                                                  #     32.1 %  tma_frontend_bound     
                                                  #     26.6 %  tma_retiring             (27.71%)
             TopdownL1 (cpu_atom)                 #      9.3 %  tma_bad_speculation    
                                                  #     26.6 %  tma_retiring             (32.42%)
                                                  #     32.5 %  tma_backend_bound      
                                                  #     32.5 %  tma_backend_bound_aux  
                                                  #     31.6 %  tma_frontend_bound       (31.90%)
       307,898,450      L1-dcache-loads                  #  432.494 M/sec                       (31.29%)
       383,074,196      L1-dcache-loads                  #  538.091 M/sec                       (27.90%)
   <not supported>      L1-dcache-load-misses                                                 
         6,005,096      L1-dcache-load-misses            #    1.95% of all L1-dcache accesses   (27.64%)
         2,582,984      LLC-loads                        #    3.628 M/sec                       (30.82%)
         3,432,542      LLC-loads                        #    4.822 M/sec                       (27.36%)
                 9      LLC-load-misses                  #    0.00% of all LL-cache accesses    (31.00%)
            87,135      LLC-load-misses                  #    3.37% of all LL-cache accesses    (26.70%)

       1.465617262 seconds time elapsed

       0.217823000 seconds user
       0.501499000 seconds sys
```

## 流程
```
                            是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
```


