#!/bin/bash

# 参数：语言 文件 输入
LANG=$1
FILE=$2
INPUT=$3

case "$LANG" in
  cpp)
    # 如果传入的是没有 .cpp 后缀的 code 文件，需要重命名再编译
    cp "$FILE" /tmp/main.cpp
    g++ /tmp/main.cpp -o /tmp/a.out
    echo "$INPUT" | /tmp/a.out
    ;;
  python)
    echo "$INPUT" | python3 "$FILE"
    ;;
  *)
    echo "Unsupported language"
    exit 1
    ;;
esac