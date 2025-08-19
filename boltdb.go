package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// 创建数据库文件
func creatdb(nodeID string) *bolt.DB {
	databaseFolder := "database"
	// 指定数据库文件的完整路径
	dbPath := databaseFolder + "/" + nodeID + ".db"
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
		// return nil, err // 返回 nil 和错误，而不是错误本身
		return nil
	}
	// defer db.Close() // 确保在函数返回前关闭数据库
	// return db, nil // 成功创建数据库后返回数据库实例和nil错误
	return db
}

// 写入数据到数据库
// func writeData(db *bolt.DB, key string, value []byte) error {

// // 记录开始时间
// start := time.Now()
// err := db.Update(func(tx *bolt.Tx) error {
// 	// 创建或获取名为 "bucket1" 的 bucket
// 	b, err := tx.CreateBucketIfNotExists([]byte("bucket1"))
// 	if err != nil {
// 		return err
// 	}
// 	// 将键值对写入到 bucket 中
// 	err = b.Put([]byte(key), value)
// 	return err
// })
// if err != nil {
// 	return err
// }
// // 记录结束时间
// end := time.Now()
// // 计算延迟
// latency := end.Sub(start)
// fmt.Printf("节点把请求写入db数据库的延迟为: %s\n", latency)
// return nil

func writeData(db *bolt.DB, key string, value []byte) (t time.Duration, er error) {
	// 记录开始时间
	start := time.Now()
	err := db.Update(func(tx *bolt.Tx) error {
		// 创建或获取名为 "bucket1" 的 bucket
		b, err := tx.CreateBucketIfNotExists([]byte("bucket1"))
		if err != nil {
			return err
		}
		// 将键值对写入到 bucket 中
		err = b.Put([]byte(key), value)
		return err
	})
	if err != nil {
		return 0, err
	}

	// 记录结束时间
	end := time.Now()
	// 计算延迟
	latency := end.Sub(start)
	// 输出延迟
	// fmt.Printf("把接受到的消息写入.db数据库延迟为: %s\n", latency)
	fmt.Println("把接受到的消息写入.db数据库延迟为" + fmt.Sprintf("%dms", latency.Milliseconds()))

	return latency, nil
}

// 从数据库读取数据
func readData(db *bolt.DB, key string) (string, error) {
	var value string

	err := db.View(func(tx *bolt.Tx) error {
		// 获取名为 "bucket1" 的 bucket
		b := tx.Bucket([]byte("bucket1"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// 从 bucket 中读取键对应的值
		v := b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("key not found")
		}

		// 将字节切片转换为字符串
		value = string(v)
		return nil
	})

	return value, err
}

// 从数据库删除数据
func deleteData(db *bolt.DB, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		// 获取名为 "bucket1" 的 bucket
		b := tx.Bucket([]byte("bucket1"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// 删除键值对
		return b.Delete([]byte(key))
	})
}
