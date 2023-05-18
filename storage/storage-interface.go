package storage

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tonbyte/provider-node/config"
	"github.com/tonbyte/provider-node/datatype"

	"github.com/labstack/gommon/log"
)

const bagIDString string = "BagID = "
const duplicateHash string = "duplicate hash "

func execStorageCliCommand(command string) ([]byte, error) {
	cmd := exec.Command(config.StorageConfig.SPCliPath,
		`-I`, `127.0.0.1:`+strconv.Itoa(config.StorageConfig.SPCliPort),
		`-k`, config.StorageConfig.StorageDBPath+`/cli-keys/client`,
		`-p`, config.StorageConfig.StorageDBPath+`/cli-keys/server.pub`,
		`-c`, command)

	return cmd.CombinedOutput()
}

func GetProviderInfo() (datatype.ProviderInfo, error) {
	cliQuery := `"get-provider-info" "--json"`
	output, err := execStorageCliCommand(cliQuery)
	if err != nil {
		log.Warn(fmt.Sprintf("ExecQuery() error: %s\noutput: %s", err.Error(), output))
		return datatype.ProviderInfo{}, err
	}

	jsonData := parseJsonOutput(string(output))
	if jsonData == "" {
		return datatype.ProviderInfo{}, fmt.Errorf("invalid json")
	}

	var providerInfo datatype.ProviderInfo
	err = json.Unmarshal([]byte(jsonData), &providerInfo)
	if err != nil {
		return datatype.ProviderInfo{}, err
	}

	return providerInfo, nil
}

func GetProviderParams() (datatype.ProviderParams, error) {
	cliQuery := `"get-provider-params" "--json"`
	output, err := execStorageCliCommand(cliQuery)
	if err != nil {
		log.Warn(fmt.Sprintf("ExecQuery() error: %s\noutput: %s", err.Error(), output))
		return datatype.ProviderParams{}, err
	}

	jsonData := parseJsonOutput(string(output))
	if jsonData == "" {
		return datatype.ProviderParams{}, fmt.Errorf("invalid json")
	}

	var providerParams datatype.ProviderParams
	err = json.Unmarshal([]byte(jsonData), &providerParams)
	if err != nil {
		return datatype.ProviderParams{}, err
	}

	return providerParams, nil
}

func CreateBagOutput(filePath string) string {
	cliQuery := `"create" "--" "` + filePath + `"`
	cmd := exec.Command(config.StorageConfig.SPCliPath,
		`-I`, `127.0.0.1:`+strconv.Itoa(config.StorageConfig.SPCliPort),
		`-k`, config.StorageConfig.StorageDBPath+`/cli-keys/client`,
		`-p`, config.StorageConfig.StorageDBPath+`/cli-keys/server.pub`,
		`-c`, cliQuery)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ExecQuery() error: %e", err)
	}

	return parseCreateBagOutput(string(output))
}

func NewContractMessage(bagID string, filePath string, queryID string, providerAddress string) bool {
	cliQuery := fmt.Sprintf(`"new-contract-message" "%s" "%s" "--query-id" "%s" "--provider" "%s"`, bagID, filePath, queryID, providerAddress)
	output, err := execStorageCliCommand(cliQuery)
	if err != nil {
		log.Warn("cliQuery: " + cliQuery)
		log.Warn(fmt.Sprintf("ExecQuery() error: %s\noutput: %s", err.Error(), output))
	}

	return parseNewMessageOutput(string(output))
}

func RemoveBag(bagID string) bool {
	cliQuery := fmt.Sprintf(`"remove" "%s" "--remove-files"`, bagID)
	output, err := execStorageCliCommand(cliQuery)
	fmt.Println(string(output))

	if err != nil {
		log.Warn(fmt.Sprintf("ExecQuery() error: %s\noutput: %s", err.Error(), output))
	}

	return parseRemoveBagOutput(string(output))
}

func parseCreateBagOutput(output string) string {
	bagIdBegin := strings.Index(output, bagIDString)
	if bagIdBegin != -1 {
		bagIdBegin += len(bagIDString)
	} else {
		bagIdBegin = strings.Index(output, duplicateHash) + len(duplicateHash)
	}

	if bagIdBegin == -1 {
		return ""
	}

	return output[bagIdBegin : bagIdBegin+64]
}

func parseRemoveBagOutput(output string) bool {
	return strings.Contains(output, "No such torrent") || strings.Contains(output, "Success")
}

func parseJsonOutput(output string) string {
	if !strings.Contains(output, "@type") {
		return ""
	}

	startIndex := strings.Index(output, "{")
	if startIndex == -1 {
		return ""
	}

	return output[startIndex:]
}

func parseNewMessageOutput(output string) bool {
	return strings.Contains(output, "Saved message body to file")
}
