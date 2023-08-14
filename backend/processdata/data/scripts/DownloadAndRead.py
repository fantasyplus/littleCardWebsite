import downloadexcel
import wpsInterface
import readCard
import readSell

#从腾讯文档下载excel表格并读取写入到json文件中
if __name__ == "__main__":
    sell_path=downloadexcel.DownloadSellExcel()
    wpsInterface.OpenAndSave(sell_path)
    readSell.ReadSell(sell_path)

    # card_path=downloadexcel.DownloadCardExcel()
    # wpsInterface.OpenAndSave(card_path)
    # readCard.ReadCard(card_path)