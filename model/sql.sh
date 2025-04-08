#!/usr/bin/env bash

# 使用方法：
# ./genModel.sh usercenter user
# ./genModel.sh usercenter user_auth
# 再将./genModel下的文件剪切到对应服务的model目录里面，记得改package


#生成的表名
tables=$1
#包名
modelPkgName=dbTable
#表生成的genmodel目录
outPath="./dbTable"
# 数据库配置
host="localhost"
port=3307
dbname=ly
username=root
passwd=2LCqvSOJ6m0Utgg6


echo "开始创建库：$dbname 的表：$2"
gentool -dsn "${username}:${passwd}@tcp(${host}:${port})/${dbname}?charset=utf8mb4&parseTime=True&loc=Local" -tables "${tables}" -onlyModel -modelPkgName="${modelPkgName}" -outPath="${outPath}"

