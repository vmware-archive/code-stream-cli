/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
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

// AuthenticationRequestCloud - vRA Authentication request structure for Cloud
type AuthenticationRequestCloud struct {
	RefreshToken string `json:"refreshToken"`
}

// AuthenticationResponseCloud - Authentication response structure for Cloud
type AuthenticationResponseCloud struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
}

// AuthenticationError - Authentication error structure
type AuthenticationError struct {
	Timestamp     int64  `json:"timestamp"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Error         string `json:"error"`
	ServerMessage string `json:"serverMessage"`
}

// documentsList - Code Stream Documents List structure
type documentsList struct {
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

// CodeStreamVariableResponse - Code Stream API Variable response
type CodeStreamVariableResponse struct {
	Project            string `json:"project"`
	Kind               string `json:"kind"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Version            string `json:"version"`
	CreatedBy          string `json:"createdBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Link               string `json:"_link"`
	UpdateTimeInMicros int64  `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64  `json:"_createTimeInMicros"`
	ProjectID          string `json:"_projectId"`
	Type               string `json:"type"`
	Value              string `json:"value"`
}

// CodeStreamVariableRequest - Code Stream API Variable Create Request
type CodeStreamVariableRequest struct {
	Project     string `json:"project"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

// CodeStreamPipeline - Code Stream Pipeline API
type CodeStreamPipeline struct {
	Project            string `json:"project"`
	Kind               string `json:"kind"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	CreatedBy          string `json:"createdBy"`
	UpdatedBy          string `json:"updatedBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Link               string `json:"_link"`
	UpdateTimeInMicros int64  `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64  `json:"_createTimeInMicros"`
	ProjectID          string `json:"_projectId"`
	Icon               string `json:"icon"`
	Enabled            bool   `json:"enabled"`
	Concurrency        int    `json:"concurrency"`
	Input              interface {
	} `json:"input"`
	Output interface {
	} `json:"output"`
	Starred struct {
	} `json:"starred"`
	StageOrder    []string    `json:"stageOrder"`
	Stages        interface{} `json:"stages"`
	Notifications struct {
		Email   []interface{} `json:"email"`
		Jira    []interface{} `json:"jira"`
		Webhook []interface{} `json:"webhook"`
	} `json:"notifications"`
	Options   []interface{} `json:"options"`
	Workspace struct {
		Image    string        `json:"image"`
		Path     string        `json:"path"`
		Endpoint string        `json:"endpoint"`
		Cache    []interface{} `json:"cache"`
		Limits   struct {
			CPU    float64 `json:"cpu"`
			Memory int     `json:"memory"`
		} `json:"limits"`
		AutoCloneForTrigger bool `json:"autoCloneForTrigger"`
	} `json:"workspace"`
	InputMeta  interface{}   `json:"_inputMeta"`
	OutputMeta interface{}   `json:"_outputMeta"`
	Warnings   []interface{} `json:"_warnings"`
	Rollbacks  []interface{} `json:"rollbacks"`
	Tags       []interface{} `json:"tags"`
	State      string        `json:"state"`
}

// CodeStreamCreateExecutionRequest - Code Stream Create Execution Request
type CodeStreamCreateExecutionRequest struct {
	Comments string      `json:"comments"`
	Input    interface{} `json:"input"`
}

// CodeStreamCreateExecutionResponse - Code Stream Create Execution Response
type CodeStreamCreateExecutionResponse struct {
	Comments      string      `json:"comments"`
	Source        string      `json:"source"`
	Input         interface{} `json:"input"`
	ExecutionLink string      `json:"executionLink"`
	Tags          []string    `json:"tags"`
}

// CodeStreamException - Generic exception struct
type CodeStreamException struct {
	Timestamp int64  `json:"timestamp"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
	Type      string `json:"@type"`
}

// CodeStreamEndpoint - Code Stream Create Endpoint
type CodeStreamEndpoint struct {
	Project            string      `json:"project"`
	Kind               string      `json:"kind"`
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	CreatedBy          string      `json:"createdBy"`
	UpdatedBy          string      `json:"updatedBy"`
	CreatedAt          string      `json:"createdAt"`
	UpdatedAt          string      `json:"updatedAt"`
	Link               string      `json:"_link"`
	UpdateTimeInMicros int64       `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64       `json:"_createTimeInMicros"`
	ProjectID          string      `json:"_projectId"`
	Type               string      `json:"type"`
	IsRestricted       bool        `json:"isRestricted"`
	Properties         interface{} `json:"properties"`
	IsLocked           bool        `json:"isLocked"`
	ValidationOutput   string      `json:"validationOutput"`
}

// CodeStreamCustomIntegration - Code Stream Custom Integration
type CodeStreamCustomIntegration struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Version            string `json:"version"`
	CreatedBy          string `json:"createdBy"`
	UpdatedBy          string `json:"updatedBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Link               string `json:"_link"`
	UpdateTimeInMicros int64  `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64  `json:"_createTimeInMicros"`
	Status             string `json:"status"`
	Yaml               string `json:"yaml"`
}
