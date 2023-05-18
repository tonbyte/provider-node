package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tonbyte/provider-node/config"
	"github.com/tonbyte/provider-node/datatype"
	"github.com/tonbyte/provider-node/storage"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type handler struct {
	pub     ed25519.PublicKey
	priv    ed25519.PrivateKey
	payload map[string]datatype.Payload
}

func newHandler(pub ed25519.PublicKey, priv ed25519.PrivateKey) *handler {
	h := handler{
		pub:     pub,
		priv:    priv,
		payload: make(map[string]datatype.Payload),
	}
	go h.worker()
	return &h
}

func (h *handler) worker() {
	for {
		<-time.NewTimer(time.Minute).C
		for k, v := range h.payload {
			if time.Now().Unix() > v.ExpirtionTime {
				delete(h.payload, k)
			}
		}
	}
}

func (h *handler) Status(c echo.Context) error {
	log.Info("/status")

	providerInfo, err := storage.GetProviderInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can not get provider info")
	}

	providerParams, err := storage.GetProviderParams()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "can not get provider params")
	}

	result := datatype.Status{
		Version:           Version,
		RatePerMbDay:      providerParams.RatePerMbDay,
		Span:              providerParams.MaxSpan,
		MinFileSize:       providerParams.MinimalFileSize,
		MaxFileSize:       providerParams.MaximalFileSize,
		ContractsCount:    providerInfo.ContractsCount,
		MaxContractsCount: providerInfo.ProviderConfig.MaxContracts,
		UsedSpace:         providerInfo.ContractsTotalSize,
		MaxTotalSize:      providerInfo.ProviderConfig.MaxTotalSize,
		ContractAddress:   config.StorageConfig.ContractAddress,
		HasGateway:        config.StorageConfig.HasGateway,
	}

	return c.JSON(http.StatusOK, echo.Map{
		"result": result,
	})
}

func (h *handler) UploadFile(c echo.Context) error {
	log.Info("/uploadFile")
	providerContractAddress := c.QueryParam("contract_address")

	// Do file checks.
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid form-data")
	}
	formDataFile := form.File["file"]
	if len(formDataFile) < 1 {
		return c.JSON(http.StatusBadRequest, "no files")
	}
	providerInfo, err := storage.GetProviderInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "server error")
	}
	maxSize, err := strconv.ParseInt(providerInfo.ProviderConfig.MaxTotalSize, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "server error")
	}
	usedSize, err := strconv.ParseInt(providerInfo.ContractsTotalSize, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "server error")
	}
	spaceLeft := maxSize - usedSize
	if spaceLeft < formDataFile[0].Size {
		return c.JSON(http.StatusBadRequest, "no space left")
	}

	// Save file.
	// TODO: timer to check and delete file if no contract created.
	storagePath := filepath.Join(config.StorageConfig.StorageDBPath, "contracts-storage")
	os.Mkdir(storagePath, 0750)
	fileInfo, err := saveFile(storagePath, spaceLeft, formDataFile)
	if err != nil {
		log.Warn(err.Error())
		os.Remove(fileInfo.FullPath)

		return c.JSON(http.StatusBadRequest, "can not save file")
	}

	// Add file to storage.
	newBagId := storage.CreateBagOutput(fileInfo.FullPath)
	if len(newBagId) != 64 || strings.ContainsAny(newBagId, "[]:-!?,.") {
		os.Remove(fileInfo.FullPath)
		return c.JSON(http.StatusBadRequest, "can not create bag")
	}

	// Generate storage contract message.
	rand.Seed(time.Now().UnixNano())
	queryID := fmt.Sprint(rand.Uint32())
	messagesFolder := filepath.Join(config.StorageConfig.StorageDBPath, "messages-storage")
	os.Mkdir(messagesFolder, 0750)
	messagePath := filepath.Join(messagesFolder, newBagId+"_"+queryID)
	if !storage.NewContractMessage(newBagId, messagePath, queryID, providerContractAddress) {
		storage.RemoveBag(newBagId)
		return c.JSON(http.StatusBadRequest, "can not create message")
	}

	// Generate transaction.
	messageBody, err := os.ReadFile(messagePath)
	if err != nil {
		storage.RemoveBag(newBagId)
		return c.JSON(http.StatusBadRequest, "can not read message")
	}
	messageCell, err := cell.FromBOC(messageBody)
	if err != nil {
		storage.RemoveBag(newBagId)
		return c.JSON(http.StatusBadRequest, "can not parse message")
	}
	cellString := base64.StdEncoding.EncodeToString(messageCell.ToBOCWithFlags(false))
	return c.JSON(http.StatusOK, echo.Map{
		"transaction": datatype.StorageContractInfo{
			Bin:        cellString,
			ServeBagID: newBagId,
		},
	})
}

// TODO: check if this is needed.
func (h *handler) OptionsHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "Allow: OPTIONS, GET, HEAD, POST")
}
