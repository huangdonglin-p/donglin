package main

import (
	"context"
	"donglin.framework.use/framework"
	"fmt"
	"log"
	"net/http"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(time.Second))
	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		// Do real action
		time.Sleep(3 * time.Second)
		fmt.Println("running。。。。")
		c.Json(http.StatusOK, "ok")

		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(http.StatusInternalServerError, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(http.StatusInternalServerError, "time out")
		c.SetHasTimeout()
	}
	return nil
}
