# -*- coding: UTF-8 -*-

import json
import os
from os import path
import re
import time
from datetime import datetime
from time import sleep
import click
import pandas as pd
import requests
from bs4 import BeautifulSoup
from pywpsrpc.rpcetapi import createEtRpcInstance, etapi

#关闭https警告
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


class TengXunDocument:
    def __init__(self, document_url, local_pad_id, cookie_value):
        # excel文档地址
        self.document_url = document_url
        # 此值每一份腾讯文档有一个,需要手动获取
        self.localPadId = local_pad_id
        self.headers = {
            "content-type": "application/x-www-form-urlencoded",
            "Cookie": cookie_value,
        }

    def get_now_user_index(self):
        """
        # 获取当前用户信息,供创建下载任务使用
        :return:
            # nowUserIndex = '4883730fe8b94fbdb94da26a9a63b688'
            # uid = '144115225804776585'
            # utype = 'wx'
        """
        response_body = requests.get(
            url=self.document_url, headers=self.headers, verify=False
        )
        parser = BeautifulSoup(response_body.content, "html.parser")
        global_multi_user_list = re.findall(
            re.compile("window.global_multi_user=(.*?);"), str(parser)
        )
        if global_multi_user_list:
            user_dict = json.loads(global_multi_user_list[0])
            print(user_dict)
            return user_dict["nowUserIndex"]
        return "cookie过期,请重新输入"

    def export_excel_task(self, export_excel_url):
        """
        导出excel文件任务,供查询文件数据准备进度
        :return:
        """
        body = {"docId": self.localPadId, "version": "2"}

        res = requests.post(
            url=export_excel_url, headers=self.headers, data=body, verify=False
        )
        operation_id = res.json()["operationId"]
        return operation_id

    def download_excel(self, check_progress_url,path):
        """
        下载excel文件
        :return:
        """
        # 拿到下载excel文件的url
        start_time = time.time()
        file_url = ""
        while True:
            res = requests.get(
                url=check_progress_url, headers=self.headers, verify=False
            )
            progress = res.json()["progress"]
            if progress == 100:
                file_url = res.json()["file_url"]
                break
            elif time.time() - start_time > 30:
                print("数据准备超时,请排查")
                break
        if file_url:
            self.headers["content-type"] = "application/octet-stream"
            res = requests.get(url=file_url, headers=self.headers, verify=False)
            with open(path, "wb") as f:
                f.write(res.content)
            print("filename:{},filepath:{}".format(os.path.basename(path),path))
        else:
            print("下载文件地址获取失败, 下载excel文件不成功")


def renameExcel(name):
    p = path.dirname(__file__) + "/../test_excel/"+name+"/"
    file_name=""

    # 获取文件夹内所有文件和子文件夹
    files = os.listdir(p)

    current_datetime = datetime.strftime(datetime.now(), "%Y_%m_%d")
    if len(files) == 0:
        cnt = 1
        file_name = f"{name}_{current_datetime}_{cnt}.xlsx"
    else:
        cnt = -1
        # 找到和当前日期相同的文件，在它的基础上次数+1
        # print(files)
        # 先找到所有包含当前日期的文件
        for temp_name in files:
            if temp_name.find(current_datetime) != -1:
                # 找到最后一个下划线的位置
                temp_cnt_index = temp_name.rfind("_")
                if temp_cnt_index != -1:
                    # 提取最后一个下划线后的部分
                    temp_cnt = temp_name[temp_cnt_index + 1 :]

                # 取出序号
                temp_cnt = int(temp_cnt[: temp_cnt.find(".xlsx")])
                # 找到最大的那个
                if temp_cnt > cnt:
                    cnt = temp_cnt
                

        if cnt != -1:
            cnt = int(cnt) + 1
            print("max cnt :{}".format(cnt))
            file_name = f"{name}_{current_datetime}_{cnt}.xlsx"
        #如果没有，那就是今天第一次下载，创建新的文件
        else:
            cnt = 1
            file_name = f"{name}_{current_datetime}_{cnt}.xlsx"
    return p + file_name

