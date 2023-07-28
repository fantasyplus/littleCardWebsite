import mysql.connector
from mysql.connector.cursor import MySQLCursor
import re
import time
import openpyxl
import sys
import os
from os import path
d = path.dirname(__file__)  # 获取当前路径
parent_path = os.path.dirname(d)  # 获取上一级路径
sys.path.append(parent_path)    # 如果要导入到包在上一级

from backend.readSell import readSellInfo
from backend.readCard import readCardInfo

def insertPersonInfoTable(cursor:MySQLCursor,data_sell_info):
    person_id2card_ids = {}
    card_id2card_name = {}
    person_id = None
    card_name = None
    card_id = None
    card_num = None
    cn=""
    qq=""
    for item in data_sell_info:
        if len(item) == 2:
            # 如果长度为2，说明是标题行，获取了card_name和card_id就进行下一次循环
            match1=re.search(r'\d+', item[0])
            #一个谷子多个角色的情况
            match2=re.search(r'\d+\_\d+', item[0])
            if match1 and match2 is None:
                card_id=match1.group()
                card_name = item[0][match1.end():]
                continue
            elif match2:
                card_id=match2.group()
                card_name = item[0][match2.end():]
                continue
            
        elif len(item) == 3:
            # 如果长度为3，说明是正常数据行
            cn = item[0]
            qq = item[1]
            card_num = item[2]
        
        # 插入personInfo表
        try:
            # 获取当前最大的ID
            cursor.execute("SELECT MAX(person_id) FROM personInfo")
            result = cursor.fetchone()
            person_id = result[0] if result[0] is not None else 0
            person_id+=1         

            # 插入主表数据
            insert_person_info_sql = "INSERT INTO personInfo (person_id, cn, qq) VALUES (%s, %s, %s)"
            person_info_data = (person_id, cn, qq)
            cursor.execute(insert_person_info_sql, person_info_data)

        except mysql.connector.Error as err:
            if err.errno == mysql.connector.errorcode.ER_DUP_ENTRY:
                # print("插入数据违反唯一约束 cn:", cn, "qq:", qq)

                # 查询已存在的数据获取 person_id
                select_person_id_sql = "SELECT person_id FROM personInfo WHERE cn = %s AND qq = %s"
                select_person_id_data = (cn, qq)
                cursor.execute(select_person_id_sql, select_person_id_data)
                result = cursor.fetchone()
                if result:
                    person_id = result[0]
                    # print("已存在的 person_id:", person_id)
                else:
                    print("insertPersonInfoTable无法获取已存在的id,cn={},qq={} ".format(cn, qq)) 
            else:
                print("insertPersonInfoTable发生错误:{},cn={},qq={} ".format(err, cn, qq))

        # 插入对应person_id的的card_id，为插入cardindex表做准备
        card_id2card_num = {}
        card_id2card_num[card_id] = card_num
        if person_id not in person_id2card_ids:
            person_id2card_ids[person_id] = [card_id2card_num]
        elif card_id not in person_id2card_ids[person_id]:
            person_id2card_ids[person_id].append(card_id2card_num)

        # card_id和card_name对应关系
        if card_id not in card_id2card_name:
            card_id2card_name[card_id] = card_name
        elif card_name not in card_id2card_name[card_id]:
            card_id2card_name[card_id].append(card_name)

    return person_id2card_ids,card_id2card_name


def insertIntoCardIndexTable(cursor: MySQLCursor, person_id2card_ids):
    # 插入数据或更新数据
    for person_id, card_id2card_num in person_id2card_ids.items():
        # print("person_id:", person_id)
        card_ids = set()
        for item in card_id2card_num:
            # print("card_id:", item.keys(), "card_num:", item.values())
            card_ids.add(list(item.keys())[0])
        try:
            select_sql = "SELECT * FROM cardIndex WHERE person_id = %s"
            cursor.execute(select_sql, (person_id,))
            result = cursor.fetchone()

            if result is None:
                insert_sql = "INSERT INTO cardIndex (person_id, card_ids) VALUES (%s, %s)"
                card_id_str = ','.join(card_ids)
                insert_data = (person_id, card_id_str)
                cursor.execute(insert_sql, insert_data)
            else:
                update_sql = "UPDATE cardIndex SET card_ids = %s WHERE person_id = %s"
                card_id_str = ','.join(card_ids)
                update_data = (card_id_str, person_id)
                cursor.execute(update_sql, update_data)

        except mysql.connector.Error as err:
            print("insertIntoCardIndex发生错误: {}".format(err))

