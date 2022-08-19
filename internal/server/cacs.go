package server

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"log"
	"strings"

	"github.com/dollarkillerx/2password/internal/pkg/errs"
	"github.com/dollarkillerx/2password/internal/pkg/response"
	"github.com/dollarkillerx/2password/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) all(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	newName, _ := getAstName(model.Account)
	resp, err := s.cos.Object.Get(context.Background(), newName, nil)
	if err != nil {
		response.Return(ctx, errs.NotData)
		return
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Return(ctx, errs.NotData)
		return
	}
	response.Return(ctx, string(bs))
}

func getAstName(key string) (string, string) {
	return fmt.Sprintf("2password/%s.json", key), fmt.Sprintf("2password/%s.json.bak", key)
}

type PData struct {
	Data string `json:"data"`
}

func (s *Server) update(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	newName, bakName := getAstName(model.Account)

	var bakData string
	resp, err := s.cos.Object.Get(context.Background(), newName, nil)
	if err == nil {
		defer resp.Body.Close()
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			response.Return(ctx, errs.NotData)
			return
		}
		bakData = string(bs)
	}

	var payload PData
	err = ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(404, gin.H{})
		return
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "application/json",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}

	if bakData != "" {
		_, err = s.cos.Object.Put(context.Background(), bakName, strings.NewReader(bakData), opt)
		if err != nil {
			log.Println(err)
			response.Return(ctx, errs.BadRequest)
			return
		}
	}

	if payload.Data != "" {
		_, err = s.cos.Object.Put(context.Background(), newName, strings.NewReader(payload.Data), opt)
		if err != nil {
			log.Println(err)
			response.Return(ctx, errs.BadRequest)
			return
		}
	}

	response.Return(ctx, gin.H{})
}
