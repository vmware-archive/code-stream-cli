package cmd

// AuthenticationRequest - vRA Authentication request structure
type AuthenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

// AuthenticationResponse - Authentication response structure
type AuthenticationResponse struct {
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// AuthenticationError - Authentication error structure
type AuthenticationError struct {
	Timestamp     int64  `json:"timestamp"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Error         string `json:"error"`
	ServerMessage string `json:"serverMessage"`
}

// ExecutionsList - Code Stream Execution List structure
type ExecutionsList struct {
	Count      int                    `json:"count"`
	TotalCount int                    `json:"totalCount"`
	Links      []string               `json:"links"`
	Documents  map[string]interface{} `json:"documents"`
}

// CodestreamAPIExecutions - Code Stream Execution document structure
type CodestreamAPIExecutions struct {
	Project            string        `json:"project"`
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	UpdatedAt          string        `json:"updatedAt"`
	Link               string        `json:"_link"`
	UpdateTimeInMicros int64         `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64         `json:"_createTimeInMicros"`
	ProjectID          string        `json:"_projectId"`
	Index              int           `json:"index"`
	Notifications      []interface{} `json:"notifications"`
	Comments           string        `json:"comments"`
	Icon               string        `json:"icon"`
	Starred            struct {
	} `json:"starred"`
	Input                 interface{}   `json:"input"`
	Output                interface{}   `json:"output"`
	StageOrder            []interface{} `json:"stageOrder"`
	Stages                interface{}   `json:"stages"`
	Status                string        `json:"status"`
	StatusMessage         string        `json:"statusMessage"`
	DurationInMicros      int           `json:"_durationInMicros"`
	TotalDurationInMicros int           `json:"_totalDurationInMicros"`
	RequestTimeInMicros   int64         `json:"_requestTimeInMicros"`
	ExecutedBy            string        `json:"_executedBy"`
	PipelineLink          string        `json:"_pipelineLink"`
	Nested                bool          `json:"_nested"`
	Rollback              bool          `json:"_rollback"`
	InputMeta             interface{}   `json:"_inputMeta"`
	OutputMeta            interface{}   `json:"_outputMeta"`
	WorkspaceResults      []struct {
		Status string   `json:"status"`
		Step   string   `json:"step"`
		Logs   []string `json:"logs"`
	} `json:"workspaceResults"`
	Tags []interface{} `json:"tags"`
}

// Execution - Code Stream Execution structure
type Execution struct {
	Project            string        `json:"project"`
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	UpdatedAt          string        `json:"updatedAt"`
	Link               string        `json:"_link"`
	UpdateTimeInMicros int64         `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64         `json:"_createTimeInMicros"`
	ProjectID          string        `json:"_projectId"`
	Index              int           `json:"index"`
	Notifications      []interface{} `json:"notifications"`
	Comments           string        `json:"comments"`
	Icon               string        `json:"icon"`
	Starred            struct {
	} `json:"starred"`
	Status                string `json:"status"`
	StatusMessage         string `json:"statusMessage"`
	DurationInMicros      int    `json:"_durationInMicros"`
	TotalDurationInMicros int    `json:"_totalDurationInMicros"`
	RequestTimeInMicros   int64  `json:"_requestTimeInMicros"`
	ExecutedBy            string `json:"_executedBy"`
	PipelineLink          string `json:"_pipelineLink"`
	Nested                bool   `json:"_nested"`
	Rollback              bool   `json:"_rollback"`
	WorkspaceResults      []struct {
		Status string   `json:"status"`
		Step   string   `json:"step"`
		Logs   []string `json:"logs"`
	} `json:"workspaceResults"`
	Tags []interface{} `json:"tags"`
}
