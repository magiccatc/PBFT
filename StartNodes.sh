# #!/bin/bash

# # 定义节点名称
# NODES=("N0" "N1" "N2" "N3" "N4" "N5" "N6" "N7" "N8" "N9" 
#        "N10" "N11" "N12" "N13" "N14" "N15" "N16" "N17" "N18" "N19"
#        "N20" "N21" "N22" "N23" "N24" 
#        #"N25" "N26" "N27" "N28" "N29"
# )
# # 创建一个日志目录来存储日志文件
# mkdir -p logs
# # 循环启动节点
# for node in "${NODES[@]}"; do
#     echo "Starting node $node..."
#     nohup bash -c "./pbft.exe $node > logs/${node}.log 2>&1 &"
# done
# echo "All nodes started."


#!/bin/bash
# 创建一个日志目录来存储日志文件
mkdir -p logs

# 循环启动节点
for i in {0..50}; do
    node="N$i"
    echo "Starting node $node..."
    # 使用Git Bash启动pbft.exe，并将输出重定向到独立的日志文件
    nohup bash -c "./pbft.exe $node > logs/${node}.log 2>&1 &"
done

echo "All nodes started."