def DownloadSellExcel_old():
    # excel文档地址
    document_url = "https://docs.qq.com/sheet/DRmVDY3l5aldRcGNM"
    # local_pad_id和cookie_value来自doc_info
    # 此值每一份腾讯文档有一个,需要手动获取
    local_pad_id = "300000000$FeCcyyjWQpcL"
    # 打开腾讯文档后,从抓到的接口中获取cookie信息(先按F5刷新，F12才能检测到doc_info)
    cookie_value = "RK=Rfn5jgFNTO; ptcz=aeb83b496ba5c5b68290e157192714bd4e92c996af50fc3a10d9009db743fcd9; traceid=1d6ca2be6c; TOK=1d6ca2be6caed29d; hashkey=1d6ca2be; fingerprint=1b5db134e35c4d438fd72ae972f6260312; ES2=30abdb2bffd9d4a5; low_login_enable=1; CheckKey=f5968ee7b92aa5ac207457d5FY0Rq1; uin=o0530925206; skey=@bzeimZVUB; luin=o0530925206; lskey=000100000d26691cfad4c1eab0acc14e8cd83322b12b6439dc26155e98cdc6e402bb2d7e835b21680d134d76; p_uin=o0530925206; pt4_token=M4jYG7I7zoSTGNFgdFvbHzHbT5IRWlMSeYrnd6NpmKA_; p_skey=wr6JW0E-BbHZM7AeII4HCmi6GL7josuc8MhA*WAPat8_; p_luin=o0530925206; p_lskey=000400008f9f14c64ab56f9a0cc7df9412f1bcb6b9d51b54ffe59964d8a278107b1e1b35fa5af3832cdbf68a; uid=144115210445368678; utype=qq; vfwebqq=68fefae0ac21e1c56d14078d85b218bc6e2ed4a26b2a5c6b6d8d0a69d7c7b887dcb2b62d43e365bc; DOC_QQ_APPID=101458937; DOC_QQ_OPENID=E0AAA7BC62A9988A22F6F16AD7B3FC90; env_id=gray-no2; gray_user=true; DOC_SID=3782b60a0dda49a688d535c3d5b58a518b65bc063c7f45f7bd6da3dfc1983865; SID=3782b60a0dda49a688d535c3d5b58a518b65bc063c7f45f7bd6da3dfc1983865; uid_key=EOP1mMQHGixzNnRqT3NDeitTWC9ZSVRxRDNFTjlDU3RsNHprUVBUQlltaWVlM1FWTGZzPSJIvzj1gcn0GCGPPoWEfTC0%2BOmVrXGJHTGPiyC56fGV7Y3KkzQWoIvNzCZMCot4PpY%2BqtqqGFmQshWHbTnMxxjInyPhxOwTuucgKNT4xKcG; loginTime=1690939770860; optimal_cdn_domain=docs2.gtimg.com; _tucao_custom_info=MmpKS2hqQlFTQnhNSGlidkowZU53d1lJaDF4T2Z3aXd2NUs4aHZmcnJsa0pkWmRBamsvUTlIT2F0c2RGdzJPUA%3D%3D--jZS3E40tufnZ4zwKmS7LPQ%3D%3D"
    tx = TengXunDocument(document_url, local_pad_id, cookie_value)
    now_user_index = tx.get_now_user_index()
    # 导出文件任务url，来自export_office接口
    export_excel_url = f"https://docs.qq.com/v1/export/export_office?u={now_user_index}"
    # 获取导出任务的操作id，来自query_process接口
    operation_id = tx.export_excel_task(export_excel_url)
    check_progress_url = f"https://docs.qq.com/v1/export/query_progress?u={now_user_index}&operationId={operation_id}"

    path = renameExcel("selldata")
    tx.download_excel(check_progress_url, path)

    return path

