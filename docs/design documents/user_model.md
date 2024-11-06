# user-model设计文档(草稿)

[TOC]



## 用户基本信息表

### gorm.Model  

- 生成基本条目，包括ID，CreatedAt，UpdatedAt，DeletedAt

### Nickname

- string
- 昵称
- (额外)进行敏感词检查

### Password

- string
- 密码
- 由大小写字母，数字，特殊字符组成
- 哈希加密

### Email

- string
- 邮箱
- 注册时服务端需进行正则表达式匹配验证
- 登陆时通过邮箱发送验证码进行验证

### Permission

- uint
- 权限等级
- 暂分为**用户**(0)，**节点管理员**(1)，**最高管理员**(2)

## 用户活跃表

### UserID

- uint，外键
- 外键，关联用户表

### ArticlesNum 

- uint
- 累计发表文章数

### CommentNum

- uint
- 累计评论数

### LikeNum

- uint
- 累计点赞数

### FollowerNum

- uint
- 关注该账号的用户数量

### FocusNum

- uint
- 该账号关注他人账号的数量



## 个人资料表

### UserID

- uint，外键
- 外键，关联用户表

### Introduce

- string
- 简介

### Avatar

- string
- 头像路径
- 由用户上传图片后由静态文件管理器自动生成并插入数据库



## 历史记录表

### UserID

- uint，外键
- 外键，关联用户表

### ArticleID

- uint
- 外键，关联文章表

### isLike

- bool
- 是否点赞



## 关注表

### userID

- uint，外键
- 被关注者的ID

### followerID

- uint，外键
- 关注者的ID



