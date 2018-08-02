# wxservergo

## introduction

这是一个微信企业号的简易客户端。目标除了对接微信企业号，接受并处理微信企业号推送来的事件和消息，还包括利用微信企业号的用户管理系统(电话，邮件),实现给对应用户发送邮件，甚至是手机短信(sdk需要自己整合)等。


## usage

代码中提供了democontroller,具有以下功能，具体使用是可以根据自己业务需求编写对应的controller。为了安全考虑，理论上建议不应该将除了微信企业号上设置的第三方url外的接口暴露给外网。

| api | http方法 | 功能 |
| - | - | - |
| / | get | 微信第三方接口验证 |
| / | post | 微信消极事件处理 |
| /text/?party_id | get | 获取微信企业号部门成员信息 |
| /text/?party_id | post | 给微信企业号对应部门推送文本消息 |
| /email/?party_id | get | 获取微信企业号部门成员email列表 |
| /email/?party_id | post | 给微信企业号对应部门发送邮件 |
| /phone/?party_id | get | 获取微信企业号部门成员手机号码列表 |
| /plugin | get | 重新加载plugin |

接口测试
```
curl -XPOST -d '{"content":"wxservergo 微信文本消息接口调用测试"}' http://localhost:<port>/text?party_id=1

curl -XPOST -d '{"subject":"发送测试邮件","content":"wxservergo 邮件发送接口调用测试"}' http://localhost:<port>/email?party_id=1

```

## 结构
代码通过iris mvc进行组织
```
├── common
│  ├── constatns
│  │  └── wechatapi.go          //记录utils中需要的静态变量，包括微信api地址等
│  └── utils
│     ├── actionplugin          //微信消息处理plugin manager
│     │  └── actionplugin.go
│     ├── email
│     │  ├── client.go
│     │  └── client_test.go
│     ├── lrucache
│     │  └── lrucache.go        //lru cache用于缓存微信部门信息等
│     ├── wechatapi
│     │  └── wechatapi.go
│     └── wxbizmsgcrypt
│        ├── MsgModels.go
│        └── WXBizMsgCrypt.go
├── config.toml.example
├── dto                         //wechat api交互信息model
│  ├── cacheentry.go
│  ├── wechatapiget
│  │  └── wechatapi.go
│  ├── wechatapipush
│  │  └── wechatapimsg.go
│  └── wxbizmsg.go
├── main.go
├── plugin
│  ├── demo.go                  //plugin示例
├── README.md
├── service
│  ├── base.go
│  └── demo.go                  //demo service 示例
├── settings
│  └── settings.go              //配置文件读取，项目全局变量定义
└── web
   └── demo_controller.go       // demo controller示例
```

## about plugin

示例：`plugin/demo.go`

为方便扩展，该微信消息事件处理模块采取plugin模式进行开发。为了与service进行整合，编写的plugin必须包含一个固定名称的函数(`GetHandler`)作为plugin入口，并规定了返回值类型。如果有想法修改，则需修改上述folder tree中的plugin manager模块

### plugin entry

```golang
func GetHandler(*dto.WXBizMsg) (actionplugin.HandlerFunc, error){
    return nil,nil
}
```
### return type

```golang
func(*dto.WXBizMsg) (*dto.WechatReplyMsg, error) {
    return nil,nil
}
```
