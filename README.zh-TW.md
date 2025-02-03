# ShardingMap
[English](README.md)

ShardingMap 是一個簡單高效的 Go 實現的分片哈希映射。它允許將鍵值對分布到多個分片中，以實現線程安全和高效的訪問。

## 特性

- **分片**: 將資料分布到多個分片中，以提高可擴展性和並發性。
- **線程安全**: 使用 `sync.RWMutex` 確保安全的並發讀取和寫入。
- **可自訂**: 可設置分片數量和自訂分片函數。
- **FNV 哈希**: 預設使用 FNV-1a 哈希作為分片函數。

## 安裝

```bash
go get github.com/lishank0119/shardingmap
```

## 用法

### 創建一個分片映射

```go
package main

import (
    "fmt"
    "github.com/lishank0119/shardingmap"
)

func main() {
    // 創建一個新的分片映射，帶有自訂選項
    sm := shardingmap.New[int, string](
        shardingmap.WithShardCount[int, string](32), // 設置分片數量
        shardingmap.WithShardingFunc[int, string](func(key int) int { return key }), // 自訂分片函數
    )

    // 設置一些鍵值對
    sm.Set(1, "value1")
    sm.Set(2, "value2")

    // 獲取值
    if value, ok := sm.Get(1); ok {
        fmt.Println(value) // 輸出: value1
    }

    // 刪除鍵
    sm.Delete(2)

    // 遍歷所有鍵值對
    sm.ForEach(func(k int, v string) {
        fmt.Printf("%d: %s\n", k, v)
    })

    fmt.Println("總數:", sm.Len()) // 輸出: 總數: 1
}
```

### 方法

- **New**: 創建一個新的分片映射，可以配置選項（`WithShardCount` 和 `WithShardingFunc`）。
- **Set**: 向映射中添加或更新鍵值對。
- **Get**: 獲取與鍵相關的值。
- **Delete**: 刪除映射中的鍵值對。
- **Len**: 返回映射中鍵值對的總數。
- **ForEach**: 遍歷所有鍵值對。

## 自訂

您可以自訂以下內容：

- **分片數量**: 定義分片數量（預設值為 16）。
- **分片函數**: 提供自訂的分片函數，以便進行更複雜的資料分配。

## 授權

本專案採用 MIT 授權 - 詳細內容請參見 [LICENSE](LICENSE) 檔案。