def DownloadSellExcel():
    # excel文档地址
    document_url = "https://docs.qq.com/sheet/DRlR6UHpuRlhBc2xk"
    # local_pad_id和cookie_value来自doc_info
    # 此值每一份腾讯文档有一个,需要手动获取
    local_pad_id = "300000000$FTzPznFXAsld"
    # 打开腾讯文档后,从抓到的接口中获取cookie信息(先按F5刷新，F12才能检测到doc_info)
    cookie_value = "RK=Rfn5jgFNTO; ptcz=aeb83b496ba5c5b68290e157192714bd4e92c996af50fc3a10d9009db743fcd9; fingerprint=1b5db134e35c4d438fd72ae972f6260312; low_login_enable=1; luin=o0530925206; lskey=000100000d26691cfad4c1eab0acc14e8cd83322b12b6439dc26155e98cdc6e402bb2d7e835b21680d134d76; p_luin=o0530925206; p_lskey=000400008f9f14c64ab56f9a0cc7df9412f1bcb6b9d51b54ffe59964d8a278107b1e1b35fa5af3832cdbf68a; uid=144115210445368678; utype=qq; DOC_QQ_APPID=101458937; DOC_QQ_OPENID=E0AAA7BC62A9988A22F6F16AD7B3FC90; env_id=gray-no2; gray_user=true; DOC_SID=3782b60a0dda49a688d535c3d5b58a518b65bc063c7f45f7bd6da3dfc1983865; SID=3782b60a0dda49a688d535c3d5b58a518b65bc063c7f45f7bd6da3dfc1983865; uid_key=EOP1mMQHGixzNnRqT3NDeitTWC9ZSVRxRDNFTjlDU3RsNHprUVBUQlltaWVlM1FWTGZzPSJIvzj1gcn0GCGPPoWEfTC0%2BOmVrXGJHTGPiyC56fGV7Y3KkzQWoIvNzCZMCot4PpY%2BqtqqGFmQshWHbTnMxxjInyPhxOwTuucgKNT4xKcG; loginTime=1690939770860; optimal_cdn_domain=docs2.gtimg.com; backup_cdn_domain=docs2.gtimg.com; traceid=96992fe11c; TOK=96992fe11c64bc51; hashkey=96992fe1; tgw_l7_route=e1342f1217734b6f33001f07b5ffd8d1"
    tx = TengXunDocument(document_url, local_pad_id, cookie_value)
    now_user_index = tx.get_now_user_index()
    # 导出文件任务url，来自export_office接口
    export_excel_url = f"https://docs.qq.com/v1/export/export_office?u={now_user_index}"
    # 获取导出任务的操作id，来自query_process接口
    operation_id = tx.export_excel_task(export_excel_url)
    check_progress_url = f"https://docs.qq.com/v1/export/query_progress?u={now_user_index}&operationId={operation_id}"

    path = renameExcel("selldata")
    tx.download_excel(check_progress_url, path)

    return path

def DownloadCardExcel():
    # excel文档地址
    document_url = "https://docs.qq.com/sheet/DRnNkVmlHdEtHV3Bt"
    # local_pad_id和cookie_value来自doc_info
    # 此值每一份腾讯文档有一个,需要手动获取
    local_pad_id = "300000000$FsdViGtKGWpm"
    # 打开腾讯文档后,从抓到的接口中获取cookie信息(先按F5刷新，F12才能检测到doc_info)
    cookie_value="RK=Rfn5jgFNTO; ptcz=aeb83b496ba5c5b68290e157192714bd4e92c996af50fc3a10d9009db743fcd9; fingerprint=1b5db134e35c4d438fd72ae972f6260312; low_login_enable=1; luin=o0530925206; lskey=000100000d26691cfad4c1eab0acc14e8cd83322b12b6439dc26155e98cdc6e402bb2d7e835b21680d134d76; p_luin=o0530925206; p_lskey=000400008f9f14c64ab56f9a0cc7df9412f1bcb6b9d51b54ffe59964d8a278107b1e1b35fa5af3832cdbf68a; uid=144115210445368678; utype=qq; DOC_QQ_APPID=101458937; DOC_QQ_OPENID=E0AAA7BC62A9988A22F6F16AD7B3FC90; env_id=gray-no2; gray_user=true; DOC_SID=3782b60a0dda49a688d535c3d5b58a518b65bc063c7f45f7bd6da3dfc1983865; SID=3782b60a0dda49a688d535c3d5b58a518b65bc063c7f45f7bd6da3dfc1983865; uid_key=EOP1mMQHGixzNnRqT3NDeitTWC9ZSVRxRDNFTjlDU3RsNHprUVBUQlltaWVlM1FWTGZzPSJIvzj1gcn0GCGPPoWEfTC0%2BOmVrXGJHTGPiyC56fGV7Y3KkzQWoIvNzCZMCot4PpY%2BqtqqGFmQshWHbTnMxxjInyPhxOwTuucgKNT4xKcG; loginTime=1690939770860; optimal_cdn_domain=docs2.gtimg.com; traceid=a681e0d6c0; TOK=a681e0d6c0379197; hashkey=a681e0d6; tgw_l7_route=2989f787ed3b9cae3fa89b54c6588a3a"
    tx = TengXunDocument(document_url, local_pad_id, cookie_value)
    now_user_index = tx.get_now_user_index()
    # 导出文件任务url，来自export_office接口
    export_excel_url = f"https://docs.qq.com/v1/export/export_office?u={now_user_index}"
    # 获取导出任务的操作id，来自query_process接口
    operation_id = tx.export_excel_task(export_excel_url)
    check_progress_url = f"https://docs.qq.com/v1/export/query_progress?u={now_user_index}&operationId={operation_id}"

    path = renameExcel("carddata")
    tx.download_excel(check_progress_url, path)

    return path

if __name__ == "__main__":
    DownloadSellExcel()