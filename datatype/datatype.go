package datatype

type Payload struct {
	ExpirtionTime int64
	Signature     string
}

type FileInfo struct {
	FullPath string `json:"path"`
	Size     string `json:"size"`
}

type Bag struct {
	Type           string `json:"@type"`
	Hash           string `json:"hash"`
	Flags          int    `json:"flags"`
	TotalSize      string `json:"total_size"`
	Description    string `json:"description"`
	FilesCount     string `json:"files_count"`
	IncludedSize   string `json:"included_size"`
	DirName        string `json:"dir_name"`
	DownloadedSize string `json:"downloaded_size"`
	// RootDir		  string `json:"root_dir"`			// No need to parse this
	ActiveDownload bool   `json:"active_download"`
	ActiveUpload   bool   `json:"active_upload"`
	Completed      bool   `json:"completed"`
	DownloadSpeed  string `json:"download_speed"`
	UploadSpeed    string `json:"upload_speed"`
	FatalError     string `json:"fatal_error"`
}

type ProviderInfo struct {
	Type           string `json:"@type"`
	Address        string `json:"address"`
	Balance        string `json:"balance"`
	ProviderConfig struct {
		Type         string `json:"@type"`
		MaxContracts int    `json:"max_contracts"`
		MaxTotalSize string `json:"max_total_size"`
	} `json:"config"`
	ContractsCount     int    `json:"contracts_count"`
	ContractsTotalSize string `json:"contracts_total_size"`
}

type Status struct {
	Version           string `json:"version"`
	RatePerMbDay      string `json:"rate_per_mb_day"`
	Span              int    `json:"span"`
	MinFileSize       string `json:"min_file_size"`
	MaxFileSize       string `json:"max_file_size"`
	ContractsCount    int    `json:"contracts_count"`
	MaxContractsCount int    `json:"max_contracts_count"`
	UsedSpace         string `json:"used_space"`
	MaxTotalSize      string `json:"max_total_size"`
	ContractAddress   string `json:"contract_address"`
	HasGateway        bool   `json:"has_gateway"`
}

type ProviderParams struct {
	Type            string `json:"@type"`
	AcceptNew       bool   `json:"accept_new_contracts"`
	RatePerMbDay    string `json:"rate_per_mb_day"`
	MaxSpan         int    `json:"max_span"`
	MinimalFileSize string `json:"minimal_file_size"`
	MaximalFileSize string `json:"maximal_file_size"`
}

type StorageContractInfo struct {
	Bin        string `json:"bin"`
	ServeBagID string `json:"bag_id"`
}
