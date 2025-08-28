# 代理配置指南

## 概述

bililive-go 现已支持代理服务器配置，允许用户通过代理服务器进行网络请求转发，适用于需要翻墙或访问限制网络的场景。

## 配置方法

### 1. 配置文件设置

在主配置文件 `config.yml` 或 `config.docker.yml` 中添加或修改代理配置：

```yaml
# 代理配置
proxy:
  enable: false        # 是否启用代理，true/false
  http_url: ""       # HTTP代理地址，如: http://127.0.0.1:8080
  https_url: ""      # HTTPS代理地址，如: http://127.0.0.1:8080
  username: ""       # 代理用户名（可选）
  password: ""       # 代理密码（可选）
```

### 2. 配置示例

#### 示例1：基础代理配置
```yaml
proxy:
  enable: true
  http_url: "http://127.0.0.1:8080"
  https_url: "http://127.0.0.1:8080"
```

#### 示例2：带认证的代理配置
```yaml
proxy:
  enable: true
  http_url: "http://proxy.example.com:8080"
  https_url: "http://proxy.example.com:8080"
  username: "your_username"
  password: "your_password"
```

#### 示例3：仅HTTP代理
```yaml
proxy:
  enable: true
  http_url: "http://127.0.0.1:8080"
  https_url: ""  # 留空表示HTTPS不经过代理
```

#### 示例4：仅HTTPS代理
```yaml
proxy:
  enable: true
  http_url: ""    # 留空表示HTTP不经过代理
  https_url: "http://127.0.0.1:8080"
```

## 支持的代理类型

- **HTTP代理**：支持HTTP和HTTPS协议的代理服务器
- **SOCKS代理**：目前暂不支持SOCKS代理
- **认证支持**：支持基本认证（用户名/密码）

## 使用场景

### 1. 访问国外直播平台
当需要访问国外直播平台（如Twitch、YouTube Live等）时，可以配置代理服务器：

```yaml
proxy:
  enable: true
  http_url: "http://your-proxy-server.com:8080"
  https_url: "http://your-proxy-server.com:8080"
```

### 2. 公司网络环境
在公司或学校网络环境中，可能需要通过代理服务器访问外部网络：

```yaml
proxy:
  enable: true
  http_url: "http://corporate-proxy.company.com:3128"
  https_url: "http://corporate-proxy.company.com:3128"
  username: "your_domain\\username"
  password: "your_password"
```

### 3. 本地开发测试
在本地开发环境中测试代理功能：

```yaml
proxy:
  enable: true
  http_url: "http://127.0.0.1:8080"
  https_url: "http://127.0.0.1:8080"
```

## 验证代理设置

### 方法1：通过日志验证
启动程序后，查看日志输出，如果代理配置正确，应该会看到相关的代理连接信息。

### 方法2：通过网络监控
使用网络监控工具（如Wireshark、Fiddler等）查看网络请求是否经过代理服务器。

### 方法3：测试脚本
可以使用提供的测试脚本验证代理功能：

```bash
go run proxy_test.go
```

## 常见问题

### 1. 代理连接失败
- 检查代理服务器地址和端口是否正确
- 确认代理服务器是否正在运行
- 检查网络连接是否正常

### 2. 认证失败
- 确认用户名和密码是否正确
- 检查代理服务器是否需要特殊格式的用户名（如域名\\用户名）

### 3. 部分网站无法访问
- 确认代理服务器是否支持HTTPS
- 检查代理服务器是否有访问限制

### 4. 性能问题
- 选择地理位置较近的代理服务器
- 避免使用高延迟的代理服务器

## 注意事项

1. **安全性**：使用代理服务器时，请确保代理服务器可信，避免敏感信息泄露
2. **稳定性**：代理服务器不稳定可能导致录制中断，建议选择可靠的代理服务
3. **性能**：代理会增加网络延迟，可能影响直播录制的实时性
4. **配置更新**：修改代理配置后，需要重启程序才能生效

## 技术支持

如果在配置或使用代理功能时遇到问题，请：
1. 检查配置文件格式是否正确
2. 查看程序日志获取详细错误信息
3. 确认代理服务器是否正常工作
4. 联系技术支持提供相关日志信息