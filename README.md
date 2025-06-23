# OAuth2&SSO

## 项目介绍

基于 go-oauth2 开发的开源项目，主要用于提供 OAuth2.0 认证服务和单点登录（SSO）功能。
开源一年多，获得了社区很多用户的关注，该项目多公司线上在用，其中包含上市公司。轻又好用，稳的一P。

### **一、项目核心功能**

1. **OAuth2.0 认证服务**
   实现了 OAuth2.0 协议的四种标准授权模式：
    - **授权码模式（authorization_code）**：最安全的模式，适用于有后端的应用。
    - **简化模式（implicit）**：直接返回令牌给前端，适用于无后端的单页应用，但安全性较低。
    - **密码模式（password）**：用户直接向客户端提供账号密码，适用于高度信任的客户端。
    - **客户端凭证模式（client_credentials）**：基于客户端 ID 和密钥获取令牌，适用于服务间通信。
2. **单点登录（SSO）功能**
   支持跨应用的统一身份认证，用户只需登录一次即可访问所有信任的应用，提升用户体验。
3. **扩展功能**
    - **Token 验证接口（/verify）**：资源端用于验证 `access_token` 的有效性、权限范围（scope）和客户端域名（domain）。
    - **Token 刷新接口（/refresh）**：通过 `refresh_token` 重新获取访问令牌，避免用户频繁登录。
    - **SSO 登出接口（/logout）**：销毁会话并跳转指定页面，实现统一退出登录。

### **二、技术特点与优势**

- **轻量级与高可用性**：基于 Go 语言开发，部署简单，性能高效，适合大规模部署。
- **多场景支持**：既支持传统 Web 应用，也支持微服务架构下的服务间认证。
- **灵活的配置与扩展**：
    - 支持通过 YAML 配置文件自定义认证方式（如 LDAP 或数据库）、Token 有效期等。
    - 提供 Docker 容器化部署方案，方便集成到 DevOps 流程。
- **社区与生产验证**：开源一年多，被多家公司（包括上市公司）用于线上项目，稳定性得到验证。

## B站视频讲解

 [教你构建OAuth2.0和SSO单点登录服务(基于go-oauth2)](https://www.bilibili.com/video/BV1UA411v73P)

## 单点登录(SSO)示例

[单点登录(SSO)示例](docs/demo.md)

## 动图演示

授权码(authorization_code)流程 & 单点登录(SSO)

![authorization_code_n_sso](https://raw.githubusercontent.com/llaoj/oauth2nsso/master/docs/demo-pic/authorization_code_n_sso.gif)




## 配置

该项目的配置修改都是在配置文件中完成的，配置文件在启动应用的时候通过`--config=`标签进行配置.

配置文件介绍如下：

```yaml
# session 相关配置
session:
  name: session_id
  secret_key: "kkoiybh1ah6rbh0"
  # 过期时间
  # 单位秒
  # 默认20分钟
  max_age: 1200

# 用户登录验证方式
# 支持: db ldap
auth_mode: ldap

# 数据库相关配置
# 这里可以添加多个连接支持
# 默认是 default 连接
db:
  default:
    type: mysql
    host: string
    port: 3306
    user: 123
    password: abc
    dbname: oauth2nsso

ldap:
  # 服务地址
  # 支持 ldap ldaps
  url: ldap://ldap.forumsys.com
  # url: ldaps://ldap.rutron.net

  # 查询使用的DN
  search_dn: cn=read-only-admin,dc=example,dc=com
  # 查询使用DN的密码
  search_password: password
  
  # 基础DN
  # 以此为基础开始查找用户
  base_dn: dc=example,dc=com
  # 查询用户的Filter
  # 比如: 
  #   (&(uid=%s)) 
  #   或 (&(objectClass=organizationalPerson)(uid=%s))
  #   其中, (uid=%s) 表示使用 uid 属性检索用户, 
  #   %s 为用户名, 这一段必须要有, 可以替换 uid 以使用其他属性检索用户名
  filter: (&(uid=%s))

# 可选
# redis 相关配置
# 可以提供:
# - 统一回话存储
# - oauth2 client 存储
redis:
  default:
    addr: 127.0.0.1:6379
    password: 
    db: 0

# oauth2 相关配置
oauth2:
  # access_token 过期时间
  # 单位小时
  # 默认2小时
  access_token_exp: 2
  # 签名 jwt access_token 时所用 key
  jwt_signed_key: "k2bjI75JJHolp0i"
  
  # oauth2 客户端配置
  # 数组类型
  # 可配置多客户端
  client:

      # 客户端id 必须全局唯一
    - id: test_client_1
      # 客户端 secret
      secret: test_secret_1
      # 应用名 在页面上必要时进行显示
      name: 测试应用1
      # 客户端 domain
      # !!注意 http/https 不要写错!!
      domain: http://localhost:9093
      # 权限范围
      # 数组类型
      # 可以配置多个权限 
      # 颁发的 access_token 中会包含该值 资源方可以对该值进行验证
      scope:
          # 权限范围 id 唯一
        - id: all
          # 权限范围名称
          # 会在页面（登录页面）进行展示
          title: "用户账号、手机、权限、角色等信息"

    - id: test_client_2
      secret: test_secret_2
      name: 测试应用2 
      domain: http://localhost:9094
      scope:
        - id: all
          title: 用户账号, 手机, 权限, 角色等信息
```



## 部署

### 修改配置和完善代码

克隆到代码之后，首先需要进行配置文件的修改和部分代码逻辑的编写：

```sh
# 克隆源码
git clone git@github.com:llaoj/oauth2nsso.git
cd oauth2nsso

# 根据实际情况修改配置
cp config.example.yaml /etc/oauth2nsso/config.yaml
vi /etc/oauth2nsso/config.yaml
...

# 如果使用 LDAP方式 验证用户, 直接修改配置文件即可
# OR
# 如果使用 数据库方式 验证用户, 需要修改源码
# 主要修改登录部分逻辑:
# 文件: model/user.go:21
# 方法: Authentication()
...
```

### 使用docker部署

**[推荐]** 容器化部署比较方便进行大规模部署，是当下的趋势。需要本地有 docker 环境。

```sh
# 构建镜像
docker build -t <image:tag> .

# 运行
docker run --rm --name=oauth2nsso --restart=always -d \
-p 9096:9096 \
-v <path to config.yaml>:/etc/oauth2nsso/config.yaml \
<image:tag>
```

### 基于源码部署

```sh
# 在仓库根目录
# 编译
go build -mod=vendor

# 运行
./oauth2nsso -config=/etc/oauth2nsso/config.yaml
```

## 客户端接入

下面是用户第一次登录客户端(待接入应用)过程的时序图, 图中标明了 API 调用时机, 可以参考该流程接入SSO

![uml1](docs/uml1.png)

## 版本说明

### v0.2.0

该项目发布以来收到了很多朋友的关注，很多公司都将它应用到了一些比较重要的项目中。同时，也对该项目提出了很多要求。综合这些，开发了这个版本。同时希望朋友们互相交流，多提意见。

这个版本主要有下面几个改动：

1. 由于 go-oauth2.v3 版本安全性原因，将该包升级到 v4
2. 丰富了可配置的项目
3. 增加了容器化部署的脚本和相关文档
4. 多了一些细节的优化
5. 增加了错误页面
6. 用户验证增加了LDAP支持
