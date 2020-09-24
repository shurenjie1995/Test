package main

import (
	"time"
	"context"
)

const (
	//url
	ApprivalSubscribeUrl     = "/approval/openapi/v2/subscription/subscribe"   //订阅审批事件
	ApprivalUnsubscribeUrl   = "/approval/openapi/v2/subscription/unsubscribe" //取消订阅审批事件
	ApprovalCreateUrl        = "/approval/openapi/v2/instance/create"          //创建审批实例
	ApprovalApproveUrl       = "/approval/openapi/v2/instance/approve"         //同意审批实例
	ApprovalRejectUrl        = "/approval/openapi/v2/instance/reject"          //拒绝审批实例
	ApprovalInstanceGetUrl   = "/approval/openapi/v2/instance/get"             //获取审批实例详情
	ApprovalTransferUrl      = "/approval/openapi/v2/instance/transfer"        //转交审批实例
	ApprovalCancelUrl        = "/approval/openapi/v2/instance/cancel"          //撤销审批实例
	ApprovalUploadUrl        = "/approval/openapi/v2/file/upload"              //上传文件
)


const (
	// 审批结果 approve:通过 reject:拒绝 cancel:取消
	ApprovalResult_Approve ApprovalResult = "APPROVED"
	ApprovalResult_Reject  ApprovalResult = "REJECTED"
	ApprovalResult_Cancel  ApprovalResult = "CANCELED"

	ReqParamErrorCode   = 60001 // 请求参数错误
	NoApprovalCode      = 60002 // 审批定义approval_code找不到
	NoInstanceCode      = 60003 // 审批实例instance_code找不到
	NoUserCode          = 60004 // 用户找不到
	DepVerifyFailCode   = 60005 // 部门验证失败
	FormVerifyFailCode  = 60006 // 表单验证失败
	HadSubscribeCode    = 60007 // 订阅已存在
	NoSubscribeCode     = 60008 // 订阅不存在
	PowerDeniedCode     = 60009 // 权限不足
	NoTaskIdCode        = 60010 // 审批任务task_id找不到
	NeedPayCode         = 60011 // 付费审批，免费用户不可使用
	UuidConflictCode    = 60012 // 审批实例uuid冲突
	DefNonsupportCode   = 60013 // 不支持的审批定义

	TimeOut          = time.Second * 3
)

const(
	VerificationMsg = "url_verification" //机器人url校验
	ApprovalMsg     = "event_callback"   //机器人审批事件回调
)

type ApprovalDef struct { //审批定义
	ApprovalCode string                                                        //审批定义唯一标识
	BotId        string                                                        //审批对应机器人id
	CallBackFunc func(c context.Context, backInfo *CallBackInfo) error //回调函数
}

type ApprovalTask struct { //审批任务
	ApprovalCode string       //审批实例对应审批定义
	InstanceCode string       //审批实例ID
	UserId       string       //审批发起用户
	TaskForm     []ApprovalForm  //审批任务对应的数据，转json
}

type CallBackInfo struct {
	Challenge string         `json:"challenge"`
	Token     string         `json:"token"` //校验Toke
	MsgType   string         `json:"type"`
	Ts        string         `json:"ts"` //时间戳
	Uuid      string         `json:"uuid"`
	Event     interface{}  `json:"event"`
}

type ApprovalEvent struct {
	AppId        string         `json:"app_id"`
	TenantKey    string         `json:"tenant_key"`
	EventType    string         `json:"type"`
	InstanceCode string         `json:"instance_code"` // 审批实例Code
	Status       ApprovalResult `json:"status"`        // 审批结果
	ApprovalCode string         `json:"approval_code"`
}

type OtherEvent struct {
	XxxId        string         `json:"xxx_id"`
	EventType    string         `json:"type"`
	XxxCode string         `json:"xxx_code"` // 审批实例Code
	WwCode string         `json:"ww_code"`
}

type ApprovalResult string



type NotificationTaskStatus int8

const (
	NotificationTaskStatus_Unknown         NotificationTaskStatus = 0
	NotificationTaskStatus_Offline         NotificationTaskStatus = 1
	NotificationTaskStatus_OnlineRevewing  NotificationTaskStatus = 2
	NotificationTaskStatus_Online          NotificationTaskStatus = 3
	NotificationTaskStatus_OfflineRevewing NotificationTaskStatus = 4
	NotificationTaskStatus_RefreshRevewing NotificationTaskStatus = 5
)

type PlatformType int8

const (
	PlatformType_Unknown PlatformType = 0
	PlatformType_Android PlatformType = 1
	PlatformType_iOS     PlatformType = 2
	PlatformType_Windows PlatformType = 3
	PlatformType_Mac     PlatformType = 4
)

type NotificationTaskType int8

const (
	NotificationTaskType_Unknown    NotificationTaskType = 0
	NotificationTaskType_NewFeature NotificationTaskType = 1
)

// TODO change OnlineConfiguration & DraftConfiguration to interface pointer
type NotificationTask struct {
	Id                  int64
	Description         string
	CreateAt            *time.Time
	Status              NotificationTaskStatus
	Type                NotificationTaskType
	MajorVersion        *int8
	MinorVersion        *int8
	Administrators      []string
	OnlineConfiguration interface{}
	DraftConfiguration  interface{}
	ApprovalInstanceId  *string
}

type NewFeatureNotificationTaskConfiguration struct {
	OnlineFrom *time.Time                                       `json:"online_from"`
	Details    []NewFeatureNotificationTaskConfigurationElement `json:"details"`
}

type NewFeatureNotificationTaskConfigurationElement struct {
	TargetPlatforms []PlatformType `json:"target_platforms"`
	CardPayload     CardDetail   `json:"card_payload"`
}

// TODO: move to a separate file, new_feature_approval_form.go
type ApprovalForm struct {
	ID           string                                           `json:"id"`            //任务ID
	Description  string                                           `json:"description"`   //任务描述
	UserId       string                                           `json:"user_id"`       //申请人
	CreateAt     string                                           `json:"create_at"`     //申请时间
	OnOfflineAt  string                                           `json:"onoffline_at"`     //任务上下线时间
	TaskType     NotificationTaskType                             `json:"type"`          //任务申请类型
	ApprovalType NotificationApprovalType                         `json:"approval_type"` //任务审批类型
	Version      string                                           `json:"version"`       //任务对应版本
	Details      []NewFeatureNotificationTaskConfigurationElement `json:"details"`       //素材
}

type CardDetail struct {
	PicI8n         map[string]string `json:"pic_i8n"`
	PicUrlI18n     map[string]string `json:"pic_url_i18n"`
	Title          string            `json:"title"`
	TitleI18n      map[string]string `json:"title_i18n"`
	Content        string            `json:"content"`
	ContentI18n    map[string]string `json:"content_i18n"`
	ButtonContent  string            `json:"button_content"`
	ButtonI18n     map[string]string `json:"button_i18n"`
	ButtonLinkI18n map[string]string `json:"button_link_i18n"`
}

type NotificationApprovalType string

const (
	ApprovalType_Online  NotificationApprovalType = "发起上线" //发起上线
	ApprovalType_Offline NotificationApprovalType = "发起下线" //发起下线
)


func NotificationTaskIsNil(task *NotificationTask) bool {
	if task == nil || task.CreateAt == nil || task.MajorVersion == nil || task.MinorVersion ==nil {
		return true
	}
	return false
}

