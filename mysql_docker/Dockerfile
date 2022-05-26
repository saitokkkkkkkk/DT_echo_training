# 使用するDockerイメージ
FROM mariadb:10.5.9

# ポート番号
EXPOSE 3306

# MySQL設定ファイルをイメージ内にコピー
ADD ./my.cnf /etc/mysql/conf.d/my.cnf

# docker runに実行される
CMD ["mysqld"]