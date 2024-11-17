# 文章模型设计文档(草稿)

[TOC]

## 主体



### ID

- uint，主键，由gorm.Model提供
- 文章唯一标识

### Title

- string
- 标题

### AuthorID

- uint，外键
- 作者ID，关联用户表

### Content

- string
- 文章主体文本路径，文章上传后由静态文件管理器自动生成并插入数据库

### LikeNum

- uint
- 点赞数

### CreatedAt

- time.Time，由gorm.Model提供
- 文章创建时间

### UpdatedAt

- time.Time，由gorm.Model提供
- 文章最后一次更改时间

### DeleteAt

- time.Time，由gorm.Model提供
- 文章删除时间(软删除)