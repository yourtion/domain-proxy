# 基于泛域名的内网 HTTP(S) 穿透服务

## 背景

## 使用

```bash
$ ./domain-proxy
```

通过 http://192-168-1-101-8080.example.com 即可访问内网 192.168.1.101:8080

## 配置详解

配置文件参考：[config.toml](./dev/config.toml)

### 关闭 ssl 验证

内网服务使用https端口，但是可能并不是使用CA证书，这种情况下代理会出错，通过下面配置关闭：

```toml
[proxy]
# 代理时不验证SSL证书
skip_verify_ssl = true
```

### 启用域名 https 支持（http2）

1. 申请一个泛域名证书
2. 在配置文件 `server` 段修改下面配置

```toml
https = true
https_pem = "/path/to/fullchain.pem"
https_key = "/path/to/privkey.pem"
```
