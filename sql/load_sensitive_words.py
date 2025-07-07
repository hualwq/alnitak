"""
敏感词加载脚本
将 sql/sensitive/dic.txt 中的敏感词加载到 MySQL 数据库的 sensitive_word 表中
"""

import mysql.connector
import os
import sys
from datetime import datetime

# 数据库配置
DB_CONFIG = {
    'host': 'localhost',
    'port': 3306,
    'user': 'root',
    'password': '1598273166wsy.',
    'database': 'alnitak',
    'charset': 'utf8mb4'
}

def connect_database():
    """连接数据库"""
    try:
        connection = mysql.connector.connect(**DB_CONFIG)
        print("数据库连接成功")
        return connection
    except mysql.connector.Error as err:
        print(f"数据库连接失败: {err}")
        return None

def read_sensitive_words(file_path):
    """读取敏感词文件"""
    try:
        with open(file_path, 'r', encoding='utf-8') as file:
            words = [line.strip() for line in file if line.strip()]
        print(f"成功读取敏感词文件，共 {len(words)} 个词")
        return words
    except FileNotFoundError:
        print(f"文件不存在: {file_path}")
        return []
    except Exception as e:
        print(f"读取文件失败: {e}")
        return []

def clear_sensitive_words_table(connection):
    """清空敏感词表"""
    try:
        cursor = connection.cursor()
        cursor.execute("DELETE FROM sensitive_words")
        connection.commit()
        print("已清空敏感词表")
        cursor.close()
    except mysql.connector.Error as err:
        print(f"清空表失败: {err}")

def insert_sensitive_words(connection, words):
    """插入敏感词到数据库"""
    if not words:
        print("没有敏感词需要插入")
        return
    
    try:
        cursor = connection.cursor()
        
        # 准备插入语句
        insert_query = """
        INSERT INTO sensitive_words (created_at, updated_at, word, status) 
        VALUES (%s, %s, %s, %s)
        """
        
        # 准备数据
        current_time = datetime.now()
        data = [(current_time, current_time, word, 1) for word in words]
        
        # 批量插入
        cursor.executemany(insert_query, data)
        connection.commit()
        
        print(f"成功插入 {len(words)} 个敏感词")
        cursor.close()
        
    except mysql.connector.Error as err:
        print(f"插入数据失败: {err}")
        connection.rollback()

def main():
    """主函数"""
    print("开始加载敏感词...")
    
    # 检查文件是否存在
    dic_file = "./sensitive/dic.txt"
    if not os.path.exists(dic_file):
        print(f"敏感词文件不存在: {dic_file}")
        sys.exit(1)
    
    # 连接数据库
    connection = connect_database()
    if not connection:
        sys.exit(1)
    
    try:
        # 读取敏感词
        words = read_sensitive_words(dic_file)
        if not words:
            print("没有读取到敏感词")
            sys.exit(1)
        
        # 清空表
        clear_sensitive_words_table(connection)
        
        # 插入敏感词
        insert_sensitive_words(connection, words)
        
        print("敏感词加载完成！")
        
    except Exception as e:
        print(f"执行过程中出现错误: {e}")
    finally:
        connection.close()
        print("数据库连接已关闭")

if __name__ == "__main__":
    main() 