import pandas as pd
from os import path


def read_excel_with_pandas(excel_name):
    p = path.dirname(__file__) + "/../test_excel/" + excel_name

    # 读取Excel文件中的所有sheet数据，返回一个字典，键是sheet名，值是对应的DataFrame
    all_sheets_data = pd.read_excel(p, sheet_name=None, engine="openpyxl")
    # "xlrd", "openpyxl", "odf", "pyxlsb"
    # 获取所有的sheet名
    sheet_names = list(all_sheets_data.keys())

    # 打印sheet名
    print(sheet_names)

    # 打印每个sheet的前几行数据示例
    for sheet_name in sheet_names:
        print(f"--- {sheet_name} ---")
        print(all_sheets_data[sheet_name].head())


if __name__ == "__main__":
    excel_name = "selldata_2023_08_02_1.xlsx"  # 替换为你的Excel文件名
    read_excel_with_pandas(excel_name)
