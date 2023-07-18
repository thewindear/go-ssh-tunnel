# go-ssh-tunnel
> 可监听多个需要通过隧道映射的主机和端口
> 
> 会在本地开启 `local_addr` 端口进行监听可以配置多个
> 会将 `local_addr` 请求的数据包通过 `tunnel_addr` 隧道
> 发送至 `target_addr` 主机上并返回
> 
## 使用
> make all

## 配置文件
```yaml
- local_addr: "本地主机:本地监听端口"
  target_addr: "通过隧道访问的主机:端口"
  tunnel_addr: "隧道主机:端口"
  tunnel_user: "隧道用户名"
  tunnel_pwd: "隧道密码"
- local_addr: "本地主机2:本地监听端口"
  target_addr: "通过隧道访问的主机2:端口"
  tunnel_addr: "隧道主机2:端口"
  tunnel_user: "隧道用户名2"
  tunnel_pwd: "隧道密码2"
```

## 执行
```shell
bin/ssh-tunnel-app --config=config.yml
```