// Package db provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package db

// Account defines model for Account.
type Account struct {
	AlarmInfo  *AccountAlarmInfo `json:"alarmInfo,omitempty"`
	CreatedOn  *int32            `json:"createdOn,omitempty"`
	Credits    *float64          `json:"credits,omitempty"`
	Id         *string           `json:"id,omitempty"`
	ModifiedOn *int32            `json:"modifiedOn,omitempty"`
	Password   *string           `json:"password,omitempty"`
	Username   string            `json:"username"`
}

// AccountAlarmInfo defines model for AccountAlarmInfo.
type AccountAlarmInfo struct {
	IsCreditLow          *bool `json:"isCreditLow,omitempty"`
	IsMaxNumberOfServers *bool `json:"isMaxNumberOfServers,omitempty"`
}

// Alarm defines model for Alarm.
type Alarm struct {
	AlarmType  string  `json:"alarmType"`
	Id         *string `json:"id,omitempty"`
	Name       string  `json:"name"`
	ResourceId string  `json:"resourceId"`
}

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Server defines model for Server.
type Server struct {
	AccountId          string           `json:"accountId"`
	AlarmInfo          *ServerAlarmInfo `json:"alarmInfo,omitempty"`
	CoreNumber         *string          `json:"coreNumber,omitempty"`
	CreatedOn          *int64           `json:"createdOn,omitempty"`
	Hostname           *string          `json:"hostname,omitempty"`
	Id                 *string          `json:"id,omitempty"`
	License            *int32           `json:"license,omitempty"`
	MemoryAmount       *string          `json:"memoryAmount,omitempty"`
	ModifiedOn         *int64           `json:"modifiedOn,omitempty"`
	Plan               *string          `json:"plan,omitempty"`
	PlanIpV4Bytes      *string          `json:"planIpV4Bytes,omitempty"`
	PlanIpV6Bytes      *string          `json:"planIpV6Bytes,omitempty"`
	ServerCreationTime *int32           `json:"serverCreationTime,omitempty"`
	State              *string          `json:"state,omitempty"`
	Tags               *Tags            `json:"tags,omitempty"`
	Title              *string          `json:"title,omitempty"`
	Uuid               string           `json:"uuid"`
	Zone               *string          `json:"zone,omitempty"`
}

// ServerAlarmInfo defines model for ServerAlarmInfo.
type ServerAlarmInfo struct {
	IsErrored    *bool `json:"isErrored,omitempty"`
	IsPoweredOff *bool `json:"isPoweredOff,omitempty"`
}

// Tags defines model for Tags.
type Tags struct {
	Tag *[]string `json:"tag,omitempty"`
}

// FindAccountsParams defines parameters for FindAccounts.
type FindAccountsParams struct {
	// Username of the account
	Username *string `json:"username,omitempty"`
}

// AddAccountsJSONBody defines parameters for AddAccounts.
type AddAccountsJSONBody Account

// AddAccountByIdJSONBody defines parameters for AddAccountById.
type AddAccountByIdJSONBody Account

// FindAlarmsParams defines parameters for FindAlarms.
type FindAlarmsParams struct {
	// type of the alarm
	AlarmType *string `json:"alarmType,omitempty"`

	// uuid of the associated resource
	ResourceId *string `json:"resourceId,omitempty"`
}

// AddAlarmsJSONBody defines parameters for AddAlarms.
type AddAlarmsJSONBody []interface{}

// AddAlarmJSONBody defines parameters for AddAlarm.
type AddAlarmJSONBody Alarm

// FindServersParams defines parameters for FindServers.
type FindServersParams struct {
	// Id of the accoun the servers are part of
	AccountId *string `json:"accountId,omitempty"`
}

// AddServersJSONBody defines parameters for AddServers.
type AddServersJSONBody []interface{}

// AddServerJSONBody defines parameters for AddServer.
type AddServerJSONBody Server

// AddAccountsJSONRequestBody defines body for AddAccounts for application/json ContentType.
type AddAccountsJSONRequestBody AddAccountsJSONBody

// AddAccountByIdJSONRequestBody defines body for AddAccountById for application/json ContentType.
type AddAccountByIdJSONRequestBody AddAccountByIdJSONBody

// AddAlarmsJSONRequestBody defines body for AddAlarms for application/json ContentType.
type AddAlarmsJSONRequestBody AddAlarmsJSONBody

// AddAlarmJSONRequestBody defines body for AddAlarm for application/json ContentType.
type AddAlarmJSONRequestBody AddAlarmJSONBody

// AddServersJSONRequestBody defines body for AddServers for application/json ContentType.
type AddServersJSONRequestBody AddServersJSONBody

// AddServerJSONRequestBody defines body for AddServer for application/json ContentType.
type AddServerJSONRequestBody AddServerJSONBody