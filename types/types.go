package types

type Function struct {
	Name       string        `json:"name" binding:"required"`
	Parameters map[string]interface{} `json:"parameters" binding:"required"`
}

type ObjectExample struct {
	Images   []string      	`json:"images" binding:"required"`
	Function []Function 	`json:"functions" binding:"required"`
}

type JsonExample struct {
	Objects []ObjectExample `json:"objects" binding:"required"`
}