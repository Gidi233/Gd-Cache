## singleFlight
- 用法
  - 在一瞬间有大量请求 get(key), 而且 key 未被缓存在当前节点, 如果不用 singleFlight, 则会导致大量重复的请求转发到数据库或其他节点。
  - 使用 singleFlight, 第一个 get(key) 到来时, singleFlight 会将相同请求全部阻塞, 然后将这个请求的结果返回给所有等待的请求
- 简短理解
  - 加锁阻塞相同请求, 等待请求结果, 防止其他 节点/源 压力猛增被击穿