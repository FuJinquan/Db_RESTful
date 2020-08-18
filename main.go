package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Note note struct
type Note struct {
	gorm.Model
	Title string `json:"title"`
	Text  string `json:"text"`
}

//Response response struct
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"messsage"`
}

//StandardError 标准错误，包含错误码和错误信息
type StandardError struct {
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errormsg"`
}

var (
	//ErrSuccess 操作成功
	ErrSuccess = StandardError{0, "成功"}
	//ErrUnrecognized 未知错误
	ErrUnrecognized = StandardError{-1, "未知错误"}
	//ErrTitleIsNil 标题为空错误
	ErrTitleIsNil = StandardError{1000, "标题为空错误！"}
	//ErrRecordIsNil 记录为空错误
	ErrRecordIsNil = StandardError{1001, "记录为空错误"}
	//ErrDbIsNil 数据库为空错误
	ErrDbIsNil = StandardError{1002, "数据库为空错误！"}
	//ErrCreateTable 创建数据表错误
	ErrCreateTable = StandardError{2000, "创建数据表错误！"}
	//ErrCreateDb 创建数据库失败
	ErrCreateDb = StandardError{2001, "创建数据库失败！"}
	//ErrSearch 查询数据库错误
	ErrSearch = StandardError{3000, "查询数据库错误!"}
	//ErrDelete 删除记录错误
	ErrDelete = StandardError{4000, "删除记录错误！"}
	//ErrUpdate 更新记录错误
	ErrUpdate = StandardError{5000, "记录更新错误！"}
)

//IsContain 返回一个bool值，判断某元素是否在数组中
func IsContain(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func main() {

	//连接mysql数据库
	db, err := gorm.Open("mysql", "root:950120fjq@/note?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	//设置数据库表名为单数
	db.SingularTable(true)
	defer db.Close()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//建立note表
	if err := db.AutoMigrate(&Note{}).Error; err != nil {
		panic("创建数据库失败！")
	}
	r := gin.Default()
	//添加
	//非id字段用body传
	r.POST("/note", func(c *gin.Context) {
		title := c.PostForm("title")
		text := c.PostForm("text")

		//错误是业务逻辑的一部分
		// 错误逻辑
		if len(title) < 1 {
			response := Response{
				Code:    ErrTitleIsNil.ErrorCode,
				Data:    nil,
				Message: ErrTitleIsNil.ErrorMsg,
			}
			c.JSON(422, response)
			return
		}

		note := Note{
			Title: title,
			Text:  text,
		}
		//添加记录
		//正常逻辑
		//正常逻辑下的错误逻辑
		if err := db.Create(&note).Error; err != nil {
			//这里最好不要用panic，panic主要是用于当前错误会使整个程序宕掉的时候
			// panic("记录创建失败！")
			response := Response{
				Code:    ErrCreateTable.ErrorCode,
				Data:    nil,
				Message: ErrCreateTable.ErrorMsg,
			}
			c.JSON(422, response)
			return
		}
		//正常逻辑下的正确逻辑
		response := Response{
			Code:    ErrSuccess.ErrorCode,
			Data:    note,
			Message: ErrSuccess.ErrorMsg,
		}
		//状态码201，用户新建或修改数据成功
		c.JSON(201, response)
	})
	//查找
	//按id查找
	r.GET("/note/:id", func(c *gin.Context) {
		var note Note
		id := c.Param("id")
		intID, _ := strconv.Atoi(id)
		//intID必须是整形的
		if err := db.Find(&note, intID).Error; err != nil {
			//panic("查询数据库失败")
			response := Response{
				Code:    ErrSearch.ErrorCode,
				Data:    nil,
				Message: ErrSearch.ErrorMsg,
			}
			c.JSON(422, response)
			return
		}
		response := Response{
			Code:    ErrSuccess.ErrorCode,
			Data:    note,
			Message: ErrSuccess.ErrorMsg,
		}
		c.JSON(200, response)

	})
	//查找所有记录
	r.GET("/note", func(c *gin.Context) {
		var notes []Note
		if err := db.Find(&notes).Error; err != nil {
			//panic("查询数据库失败")
			response := Response{
				Code:    ErrSearch.ErrorCode,
				Data:    nil,
				Message: ErrSearch.ErrorMsg,
			}
			c.JSON(500, response)
			return
		}
		if len(notes) < 1 {
			response := Response{
				Code:    ErrDbIsNil.ErrorCode,
				Data:    nil,
				Message: ErrDbIsNil.ErrorMsg,
			}
			c.JSON(404, response)
			return
		}
		for _, items := range notes {
			response := Response{
				Code:    ErrSuccess.ErrorCode,
				Data:    items,
				Message: ErrSuccess.ErrorMsg,
			}
			c.JSON(200, response)

		}

	})
	//删除
	r.DELETE("/note/:id", func(c *gin.Context) {
		var note Note
		id := c.Param("id")
		intID, _ := strconv.Atoi(id)
		note.ID = uint(intID)
		//添加Unscoped()方法，直接在数据库内删除记录。而不是软删除
		if err := db.Unscoped().Delete(&note).Error; err != nil {
			//panic("查询记录失败")
			response := Response{
				Code:    ErrDelete.ErrorCode,
				Data:    nil,
				Message: ErrDelete.ErrorMsg,
			}
			c.JSON(500, response)
			return
		}
		//fmt.Println(note)
		response := Response{
			Code:    ErrSuccess.ErrorCode,
			Data:    note,
			Message: ErrSuccess.ErrorMsg,
		}
		c.JSON(200, response)

	})
	//更新
	//非id字段用body传输
	r.PUT("/note/:id", func(c *gin.Context) {
		var note Note
		id := c.Param("id")
		title := c.PostForm("title")
		text := c.PostForm("text")
		intID, _ := strconv.Atoi(id)
		note.ID = uint(intID)
		if err := db.Model(&note).Updates(Note{Title: title, Text: text}).Error; err != nil {
			//panic("数据更新失败")
			response := Response{
				Code:    ErrUpdate.ErrorCode,
				Data:    note,
				Message: ErrUpdate.ErrorMsg,
			}
			c.JSON(500, response)
			return
		}
		//fmt.Println(note)
		response := Response{
			Code:    ErrSuccess.ErrorCode,
			Data:    note,
			Message: ErrSuccess.ErrorMsg,
		}
		//状态码201，用户新建或修改数据成功
		c.JSON(200, response)

	})
	r.Run()
}
