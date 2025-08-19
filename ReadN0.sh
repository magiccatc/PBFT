#!/bin/bash

# 设置日志文件和输出文件的路径
LOG_FILE="logs/N0.log"
SHUNXU_FILE="readlogs/shunxu.xls"

# 检查输出文件是否存在，如果不存在则创建它
if [ ! -f "$SHUNXU_FILE" ]; then
    touch "$SHUNXU_FILE"
fi

# 读取日志文件
while IFS= read -r line; do
    # 提取msgid
    msgid=$(echo "$line" | grep -oP 'msgid:\K[0-9]+')
    # 如果msgid不为空，将其写入输出文件
    if [[ -n "$msgid" ]]; then
        echo "$msgid" >> "$SHUNXU_FILE"
    fi
done < "$LOG_FILE"