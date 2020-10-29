package GetData

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type Data struct {
	Id   int
	Name string
}

type TagData struct {
	Name        string
	ListNameTag []string
}

type JsonResult struct {
	Name string
	Des  string
	Time int
	View int
	Vote int
}

var Name = "name"
var Des = "des"
var Time = "time"
var View = "view"
var Vote = "vote"

func ConnectDB() (*gorm.DB, error) {
	dbUser := "root"
	dbPass := ""
	dbName := "hnag"
	dsn := dbUser + ":" + dbPass + "@/" + dbName
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Connect errr")
	}

	return db, err
}

func GetJson() []byte {
	defer duration(track("GetJSon"))

	db, _ := ConnectDB()
	channel := make(chan map[string][]TagData, 5)
	//Put data channel
	for i := 0; i < 5; i++ {
		go handlerJson(i, channel, db)
	}

	//Get data from channel

	var arrObject = make([]map[string][]TagData, 5)
	for i, _ := range arrObject {
		arrObject[i] = <-channel
	}

	var b []byte
	b, _ = json.Marshal(arrObject)

	return b
}

func handlerJson(stt int, channel chan map[string][]TagData, db *gorm.DB) {
	switch stt {
	case 0:
		myMap := make(map[string][]TagData)
		tagDatas, key := GetArrName(db)
		myMap[key] = tagDatas
		channel <- myMap
	case 1:
		myMap := make(map[string][]TagData)
		tagDatas, key := GetArrDes(db)
		myMap[key] = tagDatas
		channel <- myMap
	case 2:
		myMap := make(map[string][]TagData)
		tagDatas, key := GetArrTime(db)
		myMap[key] = tagDatas
		channel <- myMap
	case 3:
		myMap := make(map[string][]TagData)
		tagDatas, key := GetArrView(db)
		myMap[key] = tagDatas
		channel <- myMap
	case 4:
		myMap := make(map[string][]TagData)
		tagDatas, key := GetArrVote(db)
		myMap[key] = tagDatas
		channel <- myMap
	}
}

func GetArrName(db *gorm.DB) ([]TagData, string) {
	defer duration(track("GetArrName"))
	var data []Data
	db.Raw("SELECT id , name  FROM `recipes` WHERE `name` LIKE '%ga%' ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data, db)
	return tagDatas, Name
}

func GetArrDes(db *gorm.DB) ([]TagData, string) {
	defer duration(track("GetArrDes"))

	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE description LIKE '%bánh%' ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data, db)
	return tagDatas, Des
}

func GetArrTime(db *gorm.DB) ([]TagData, string) {
	defer duration(track("GetArrTime"))
	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE execution_time>60 AND execution_time<100 ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data, db)
	return tagDatas, Time
}

func GetArrView(db *gorm.DB) ([]TagData, string) {
	defer duration(track("GetArrView"))

	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE view>1000 ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data, db)
	return tagDatas, View
}

func GetArrVote(db *gorm.DB) ([]TagData, string) {
	defer duration(track("GetArrVote"))

	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE vote>1000 ORDER BY Rand() LIMIT 0,5").Scan(&data)

	tagDatas := getTagDatas(data, db)
	return tagDatas, Vote
}

func GetArrTag(idTag int, db *gorm.DB) []string {
	var listTag []string
	db.Table("tag").Select("tag.name").Joins("INNER JOIN relationship ON tag.id = relationship.tag_id INNER JOIN recipes ON recipes.id = relationship.recipes_id").Where("recipes.id = ?", idTag).Scan(&listTag)
	//db.Raw("SELECT tag.name FROM tag INNER JOIN relationship ON tag.id = relationship.tag_id INNER JOIN recipes ON recipes.id = relationship.recipes_id WHERE recipes.id = ?",idTag).Find(&listTag)
	return listTag
}

func handlerTag(idTag int, channel chan TagData, nameTagData string, db *gorm.DB) {
	listTag := GetArrTag(idTag, db)
	var tagData TagData
	tagData.Name = nameTagData
	tagData.ListNameTag = listTag
	channel <- tagData
}

func getTagDatas(data []Data, db *gorm.DB) []TagData {

	totalGetTag := len(data)
	channel := make(chan TagData, totalGetTag)
	//Put data channel
	for i := 0; i < totalGetTag; i++ {
		go handlerTag(data[i].Id, channel, data[i].Name, db)
	}

	//Get data from channel

	var tagDatas = make([]TagData, totalGetTag)
	for i, _ := range tagDatas {
		tagDatas[i] = <-channel
	}
	return tagDatas
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
