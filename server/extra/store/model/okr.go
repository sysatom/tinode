package model

type Objective struct {
	Id           int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"  gorm:"primaryKey"`
	UserId       int64  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" `
	Sequence     int64  `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty" `
	Title        string `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty" `
	Memo         string `protobuf:"bytes,5,opt,name=memo,proto3" json:"memo,omitempty" `
	Motive       string `protobuf:"bytes,6,opt,name=motive,proto3" json:"motive,omitempty" `
	Feasibility  string `protobuf:"bytes,7,opt,name=feasibility,proto3" json:"feasibility,omitempty" `
	IsPlan       bool   `protobuf:"varint,8,opt,name=is_plan,json=isPlan,proto3" json:"is_plan,omitempty" `
	PlanStart    int64  `protobuf:"varint,9,opt,name=plan_start,json=planStart,proto3" json:"plan_start,omitempty" `
	PlanEnd      int64  `protobuf:"varint,10,opt,name=plan_end,json=planEnd,proto3" json:"plan_end,omitempty" `
	TotalValue   int32  `protobuf:"varint,11,opt,name=total_value,json=totalValue,proto3" json:"total_value,omitempty" `
	CurrentValue int32  `protobuf:"varint,12,opt,name=current_value,json=currentValue,proto3" json:"current_value,omitempty" `
	CreatedAt    int64  `protobuf:"varint,13,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" `
	UpdatedAt    int64  `protobuf:"varint,14,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" `
	Tag          string `protobuf:"bytes,15,opt,name=tag,proto3" json:"tag,omitempty" gorm:"-"`
}

func (Objective) TableName() string {
	return "chatbot_objectives"
}

type KeyResult struct {
	Id           int64         `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"  gorm:"primaryKey"`
	UserId       int64         `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty" `
	ObjectiveId  int64         `protobuf:"varint,3,opt,name=objective_id,json=objectiveId,proto3" json:"objective_id,omitempty" `
	Sequence     int64         `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty" `
	Title        string        `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty" `
	Memo         string        `protobuf:"bytes,6,opt,name=memo,proto3" json:"memo,omitempty" `
	InitialValue int32         `protobuf:"varint,7,opt,name=initial_value,json=initialValue,proto3" json:"initial_value,omitempty" `
	TargetValue  int32         `protobuf:"varint,8,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty" `
	CurrentValue int32         `protobuf:"varint,9,opt,name=current_value,json=currentValue,proto3" json:"current_value,omitempty" `
	ValueMode    ValueModeType `protobuf:"bytes,10,opt,name=value_mode,json=valueMode,proto3" json:"value_mode,omitempty" `
	CreatedAt    int64         `protobuf:"varint,11,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" `
	UpdatedAt    int64         `protobuf:"varint,12,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" `
	Tag          string        `protobuf:"bytes,13,opt,name=tag,proto3" json:"tag,omitempty" gorm:"-"`
}

func (KeyResult) TableName() string {
	return "chatbot_key_results"
}

type ValueModeType string

const (
	ValueSumMode  ValueModeType = "sum"
	ValueLastMode ValueModeType = "last"
	ValueAvgMode  ValueModeType = "avg"
	ValueMaxMode  ValueModeType = "max"
)

type KeyResultValue struct {
	Id          int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"  gorm:"primaryKey"`
	KeyResultId int64 `protobuf:"varint,2,opt,name=key_result_id,json=keyResultId,proto3" json:"key_result_id,omitempty" `
	Value       int32 `protobuf:"varint,3,opt,name=value,proto3" json:"value,omitempty" `
	CreatedAt   int64 `protobuf:"varint,13,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" `
	UpdatedAt   int64 `protobuf:"varint,14,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" `
}

func (KeyResultValue) TableName() string {
	return "chatbot_key_result_values"
}
