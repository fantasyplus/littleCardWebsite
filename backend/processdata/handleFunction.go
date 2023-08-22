package processdata

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func DownloadAndRead(name string) (card_data_path string, sell_data_path string) {
	// 设置要执行的Python脚本路径和参数
	pythonScript := "./data/scripts/" + name
    cmd := exec.Command("python3", pythonScript)
    
    // 创建管道来捕获标准输出
    stdoutPipe, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println("无法创建标准输出管道:", err)
        return
    }
    
    // 开始执行命令
    err = cmd.Start()
    if err != nil {
        fmt.Println("执行Python程序时出错:", err)
        return
    }
    
	card_data_path = ""
	sell_data_path = ""
    // 逐行读取标准输出并解析文件路径
    scanner := bufio.NewScanner(stdoutPipe)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "filename:") {
            parts := strings.Split(line, ",")
            if len(parts) == 2 {
                filename := strings.TrimPrefix(parts[0], "filename:")
                filepath := strings.TrimPrefix(parts[1], "filepath:")
				if strings.Contains(filename, "carddata") {
					card_data_path = filepath
				} else if strings.Contains(filename, "selldata") {
					sell_data_path = filepath
				}
				// fmt.Println(filename, filepath)
            }
        }
    }
    
    // 等待命令执行完毕
    err = cmd.Wait()
    if err != nil {
        fmt.Println("命令执行完毕时出错:", err)
    }

	return card_data_path, sell_data_path
}

func UpdateDataBase(){
	db := ConnectDB()
	CreateTable(db)
	person_id2card_ids2card_num, card_id2card_name, card_id2person_ids2status := InsertPersonInfoTable(db)
	InsertCardIndexTable(db, person_id2card_ids2card_num)
	InsertCardInfoTable(db)
	InsertCardNoTable(db, person_id2card_ids2card_num, card_id2card_name, card_id2person_ids2status)
}