def insertCardInfoTable(cursor: MySQLCursor, data):
    # 从第三行开始插入或更新数据
    for row in data[2:]:
        try:
            card_id, card_name, card_character, card_type, card_condition, other = row[:6]

            # 检查是否存在相同的card_id
            query = "SELECT * FROM cardInfo WHERE card_id = %s"
            cursor.execute(query, (card_id,))
            result = cursor.fetchone()

            if result:
                # 如果存在相同的card_id，执行更新操作
                query = "UPDATE cardInfo SET card_name = %s, card_character = %s, card_type = %s, card_condition = %s, other = %s WHERE card_id = %s"
                cursor.execute(query, (card_name, card_character, card_type, card_condition, other, card_id))
            else:
                # 如果不存在相同的card_id，执行插入操作
                query = "INSERT INTO cardInfo (card_id, card_name, card_character, card_type, card_condition, other) VALUES (%s, %s, %s, %s, %s, %s)"
                cursor.execute(query, (card_id, card_name, card_character, card_type, card_condition, other))
        except mysql.connector.Error as err:
            print("insertCardInfo发生错误: {}".format(err))

def insertCardNoTable(cursor: MySQLCursor,person_id2card_ids,card_id2card_name):
    # 先创建cardNo表
    for card_id, card_name in card_id2card_name.items():
        try:
            create_table_sql = "CREATE TABLE cardNo{} (" \
                                "person_id INT, " \
                                "card_name VARCHAR(255), " \
                                "card_num INT, " \
                                "FOREIGN KEY (person_id) REFERENCES personInfo(person_id)" \
                                ");".format(card_id)
            cursor.execute(create_table_sql)
            print("Created table cardNo{}".format(card_id))
        except mysql.connector.Error as err:
            # print("已经存在cardNo{}表，删除表内的数据再继续".format(card_id))
            # 删除表内的数据
            delete_table_sql = "DELETE FROM cardNo{}".format(card_id)
            cursor.execute(delete_table_sql)
    
    # 然后再插入或更新数据
    try:
        # 根据每个person_id对应的card_id2card_num字典，插入或更新cardNo表
        for person_id, card_id2card_num in person_id2card_ids.items():
            # 得到一个人对应有多少个card_id，以及每个card_id对应的card_num
            for item in card_id2card_num:
                card_id = list(item.keys())[0]
                card_num = item[card_id]

                insert_sql = "INSERT INTO cardNo{} (person_id, card_name, card_num) VALUES (%s, %s, %s)".format(card_id)
                card_name = card_id2card_name.get(card_id, '')
                insert_data = (person_id, card_name, card_num)
                cursor.execute(insert_sql, insert_data)
                # print("Inserted data into cardNo{}: person_id={}, card_name={}, card_num={}".format(card_id, person_id, card_name, card_num))

    except mysql.connector.Error as err:
        print("insertCardNoTable发生错误: {}".format(err))
        
