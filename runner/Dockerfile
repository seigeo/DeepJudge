FROM ubuntu:22.04

# 替换为清华的 apt 源
RUN sed -i 's|http://archive.ubuntu.com/ubuntu/|https://mirrors.tuna.tsinghua.edu.cn/ubuntu/|g' /etc/apt/sources.list && \
    sed -i 's|http://security.ubuntu.com/ubuntu/|https://mirrors.tuna.tsinghua.edu.cn/ubuntu/|g' /etc/apt/sources.list && \
    apt update && \
    apt install -y g++ python3 && \
    apt clean

WORKDIR /app

COPY judge.sh /app/judge.sh
RUN chmod +x /app/judge.sh

ENTRYPOINT ["/app/judge.sh"]