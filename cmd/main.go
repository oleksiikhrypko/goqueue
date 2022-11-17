package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	storage "goqueue/pkg/storage"
	klist "goqueue/pkg/storage/klist"
	listdb "goqueue/pkg/storage/leveldb"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/net/context"
)

var (
	Version   = ""
	Branch    = ""
	Commit    = ""
	BuildTime = ""
)

// for tests
func main() {
	now := time.Now()
	log.Printf("start version:%s branch:%s commit:%s buildTime:%s", Version, Branch, Commit, BuildTime)

	db, err := leveldb.OpenFile("tmp/test2", nil)
	defer db.Close()
	if err != nil {
		log.Println("init db failed", err)
		return
	}
	wg := &sync.WaitGroup{}
	ctx := context.Background()

	list, err := klist.New(ctx, "list", listdb.NewDB(ctx, db))
	if err != nil {
		log.Println("failed to init klist", err)
	}

	wg.Add(1)
	go addItemsToList(wg, list, 1000000)
	// time.Sleep(time.Second)
	// go readList(wg, list, "r1")
	// go readList(wg, list, "r1")
	// time.Sleep(100 * time.Microsecond)
	// go readList(wg, list, "r1")
	wg.Wait()
	log.Println("[write] done", time.Now().Sub(now))

	// wg.Add(1)
	// go readList(wg, list, "r1")
	// wg.Wait()
	log.Println("[read] done", time.Now().Sub(now))

	// readAllKeys(db)
}

func addItemsToList(wg *sync.WaitGroup, list storage.List, len int) {
	for i := 0; i < len; i++ {
		err := list.Add([]byte(fmt.Sprintf("key" + strconv.Itoa(i))))
		if err != nil {
			log.Println("add item failed", err)
		}
		// log.Println("idx:", "i:", i)
	}
	wg.Done()
}

func readList(wg *sync.WaitGroup, list storage.List, fpfx string) {
	var err error
	item := list.GetFirst()
	for item != nil {
		log.Println(fpfx, "--->", string(item))
		item, err = list.GetNext(item)
		if err != nil {
			log.Println("failed to read next element from list", err)
			return
		}
	}
	wg.Done()
}

func readList2(wg *sync.WaitGroup, list storage.List, fpfx string) {
	item, err := list.Pop()
	if err != nil {
		log.Println("failed to read first element from list", err)
		return
	}
	for item != nil {
		log.Println(fpfx, "--->", string(item))
		item, err = list.Pop()
		if err != nil {
			log.Println("failed to read next element from list", err)
			return
		}
	}
	wg.Done()
}

func readListBack(wg *sync.WaitGroup, list storage.List, fpfx string) {
	var err error
	item := list.GetLast()
	for item != nil {
		log.Println(fpfx, "<---", string(item))
		item, err = list.GetPrev(item)
		if err != nil {
			log.Println("failed to read prev element from list", err)
			return
		}
	}
	wg.Done()
}

func readAllKeys(db *leveldb.DB) {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		log.Println("key", string(key), "value", string(value))
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Println("readAllKeys error:", err)
	}
}
