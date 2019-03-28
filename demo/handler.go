package main

import (
	"github.com/linghaihui/kuafu"
	"strconv"
)

type Book struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var bookList = make(map[int]Book)

// curl -X PUT '127.0.0.1:9999/add' -H 'content-type:application/json' -d '{"id":1, "name": "唐诗三百首"}'
func AddBook(ctx *kuafu.Context) {
	book := Book{}
	ctx.BindJSON(&book)
	bookList[book.Id] = book
	ctx.JsonResponse(200, book)
}

// curl '127.0.0.1:9999/get/1'
func GetBook(ctx *kuafu.Context) {
	id, _ := strconv.Atoi(ctx.GetParam("id"))
	if book, ok := bookList[id]; ok {
		ctx.JsonResponse(200, book)
	} else {
		log.Debug("cannot find book")
		ctx.ResponseAnyCode(404, "cannot find book")
	}
}

// curl -X DELETE '127.0.0.1:9999/delete/1'
func DeleteBook(ctx *kuafu.Context) {
	id, _ := strconv.Atoi(ctx.GetParam("id"))
	if book, ok := bookList[id]; ok {
		delete(bookList, id)
		ctx.JsonResponse(200, book)
	} else {
		ctx.Response404()
	}
}
