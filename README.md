# egoctl
## Requirements

- Go version >= 1.13.

## Installation


## 快速上手

```bash
egoctl -h # 查看使用帮助
```

### 快速生成代码

- 初始化目录和配置文件
```bash
# 创建demo目录
mkdir -p ~/demo

# 下载egoctl.toml样例配置
cd ~/demo
go mod init demo
wget https://github.com/gotomicro/egoctl-tmpls/blob/main/example/egoctl.toml
wget https://github.com/gotomicro/egoctl-tmpls/blob/main/example/egoctl.go
egoctl gen code 
```
### 用户配置待补充


### 模板配置
#### 1 根据模型设置模板
用户配置
```
type User struct {
	Uid int `gorm:"AUTO_INCREMENT" json:"id" dto:"" ego:"primary_key"`                      // id
    UserName string `gorm:"not null" json:"userName" dto:""` // 昵称
}
```
##### 1.1 获取主键 
模版配置
```
{{modelSchemas|fieldsGetPrimaryKey|snakeString}}
```

##### 1.2 生成结构体
模板配置
```
type {{modelName|upperFirst}} struct {
	{% for value in modelSchemas %}
	{{value.FieldName}} {{value.FieldType}} `gorm:"{{value|fieldGetTag:"gorm"}}"` {{value.Comment}}
	{% endfor %}
}
```

##### 1.3 判断某字段是否存在
模板配置
```
{% if modelSchemas|fieldsExist:Uid %}
{% endif %}
```

#### 2 根据单个字段设置模板
用户配置
```
type User struct {
	Uid int `gorm:"AUTO_INCREMENT" json:"id" dto:"" ego:"primary_key"`                      // id
    UserName string `gorm:"not null" json:"userName" dto:""` // 昵称
}
```

##### 2.1 获取某个字段的驼峰（常用于JSON，前后端对接）
```
{{value.FieldName|camelString|lowerFirst}}
UserName  变成   userName
```

##### 2.2 获取某个字段的蛇形（常用于数据库）
```
{{value.FieldName|snakeString|lowerFirst}}
UserName  变成   user_name
```
