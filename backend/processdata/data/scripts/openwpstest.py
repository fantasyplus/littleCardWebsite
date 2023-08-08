# 第一步先import所需模块（rpcxxxapi，xxx为对应项目的名字）
# rpcwpsapi模块为WPS文字项目的开发接口
# rpcwppapi则是WPS演示的
# rpcetapi毫无疑问就是WPS表格的了
# 另外还有common模块，为前三者的公共接口模块，通常不能单独使用

# 调起WPS必需通过createXXXRpcInstance接口，所以导入它是必需的
# 以WPS文字为例
from time import sleep
from pywpsrpc.rpcetapi import (createEtRpcInstance,etapi)
from pywpsrpc import RpcIter

# 这里仅创建RPC实例
hr, rpc = createEtRpcInstance()

# 注意：
# WPS开发接口的返回值第一个总是HRESULT（无返回值的除外）
# 通常不为0的都认为是调用失败（0 == common.S_OK）
# 可以使用common模块里的FAILED或者SUCCEEDED去判断

# 通过rpc实例调起WPS进程
hr, app = rpc.getEtApplication()

app.Visible=False
workbooks = app.Workbooks

hr, workbook = workbooks.Open('/home/web/web/web-project/backend/processdata/data/scripts/selldata_2023_08_04_1.xlsx')
workbook.SaveAs('/home/web/web/web-project/backend/processdata/data/test_excel/selldata_2023_08_04_1.xlsx')
workbook.Close()
app.Quit()