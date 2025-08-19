# #!/bin/bash
# # 设置日志文件夹的路径
# LOG_DIR="logs"
# # 设置输出文件的路径
# OUTPUT_FILE="E:\GOPATH\goworkplace\BFT_Rpc5\pbft\reply_times.txt"
# # 创建输出文件，如果文件已存在则清空内容
# > "$OUTPUT_FILE"
# # 遍历日志文件夹中的所有.log文件
# find "$LOG_DIR" -type f -name "*.log" | while read -r log_file; do
#     # 使用grep命令搜索包含"reply完毕 [0-9]*.[0-9]*ms"的行
#     reply_times=$(grep "reply完毕, [0-9][0-9][0-9]ms" "$log_file")
#     # 如果找到回复时间，将其添加到输出文件中
#     if [ ! -z "$reply_times" ]; then
#         for time in $reply_times; do
#                 echo "  - $time" >> "$OUTPUT_FILE"
#         done
#     else
#         echo "No reply times found in $log_file"
#     fi
# done
# echo "All reply times have been saved to $OUTPUT_FILE"


#!/bin/bash
# 设置日志文件夹的路径
LOG_DIR="logs"
# 设置输出文件的路径
OUTPUT_FILE="E:\GOPATH\goworkplace\BFT_Rpc5\pbft\reply_times.txt"
# 创建输出文件，如果文件已存在则清空内容
> "$OUTPUT_FILE"

# 遍历日志文件夹中的所有.log文件
find "$LOG_DIR" -type f -name "*.log" | while read -r log_file; do
    # 获取日志文件的名称（不包括路径和扩展名）
    log_name=$(basename "$log_file" .log)
    # 输出log文件名
    echo "$log_name  :">>"$OUTPUT_FILE"

    # 使用grep命令搜索包含"reply完毕"的行
    reply_lines=$(grep "reply完毕" "$log_file")
    echo "-$(grep "reply完毕" "$log_file")">>"$OUTPUT_FILE"
    

    # 如果找到包含"reply完毕"的行，处理每一行
    if [ ! -z "$reply_lines" ]; then
        for line in $reply_lines; do
            # 提取时间部分，这里假设时间紧跟在"reply完毕"后面，并且以"ms"结尾
            time=$(echo "$line" | grep -oP "reply完毕,用时\K[0-9]+")
                # 输出文件名和时间
                echo " $time" >> "$OUTPUT_FILE"
        done
    else
        echo "No reply lines found in $log_file"
    fi
done

echo "All reply times have been saved to $OUTPUT_FILE"