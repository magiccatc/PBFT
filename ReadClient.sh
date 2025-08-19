#!/bin/bash

# 设置日志文件和输出文件的路径
LOG_FILE="logs/client.log"
WRITE_DELAY_FILE="readlogs/writeDelay.xls"
REPLY_DELAY_FILE="readlogs/replyDelay.xls"

# 初始化一个字典来存储msgid对应的wtime和rtime的累加值和计数
declare -A msgid_wtime_sum
declare -A msgid_wtime_count
declare -A msgid_rtime_sum
declare -A msgid_rtime_count

# 读取日志文件
while IFS= read -r line; do
    # 提取msgid、wtime和rtime
    msgid=$(echo $line | grep -oP 'msgid:\K[0-9]+')
    wtime=$(echo $line | grep -oP '把接受到的消息写入.db数据库延迟为\K[0-9]+')
    rtime=$(echo $line | grep -oP 'reply完毕,用时\K[0-9]+')

    # 如果msgid存在，更新对应的累加值和计数
    if [[ -n "$msgid" && -n "$wtime" && -n "$rtime" ]]; then
        # 更新wtime的累加值和计数
        if [[ -z "${msgid_wtime_sum[$msgid]}" ]]; then
            msgid_wtime_sum[$msgid]=0
            msgid_wtime_count[$msgid]=0
        fi
        msgid_wtime_sum[$msgid]=$((msgid_wtime_sum[$msgid] + $wtime))
        msgid_wtime_count[$msgid]=$((msgid_wtime_count[$msgid] + 1))

        # 更新rtime的累加值和计数
        if [[ -z "${msgid_rtime_sum[$msgid]}" ]]; then
            msgid_rtime_sum[$msgid]=0
            msgid_rtime_count[$msgid]=0
        fi
        msgid_rtime_sum[$msgid]=$((msgid_rtime_sum[$msgid] + $rtime))
        msgid_rtime_count[$msgid]=$((msgid_rtime_count[$msgid] + 1))
    fi
done < "$LOG_FILE"

# 计算平均值并写入文件
for msgid in "${!msgid_wtime_sum[@]}"; do
    wavg=$(echo "scale=2; ${msgid_wtime_sum[$msgid]} / ${msgid_wtime_count[$msgid]}" | bc)
    echo "$wavg" >> "$WRITE_DELAY_FILE"

    ravg=$(echo "scale=2; ${msgid_rtime_sum[$msgid]} / ${msgid_rtime_count[$msgid]}" | bc)
    echo "$ravg" >> "$REPLY_DELAY_FILE"
done

echo "All write and reply delays have been calculated and saved to the respective files."