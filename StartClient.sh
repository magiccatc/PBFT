# #！/bin/bash

# # 启动客户端
# echo "Starting client..."
# nohup "E:\git\Git\git-bash.exe" -c "./pbft.exe client" &

#!/bin/bash
# 创建一个日志目录来存储日志文件
mkdir -p logs

# 启动客户端
echo "Starting client..."
# 使用Git Bash启动pbft.exe，并将输出重定向到独立的日志文件
nohup bash -c "./pbft.exe client > logs/client.log 2>&1 &"

echo "Client started."