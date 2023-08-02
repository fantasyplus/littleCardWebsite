import openpyxl
import re
from os import path
import json


def processString(string):
    pattern = r"cn\+群内qq:(.*)"
    result = re.search(pattern, string)
    cn_qq_str = ""
    cn = ""
    qq = ""

    if result:
        cn_qq_str = result.group(1)
    else:  # 如果没有cn+群内qq:，则直接返回，第一行标题的特殊处理
        return cn, qq

    # 特殊情况，cn也全是数字
    pattern = r"^(\d+)[\+\＋](\d+)$"
    match = re.match(pattern, cn_qq_str)
    if match:
        cn = match.group(1)
        qq = match.group(2)
        return cn, qq

    pattern = r"(.*?)\s*(?=\d{6,})"
    result = re.findall(pattern, cn_qq_str)

    cn = result[0]
    # 去除末尾的 + 号
    if cn[-1] == "+" or cn[-1] == "＋":
        cn = cn[:-1]

    qq = re.search(r"\d{6,}", cn_qq_str).group()

    return cn, qq


def processRow(row, single_sheet_data):
    cn = ""
    qq = ""
    num = 0
    # 单种谷子的单行的数据
    row_data = []
    for i in range(len(row)):
        # 读取第一列，处理谷子和qq
        if i == 0:
            cn, qq = processString(row[i].value)

            # 标题行
            if cn == "" and qq == "":
                row_data.append(row[i].value)
            # cn，qq行
            else:
                row_data.extend([cn, qq])
        # 读取第二列，处理数量
        if i == 1 and row[i].value is not None and isinstance(row[i].value, int):
            num = int(row[i].value)
            row_data.append(num)
        elif i == 1 and isinstance(row[i].value, str):
            # 可能是“数量”这两个字（在第一行），是中文的话直接插入
            pattern = re.compile(r"[\u4e00-\u9fff]+")
            is_chinese = bool(pattern.search(row[i].value))
            if is_chinese:
                row_data.append(row[i].value)
            # 不是的话就是可能为字符串形式的数字，转换为int
            else:
                row_data.append(int(row[i].value))
        elif i == 1 and row[i].value is None:
            return
        # 读取第三列，处理状态
        if i == 2:
            if row[i].value is None:
                row_data.append("none")
            else:
                row_data.append(row[i].value)

    single_sheet_data.append(row_data)


def readSellInfo(excel_name):
    p = path.dirname(__file__) + "/../test_excel/" + excel_name
    # 读取Excel文件
    wb = openpyxl.load_workbook(p, data_only=True)
    sheet_names = wb.sheetnames
    # sheet_names=[wb.sheetnames[16]]
    print(sheet_names)

    # 一个表格的所有谷子信息
    # 每一个元素对应一种谷子的信息
    # 如果长度为1，说明是一个谷子对应一个角色的情况（正常情况）
    sheetdatas = []

    for sheet_name in sheet_names:
        sheet = wb[sheet_name]  # 修改为实际的工作表名

        # 获取状态所在的列
        condition_index = 0
        for content in sheet[1]:
            if content.value == "状态":
                condition_index = content.col_idx
                break

        # 一个谷子对应多个角色的情况
        # 一行大于三列：谷子名，角色1，角色2....角色N，数量，状态
        if condition_index > 3:
            # 表示要循环几次，即有一对多的情况下有几个角色
            count_index = condition_index - 3  # 比如5-3=2

            for i in range(count_index + 1):
                # 单种谷子的数据
                single_sheet_data = []
                name = None
                character = None
                # 读取整个工作表的数据
                for j, row in enumerate(sheet.iter_rows()):
                    # 如果到达某一行为空，则遍历完毕
                    if row[0].value is None:
                        break
                    # 如果是第一行，保存谷子名和角色名，跳过
                    if j == 0:
                        name = row[0].value
                        character = row[i + 1].value
                        continue

                    row = [row[0], row[i + 1], row[condition_index - 1]]
                    processRow(row, single_sheet_data)

                # 把第一行的角色名和谷子名加起来，再按照顺序加上card_id（19-1,19-2,19-3...）
                match = re.match(r"\d+", name)
                card_id_end_pos = match.end()

                name = (
                    name[:card_id_end_pos] + "_" + str(i + 1) + name[card_id_end_pos:]
                )
                name += "-" + character
                amount = "数量"

                single_sheet_data.insert(0, [name, amount])

                # 插入单个子表数据
                sheetdatas.extend(single_sheet_data)

        # 一行等于三列：谷子名，数量，状态
        elif condition_index == 3:
            # 单种谷子的数据
            single_sheet_data = []

            # 读取整个工作表的数据
            for row in sheet.iter_rows():
                # 如果到达某一行为空，则遍历完毕
                if row[0].value is None:
                    break

                processRow(row, single_sheet_data)

            # 插入单个子表数据
            sheetdatas.extend(single_sheet_data)

    # 关闭Excel文件
    wb.close()
    return sheetdatas


def writeJsonFile(json_name):
    # 分割每个谷子的数据，然后存入字典
    split_points = []
    for i in range(len(excel_data)):
        if len(excel_data[i]) == 2:
            # print(excel_data[i])
            split_points.append(i)
        if len(excel_data[i]) == 3:
            # print(excel_data[i])
            split_points.append(i)

    dict_data = {}
    for i in range(len(split_points) - 1):
        split_index = split_points[i]
        id_and_name = excel_data[split_index][0]

        # 如果长度为2，说明是标题行，获取了card_name和card_id就进行下一次循环
        match1 = re.search(r"\d+", id_and_name)
        # 一个谷子多个角色的情况
        match2 = re.search(r"\d+\_\d+", id_and_name)
        if match1 and match2 is None:
            card_id = match1.group()
            card_name = id_and_name[match1.end() :]

            dict_data[card_id] = excel_data[split_index : split_points[i + 1]]
        elif match2:
            card_id = match2.group()
            card_name = id_and_name[match2.end() :]
            dict_data[card_id] = excel_data[split_index : split_points[i + 1]]
    # print(dict_data)

    p = path.dirname(__file__) + "/../json/" + json_name

    # Convert list to JSON and save it to the file
    with open(p, "w", encoding="utf-8") as json_file:
        json.dump(dict_data, json_file, ensure_ascii=False, indent=4)


if __name__ == "__main__":
    file_name = "selldata_2023_08_02_1.xlsx"

    excel_data = readSellInfo(file_name)

    # writeJsonFile("selldata.json")
