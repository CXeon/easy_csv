package easycsv

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

var dataItem = []string{
	"SN0abcd",
	"MAC00000",
	"MeshID00",
	"DevType00",
	"Product00",
	"Ver00",
	"ProtoMQTT",
	"TimeOL",
	"TimeOffline",
	"Phone",
	"Email",
	"STATE",
}

// Memo: 经测试验证200万行记录能够写进csv,并非之前规格定义说只能写入6万行，
// 只是用WPS打开文件会提示超过长度限制，在我电脑最终展示1048576行。用goland打开没有出现截断。
func TestCsvBase(t *testing.T) {
	fileName := "./TestCsvBase.csv"

	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	//构造数据测试写入
	rowLimit := 2000000 //限制总行数
	batch := 10000      //每批次写入行数
	rowBatch := make([][]string, 0)

	batchCount := 0 //批次计数
	for i := 0; i < rowLimit; i++ {
		num := i + 1
		numStr := strconv.Itoa(num)

		newDataItem := make([]string, len(dataItem))
		for j := 0; j < len(dataItem); j++ {
			newDataItem[j] = dataItem[j] + numStr
		}

		rowBatch = append(rowBatch, newDataItem)

		if num%batch == 0 || num == rowLimit {
			batchCount++
			t.Logf("当前写入批次[%d],写入数据行数[%d]\n", batchCount, len(rowBatch))

			err = writer.WriteAll(rowBatch)
			if err != nil {
				t.Error(err)
				return
			}
			writer.Flush()
			rowBatch = make([][]string, 0)
		}
	}

	return

}