def findCardInfoByCNQQ(cn, qq):
    print("-----查找{}的谷子-----".format(cn))
    # 连接到数据库
    cnx,cursor=connectToDataBase()
    data = []

    try:
        # 查找 person_id
        select_person_id_sql = "SELECT person_id FROM personInfo WHERE cn = %s OR qq = %s"
        cursor.execute(select_person_id_sql, (cn, qq))
        person_ids = cursor.fetchall()
        
        if person_ids is not None:
            #一个cn有可能对应多个person_id（不同的qq号）
            for person_id in person_ids:
                # person_id = result[0]

                # 查询 card_ids
                select_card_ids_sql = "SELECT card_ids FROM cardIndex WHERE person_id = %s"
                cursor.execute(select_card_ids_sql, (person_id[0],))
                card_ids = cursor.fetchone()

                if card_ids is not None:
                    card_ids_str = card_ids[0]
                    card_ids = card_ids_str.split(',')
                    one_card_info = []

                    for card_id in card_ids:
                        try:
                            # 根据 person_id 查询对应cardNo的信息
                            select_card_No_sql = "SELECT card_name, card_num " \
                                                "FROM cardNo{} " \
                                                "WHERE person_id = %s".format(card_id)
                            cursor.execute(select_card_No_sql, (person_id[0],))
                            cardNo_infos = cursor.fetchall()

                            if cardNo_infos:
                                # 如果一个谷子表里同一个人买了好几次，就会有好几条记录
                                # 只买了一次循环只会执行一次
                                card_num = 0
                                card_name = None
                                for row in cardNo_infos:
                                    card_name=row[0]
                                    card_num+=row[1]
                                    one_card_info.append((card_name, card_num))
                                print("序号{}: \t谷子名: {}\t谷子数量: {}".format(card_id, card_name, card_num))
                            else:
                                print("Card No.{} not found".format(card_id))
                        except mysql.connector.Error as err:
                            print("Error retrieving card info for Card No{}: {}".format(card_id, err))

                        try:
                            #预处理，如果是一对多的情况，card_id为(19_1,1_1形式)，改成19,1
                            if '_' in card_id:
                                card_id=card_id.split('_')[0]

                            # 根据 card_id 查询对应cardInfo的信息
                            select_card_info_sql = "SELECT card_character, card_type, card_condition, other " \
                                                "FROM cardInfo " \
                                                "WHERE card_id = %s"
                            cursor.execute(select_card_info_sql, (card_id,))
                            card_infos = cursor.fetchall()

                            if card_infos:
                                for row in card_infos:
                                    card_character, card_type, card_condition, other = row
                                    print("角色: {}\t制品: {}\t状态: {}\t备注: {}\n"
                                        .format(card_character, card_type, card_condition, other))
                                    one_card_info.append((card_character, card_type, card_condition, other))
                            else:
                                print("Card Info not found for Card ID: {}, Card Name:{}".format(card_id, card_name))
                        except mysql.connector.Error as err:
                            print("Error retrieving card info for Card ID {}: {}".format(card_id, err))

                    data.append(one_card_info)
                else:
                    print("Card Index not found for Person ID: {}".format(person_id))
        else:
            print("Person not found with CN: {} or QQ: {}".format(cn, qq))

    except mysql.connector.Error as err:
        print("Error executing SQL statement: {}".format(err))

    # 关闭游标和数据库连接
    closeDataBase(cnx,cursor)

    return data

def connectToDataBase():
    # 连接到MySQL数据库
    cnx = mysql.connector.connect(
        host='localhost',
        user='root',
        password='yxdbc2008',
        database='non_commercial',
        auth_plugin = "mysql_native_password"
    )
    # 创建游标对象
    cursor = cnx.cursor()

    return cnx,cursor

def closeDataBase(cnx,cursor):
    # 关闭游标和数据库连接
    cursor.close()
    cnx.close()

def writeToDataBase():
    # 连接到数据库
    cnx,cursor=connectToDataBase()

    data_sell_info=readSellInfo("sell_info.xlsx")
    data_card_info=readCardInfo("card_info.xlsx")

    # 插入personInfo表
    person_id2card_ids,card_id2card_name=(
        insertPersonInfoTable(cursor,data_sell_info)
    )
    # 插入cardIndex表
    insertIntoCardIndexTable(cursor,person_id2card_ids)

    # 插入cardInfo表
    insertCardInfoTable(cursor,data_card_info)

    # 按照person_id2card_ids依次创建并插入cardNo表
    insertCardNoTable(cursor,person_id2card_ids,card_id2card_name)


    # 提交事务
    cnx.commit()

    # 关闭数据库
    closeDataBase(cnx,cursor)

    return "导入数据库成功！"

if __name__ == "__main__":

    #单次读取写入数据库
    # writeToDataBase()

    # 计算时间
    start = time.time()
    findCardInfoByCNQQ('寒茹怜', '')
    end = time.time()
    print("耗时{}秒".format(end - start))