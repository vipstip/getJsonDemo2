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

type TypeData struct {
	Type        string
	ListTagData []TagData
}

type JsonResult struct {
	Title string
	Datas []TypeData
}

var Name = "name"
var Des = "des"
var Time = "time"
var View = "view"
var Vote = "vote"
var db *gorm.DB

func init() {
	db, _ = ConnectDB()
}

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
	channel := make(chan TypeData, 5)
	//Put data channel
	for i := 0; i < 5; i++ {
		go handlerJson(i, channel)
	}

	//Get data from channel

	var arrObject = make([]TypeData, 5)
	for i, _ := range arrObject {
		arrObject[i] = <-channel
	}

	var jsonResult JsonResult
	jsonResult.Title = "Get Data From Json"
	jsonResult.Datas = arrObject

	var b []byte
	b, _ = json.Marshal(jsonResult)

	return b
}

func handlerJson(stt int, channel chan TypeData) {
	switch stt {
	case 0:
		tagDatas, types := GetArrName()
		var typeData TypeData
		typeData.Type = types
		typeData.ListTagData = tagDatas
		channel <- typeData
	case 1:
		tagDatas, types := GetArrDes()
		var typeData TypeData
		typeData.Type = types
		typeData.ListTagData = tagDatas
		channel <- typeData
	case 2:
		tagDatas, types := GetArrTime()
		var typeData TypeData
		typeData.Type = types
		typeData.ListTagData = tagDatas
		channel <- typeData
	case 3:
		tagDatas, types := GetArrView()
		var typeData TypeData
		typeData.Type = types
		typeData.ListTagData = tagDatas
		channel <- typeData
	case 4:
		tagDatas, types := GetArrVote()
		var typeData TypeData
		typeData.Type = types
		typeData.ListTagData = tagDatas
		channel <- typeData
	}
}

func GetArrName() ([]TagData, string) {
	defer duration(track("GetArrName"))
	var data []Data
	db.Raw("SELECT id , name  FROM `recipes` WHERE `name` LIKE '%ga%' ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data)
	return tagDatas, Name
}

func GetArrDes() ([]TagData, string) {
	defer duration(track("GetArrDes"))

	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE description LIKE '%bánh%' ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data)
	return tagDatas, Des
}

func GetArrTime() ([]TagData, string) {
	defer duration(track("GetArrTime"))
	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE execution_time>60 AND execution_time<100 ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data)
	return tagDatas, Time
}

func GetArrView() ([]TagData, string) {
	defer duration(track("GetArrView"))

	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE view>1000 ORDER BY Rand() LIMIT 0,5").Scan(&data)
	tagDatas := getTagDatas(data)
	return tagDatas, View
}

func GetArrVote() ([]TagData, string) {
	defer duration(track("GetArrVote"))

	var data []Data
	db.Raw("SELECT id,name  FROM recipes WHERE vote>1000 ORDER BY Rand() LIMIT 0,5").Scan(&data)

	tagDatas := getTagDatas(data)
	return tagDatas, Vote
}

func GetArrTag(idTag int) []string {
	var listTag []string
	db.Table("tag").Select("tag.name").Joins("INNER JOIN relationship ON tag.id = relationship.tag_id INNER JOIN recipes ON recipes.id = relationship.recipes_id").Where("recipes.id = ?", idTag).Scan(&listTag)
	//db.Raw("SELECT tag.name FROM tag INNER JOIN relationship ON tag.id = relationship.tag_id INNER JOIN recipes ON recipes.id = relationship.recipes_id WHERE recipes.id = ?",idTag).Find(&listTag)
	return listTag
}

func handlerTag(idTag int, channel chan TagData, nameTagData string) {
	listTag := GetArrTag(idTag)
	var tagData TagData
	tagData.Name = nameTagData
	tagData.ListNameTag = listTag
	channel <- tagData
}

func getTagDatas(data []Data) []TagData {

	totalGetTag := len(data)
	channel := make(chan TagData, totalGetTag)
	//Put data channel
	for i := 0; i < totalGetTag; i++ {
		go handlerTag(data[i].Id, channel, data[i].Name)
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
