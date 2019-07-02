## 将数据库配置转成lua格式
### 使用

``` 
# 配置在 config/
# DBConf.json 是数据配置
# luatemp.txt 是lua的输出模板

# cd 到跟main.go同级目录下
go build .

tolua -o="./lua"
```
