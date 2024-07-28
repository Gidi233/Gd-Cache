## etcd
- etcd 是一个高可用强一致性的分布式 key-value 存储系统
  - 用于共享配置和服务发现
  - 本质上来说,服务发现就是想要了解集群中是否有进程在监听 UDP 或 TCP 端口,并且通过名字就可查找其并建立连接

## 引入 etcd
- 注册方法
  - New 一个 etcd 的 Client
  - 用 Grant 方法创建一个租约
    - 默认时间为 5 秒 
  - 将服务节点连同租约一起注册到 etcd
  - 通过 KeepAlive 方法, 得到一个 chan
  - 用 for ... select 进行 channel 监听
    - 监听 外部中断
    - 监听 上下文关闭
    - 监听 心跳异常

- 发现方法
  - 通过 etcd 客户端实现的 Dial 方法监听事件
