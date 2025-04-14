// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.1
// source: grpc/tradeapi/v1/orders/orders_service.proto

package orders_service

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	decimal "google.golang.org/genproto/googleapis/type/decimal"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	//side "trade_api/v1/side"
	side "github.com/Ruvad39/go-finam-grpc/trade_api/v1"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Тип заявки
type OrderType int32

const (
	// Значение не указано
	OrderType_ORDER_TYPE_UNSPECIFIED OrderType = 0
	// Рыночная
	OrderType_ORDER_TYPE_MARKET OrderType = 1
	// Лимитная
	OrderType_ORDER_TYPE_LIMIT OrderType = 2
	// Стоп заявка рыночная
	OrderType_ORDER_TYPE_STOP OrderType = 3
	// Стоп заявка лимитная
	OrderType_ORDER_TYPE_STOP_LIMIT OrderType = 4
	// Мульти лег заявка
	OrderType_ORDER_TYPE_MULTI_LEG OrderType = 5
)

// Enum value maps for OrderType.
var (
	OrderType_name = map[int32]string{
		0: "ORDER_TYPE_UNSPECIFIED",
		1: "ORDER_TYPE_MARKET",
		2: "ORDER_TYPE_LIMIT",
		3: "ORDER_TYPE_STOP",
		4: "ORDER_TYPE_STOP_LIMIT",
		5: "ORDER_TYPE_MULTI_LEG",
	}
	OrderType_value = map[string]int32{
		"ORDER_TYPE_UNSPECIFIED": 0,
		"ORDER_TYPE_MARKET":      1,
		"ORDER_TYPE_LIMIT":       2,
		"ORDER_TYPE_STOP":        3,
		"ORDER_TYPE_STOP_LIMIT":  4,
		"ORDER_TYPE_MULTI_LEG":   5,
	}
)

func (x OrderType) Enum() *OrderType {
	p := new(OrderType)
	*p = x
	return p
}

func (x OrderType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderType) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[0].Descriptor()
}

func (OrderType) Type() protoreflect.EnumType {
	return &file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[0]
}

func (x OrderType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderType.Descriptor instead.
func (OrderType) EnumDescriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{0}
}

// Срок действия заявки
type TimeInForce int32

const (
	// Значение не указано
	TimeInForce_TIME_IN_FORCE_UNSPECIFIED TimeInForce = 0
	// До конца дня
	TimeInForce_TIME_IN_FORCE_DAY TimeInForce = 1
	// Good till cancel
	TimeInForce_TIME_IN_FORCE_GOOD_TILL_CANCEL TimeInForce = 2
	// Good till crossing
	TimeInForce_TIME_IN_FORCE_GOOD_TILL_CROSSING TimeInForce = 3
	// Extended hours
	TimeInForce_TIME_IN_FORCE_EXT TimeInForce = 4
	// На открытии биржи
	TimeInForce_TIME_IN_FORCE_ON_OPEN TimeInForce = 5
	// На закрытии биржи
	TimeInForce_TIME_IN_FORCE_ON_CLOSE TimeInForce = 6
	// Immediate or Cancel
	TimeInForce_TIME_IN_FORCE_IOC TimeInForce = 7
	// Fill or Kill
	TimeInForce_TIME_IN_FORCE_FOK TimeInForce = 8
)

// Enum value maps for TimeInForce.
var (
	TimeInForce_name = map[int32]string{
		0: "TIME_IN_FORCE_UNSPECIFIED",
		1: "TIME_IN_FORCE_DAY",
		2: "TIME_IN_FORCE_GOOD_TILL_CANCEL",
		3: "TIME_IN_FORCE_GOOD_TILL_CROSSING",
		4: "TIME_IN_FORCE_EXT",
		5: "TIME_IN_FORCE_ON_OPEN",
		6: "TIME_IN_FORCE_ON_CLOSE",
		7: "TIME_IN_FORCE_IOC",
		8: "TIME_IN_FORCE_FOK",
	}
	TimeInForce_value = map[string]int32{
		"TIME_IN_FORCE_UNSPECIFIED":        0,
		"TIME_IN_FORCE_DAY":                1,
		"TIME_IN_FORCE_GOOD_TILL_CANCEL":   2,
		"TIME_IN_FORCE_GOOD_TILL_CROSSING": 3,
		"TIME_IN_FORCE_EXT":                4,
		"TIME_IN_FORCE_ON_OPEN":            5,
		"TIME_IN_FORCE_ON_CLOSE":           6,
		"TIME_IN_FORCE_IOC":                7,
		"TIME_IN_FORCE_FOK":                8,
	}
)

func (x TimeInForce) Enum() *TimeInForce {
	p := new(TimeInForce)
	*p = x
	return p
}

func (x TimeInForce) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TimeInForce) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[1].Descriptor()
}

func (TimeInForce) Type() protoreflect.EnumType {
	return &file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[1]
}

func (x TimeInForce) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TimeInForce.Descriptor instead.
func (TimeInForce) EnumDescriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{1}
}

// Условие срабатывания стоп заявки
type StopCondition int32

const (
	// Значение не указано
	StopCondition_STOP_CONDITION_UNSPECIFIED StopCondition = 0
	// Цена срабатывания больше текущей цены
	StopCondition_STOP_CONDITION_LAST_UP StopCondition = 1
	// Цена срабатывания меньше текущей цены
	StopCondition_STOP_CONDITION_LAST_DOWN StopCondition = 2
)

// Enum value maps for StopCondition.
var (
	StopCondition_name = map[int32]string{
		0: "STOP_CONDITION_UNSPECIFIED",
		1: "STOP_CONDITION_LAST_UP",
		2: "STOP_CONDITION_LAST_DOWN",
	}
	StopCondition_value = map[string]int32{
		"STOP_CONDITION_UNSPECIFIED": 0,
		"STOP_CONDITION_LAST_UP":     1,
		"STOP_CONDITION_LAST_DOWN":   2,
	}
)

func (x StopCondition) Enum() *StopCondition {
	p := new(StopCondition)
	*p = x
	return p
}

func (x StopCondition) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StopCondition) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[2].Descriptor()
}

func (StopCondition) Type() protoreflect.EnumType {
	return &file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[2]
}

func (x StopCondition) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StopCondition.Descriptor instead.
func (StopCondition) EnumDescriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{2}
}

// Статус заявки
type OrderStatus int32

const (
	// Неопределенное значение
	OrderStatus_ORDER_STATUS_UNSPECIFIED OrderStatus = 0
	// NEW
	OrderStatus_ORDER_STATUS_NEW OrderStatus = 1
	// PARTIALLY_FILLED
	OrderStatus_ORDER_STATUS_PARTIALLY_FILLED OrderStatus = 2
	// FILLED
	OrderStatus_ORDER_STATUS_FILLED OrderStatus = 3
	// DONE_FOR_DAY
	OrderStatus_ORDER_STATUS_DONE_FOR_DAY OrderStatus = 4
	// CANCELED
	OrderStatus_ORDER_STATUS_CANCELED OrderStatus = 5
	// REPLACED
	OrderStatus_ORDER_STATUS_REPLACED OrderStatus = 6
	// PENDING_CANCEL
	OrderStatus_ORDER_STATUS_PENDING_CANCEL OrderStatus = 7
	// REJECTED
	OrderStatus_ORDER_STATUS_REJECTED OrderStatus = 9
	// SUSPENDED
	OrderStatus_ORDER_STATUS_SUSPENDED OrderStatus = 10
	// PENDING_NEW
	OrderStatus_ORDER_STATUS_PENDING_NEW OrderStatus = 11
	// EXPIRED
	OrderStatus_ORDER_STATUS_EXPIRED OrderStatus = 13
	// FAILED
	OrderStatus_ORDER_STATUS_FAILED OrderStatus = 16
	// FORWARDING
	OrderStatus_ORDER_STATUS_FORWARDING OrderStatus = 17
	// WAIT
	OrderStatus_ORDER_STATUS_WAIT OrderStatus = 18
	// DENIED_BY_BROKER
	OrderStatus_ORDER_STATUS_DENIED_BY_BROKER OrderStatus = 19
	// REJECTED_BY_EXCHANGE
	OrderStatus_ORDER_STATUS_REJECTED_BY_EXCHANGE OrderStatus = 20
	// WATCHING
	OrderStatus_ORDER_STATUS_WATCHING OrderStatus = 21
	// EXECUTED
	OrderStatus_ORDER_STATUS_EXECUTED OrderStatus = 22
	// DISABLED
	OrderStatus_ORDER_STATUS_DISABLED OrderStatus = 23
	// LINK_WAIT
	OrderStatus_ORDER_STATUS_LINK_WAIT OrderStatus = 24
	// SL_GUARD_TIME
	OrderStatus_ORDER_STATUS_SL_GUARD_TIME OrderStatus = 27
	// SL_EXECUTED
	OrderStatus_ORDER_STATUS_SL_EXECUTED OrderStatus = 28
	// SL_FORWARDING
	OrderStatus_ORDER_STATUS_SL_FORWARDING OrderStatus = 29
	// TP_GUARD_TIME
	OrderStatus_ORDER_STATUS_TP_GUARD_TIME OrderStatus = 30
	// TP_EXECUTED
	OrderStatus_ORDER_STATUS_TP_EXECUTED OrderStatus = 31
	// TP_CORRECTION
	OrderStatus_ORDER_STATUS_TP_CORRECTION OrderStatus = 32
	// TP_FORWARDING
	OrderStatus_ORDER_STATUS_TP_FORWARDING OrderStatus = 33
	// TP_CORR_GUARD_TIME
	OrderStatus_ORDER_STATUS_TP_CORR_GUARD_TIME OrderStatus = 34
)

// Enum value maps for OrderStatus.
var (
	OrderStatus_name = map[int32]string{
		0:  "ORDER_STATUS_UNSPECIFIED",
		1:  "ORDER_STATUS_NEW",
		2:  "ORDER_STATUS_PARTIALLY_FILLED",
		3:  "ORDER_STATUS_FILLED",
		4:  "ORDER_STATUS_DONE_FOR_DAY",
		5:  "ORDER_STATUS_CANCELED",
		6:  "ORDER_STATUS_REPLACED",
		7:  "ORDER_STATUS_PENDING_CANCEL",
		9:  "ORDER_STATUS_REJECTED",
		10: "ORDER_STATUS_SUSPENDED",
		11: "ORDER_STATUS_PENDING_NEW",
		13: "ORDER_STATUS_EXPIRED",
		16: "ORDER_STATUS_FAILED",
		17: "ORDER_STATUS_FORWARDING",
		18: "ORDER_STATUS_WAIT",
		19: "ORDER_STATUS_DENIED_BY_BROKER",
		20: "ORDER_STATUS_REJECTED_BY_EXCHANGE",
		21: "ORDER_STATUS_WATCHING",
		22: "ORDER_STATUS_EXECUTED",
		23: "ORDER_STATUS_DISABLED",
		24: "ORDER_STATUS_LINK_WAIT",
		27: "ORDER_STATUS_SL_GUARD_TIME",
		28: "ORDER_STATUS_SL_EXECUTED",
		29: "ORDER_STATUS_SL_FORWARDING",
		30: "ORDER_STATUS_TP_GUARD_TIME",
		31: "ORDER_STATUS_TP_EXECUTED",
		32: "ORDER_STATUS_TP_CORRECTION",
		33: "ORDER_STATUS_TP_FORWARDING",
		34: "ORDER_STATUS_TP_CORR_GUARD_TIME",
	}
	OrderStatus_value = map[string]int32{
		"ORDER_STATUS_UNSPECIFIED":          0,
		"ORDER_STATUS_NEW":                  1,
		"ORDER_STATUS_PARTIALLY_FILLED":     2,
		"ORDER_STATUS_FILLED":               3,
		"ORDER_STATUS_DONE_FOR_DAY":         4,
		"ORDER_STATUS_CANCELED":             5,
		"ORDER_STATUS_REPLACED":             6,
		"ORDER_STATUS_PENDING_CANCEL":       7,
		"ORDER_STATUS_REJECTED":             9,
		"ORDER_STATUS_SUSPENDED":            10,
		"ORDER_STATUS_PENDING_NEW":          11,
		"ORDER_STATUS_EXPIRED":              13,
		"ORDER_STATUS_FAILED":               16,
		"ORDER_STATUS_FORWARDING":           17,
		"ORDER_STATUS_WAIT":                 18,
		"ORDER_STATUS_DENIED_BY_BROKER":     19,
		"ORDER_STATUS_REJECTED_BY_EXCHANGE": 20,
		"ORDER_STATUS_WATCHING":             21,
		"ORDER_STATUS_EXECUTED":             22,
		"ORDER_STATUS_DISABLED":             23,
		"ORDER_STATUS_LINK_WAIT":            24,
		"ORDER_STATUS_SL_GUARD_TIME":        27,
		"ORDER_STATUS_SL_EXECUTED":          28,
		"ORDER_STATUS_SL_FORWARDING":        29,
		"ORDER_STATUS_TP_GUARD_TIME":        30,
		"ORDER_STATUS_TP_EXECUTED":          31,
		"ORDER_STATUS_TP_CORRECTION":        32,
		"ORDER_STATUS_TP_FORWARDING":        33,
		"ORDER_STATUS_TP_CORR_GUARD_TIME":   34,
	}
)

func (x OrderStatus) Enum() *OrderStatus {
	p := new(OrderStatus)
	*p = x
	return p
}

func (x OrderStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[3].Descriptor()
}

func (OrderStatus) Type() protoreflect.EnumType {
	return &file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes[3]
}

func (x OrderStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderStatus.Descriptor instead.
func (OrderStatus) EnumDescriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{3}
}

// Запрос на получение конкретного ордера
type GetOrderRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// Идентификатор заявки
	OrderId       string `protobuf:"bytes,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetOrderRequest) Reset() {
	*x = GetOrderRequest{}
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderRequest) ProtoMessage() {}

func (x *GetOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderRequest.ProtoReflect.Descriptor instead.
func (*GetOrderRequest) Descriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetOrderRequest) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *GetOrderRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

// Информация о заявке
type Order struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// Символ инструмента
	Symbol string `protobuf:"bytes,2,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// Количество в шт.
	Quantity *decimal.Decimal `protobuf:"bytes,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	// Сторона (long или short)
	Side side.Side `protobuf:"varint,4,opt,name=side,proto3,enum=grpc.tradeapi.v1.Side" json:"side,omitempty"`
	// Тип заявки
	Type OrderType `protobuf:"varint,5,opt,name=type,proto3,enum=grpc.tradeapi.v1.orders.OrderType" json:"type,omitempty"`
	// Срок действия заявки
	TimeInForce TimeInForce `protobuf:"varint,6,opt,name=time_in_force,json=timeInForce,proto3,enum=grpc.tradeapi.v1.orders.TimeInForce" json:"time_in_force,omitempty"`
	// Необходимо для лимитной и стоп лимитной заявки
	LimitPrice *decimal.Decimal `protobuf:"bytes,7,opt,name=limit_price,json=limitPrice,proto3" json:"limit_price,omitempty"`
	// Необходимо для стоп рыночной и стоп лимитной заявки
	StopPrice *decimal.Decimal `protobuf:"bytes,8,opt,name=stop_price,json=stopPrice,proto3" json:"stop_price,omitempty"`
	// Необходимо для стоп рыночной и стоп лимитной заявки
	StopCondition StopCondition `protobuf:"varint,9,opt,name=stop_condition,json=stopCondition,proto3,enum=grpc.tradeapi.v1.orders.StopCondition" json:"stop_condition,omitempty"`
	// Необходимо для мульти лег заявки
	Legs []*Leg `protobuf:"bytes,10,rep,name=legs,proto3" json:"legs,omitempty"`
	// Уникальный идентификатор заявки. Автоматически генерируется, если не отправлен. (максимум 20 символов)
	ClientOrderId string `protobuf:"bytes,11,opt,name=client_order_id,json=clientOrderId,proto3" json:"client_order_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Order) Reset() {
	*x = Order{}
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{1}
}

func (x *Order) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Order) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Order) GetQuantity() *decimal.Decimal {
	if x != nil {
		return x.Quantity
	}
	return nil
}

func (x *Order) GetSide() side.Side {
	if x != nil {
		return x.Side
	}
	return side.Side(0)
}

func (x *Order) GetType() OrderType {
	if x != nil {
		return x.Type
	}
	return OrderType_ORDER_TYPE_UNSPECIFIED
}

func (x *Order) GetTimeInForce() TimeInForce {
	if x != nil {
		return x.TimeInForce
	}
	return TimeInForce_TIME_IN_FORCE_UNSPECIFIED
}

func (x *Order) GetLimitPrice() *decimal.Decimal {
	if x != nil {
		return x.LimitPrice
	}
	return nil
}

func (x *Order) GetStopPrice() *decimal.Decimal {
	if x != nil {
		return x.StopPrice
	}
	return nil
}

func (x *Order) GetStopCondition() StopCondition {
	if x != nil {
		return x.StopCondition
	}
	return StopCondition_STOP_CONDITION_UNSPECIFIED
}

func (x *Order) GetLegs() []*Leg {
	if x != nil {
		return x.Legs
	}
	return nil
}

func (x *Order) GetClientOrderId() string {
	if x != nil {
		return x.ClientOrderId
	}
	return ""
}

// Лег
type Leg struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Символ инструмента
	Symbol string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// Количество
	Quantity *decimal.Decimal `protobuf:"bytes,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	// Сторона
	Side          side.Side `protobuf:"varint,3,opt,name=side,proto3,enum=grpc.tradeapi.v1.Side" json:"side,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Leg) Reset() {
	*x = Leg{}
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Leg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Leg) ProtoMessage() {}

func (x *Leg) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Leg.ProtoReflect.Descriptor instead.
func (*Leg) Descriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{2}
}

func (x *Leg) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Leg) GetQuantity() *decimal.Decimal {
	if x != nil {
		return x.Quantity
	}
	return nil
}

func (x *Leg) GetSide() side.Side {
	if x != nil {
		return x.Side
	}
	return side.Side(0)
}

// Состояние заявки
type OrderState struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор заявки
	OrderId string `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	// Идентификатор исполнения
	ExecId string `protobuf:"bytes,2,opt,name=exec_id,json=execId,proto3" json:"exec_id,omitempty"`
	// Статус заявки
	Status OrderStatus `protobuf:"varint,3,opt,name=status,proto3,enum=grpc.tradeapi.v1.orders.OrderStatus" json:"status,omitempty"`
	// Заявка
	Order *Order `protobuf:"bytes,4,opt,name=order,proto3" json:"order,omitempty"`
	// Дата и время выставления заявки
	TransactAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=transact_at,json=transactAt,proto3" json:"transact_at,omitempty"`
	// Дата и время принятия заявки
	AcceptAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=accept_at,json=acceptAt,proto3" json:"accept_at,omitempty"`
	// Дата и время  отмены заявки
	WithdrawAt    *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=withdraw_at,json=withdrawAt,proto3" json:"withdraw_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrderState) Reset() {
	*x = OrderState{}
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderState) ProtoMessage() {}

func (x *OrderState) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderState.ProtoReflect.Descriptor instead.
func (*OrderState) Descriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{3}
}

func (x *OrderState) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *OrderState) GetExecId() string {
	if x != nil {
		return x.ExecId
	}
	return ""
}

func (x *OrderState) GetStatus() OrderStatus {
	if x != nil {
		return x.Status
	}
	return OrderStatus_ORDER_STATUS_UNSPECIFIED
}

func (x *OrderState) GetOrder() *Order {
	if x != nil {
		return x.Order
	}
	return nil
}

func (x *OrderState) GetTransactAt() *timestamppb.Timestamp {
	if x != nil {
		return x.TransactAt
	}
	return nil
}

func (x *OrderState) GetAcceptAt() *timestamppb.Timestamp {
	if x != nil {
		return x.AcceptAt
	}
	return nil
}

func (x *OrderState) GetWithdrawAt() *timestamppb.Timestamp {
	if x != nil {
		return x.WithdrawAt
	}
	return nil
}

// Запрос получения списка активных торговых заявок
type OrdersRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId     string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrdersRequest) Reset() {
	*x = OrdersRequest{}
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrdersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrdersRequest) ProtoMessage() {}

func (x *OrdersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrdersRequest.ProtoReflect.Descriptor instead.
func (*OrdersRequest) Descriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{4}
}

func (x *OrdersRequest) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

// Список активных торговых заявок
type OrdersResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Заявки
	Orders        []*OrderState `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *OrdersResponse) Reset() {
	*x = OrdersResponse{}
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrdersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrdersResponse) ProtoMessage() {}

func (x *OrdersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrdersResponse.ProtoReflect.Descriptor instead.
func (*OrdersResponse) Descriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{5}
}

func (x *OrdersResponse) GetOrders() []*OrderState {
	if x != nil {
		return x.Orders
	}
	return nil
}

// Запрос отмены торговой заявки
type CancelOrderRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// Идентификатор заявки
	OrderId       string `protobuf:"bytes,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CancelOrderRequest) Reset() {
	*x = CancelOrderRequest{}
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CancelOrderRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelOrderRequest) ProtoMessage() {}

func (x *CancelOrderRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelOrderRequest.ProtoReflect.Descriptor instead.
func (*CancelOrderRequest) Descriptor() ([]byte, []int) {
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP(), []int{6}
}

func (x *CancelOrderRequest) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *CancelOrderRequest) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

var File_grpc_tradeapi_v1_orders_orders_service_proto protoreflect.FileDescriptor

const file_grpc_tradeapi_v1_orders_orders_service_proto_rawDesc = "" +
	"\n" +
	",grpc/tradeapi/v1/orders/orders_service.proto\x12\x17grpc.tradeapi.v1.orders\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x19google/type/decimal.proto\x1a\x1bgrpc/tradeapi/v1/side.proto\"K\n" +
	"\x0fGetOrderRequest\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\x12\x19\n" +
	"\border_id\x18\x02 \x01(\tR\aorderId\"\xb3\x04\n" +
	"\x05Order\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\x12\x16\n" +
	"\x06symbol\x18\x02 \x01(\tR\x06symbol\x120\n" +
	"\bquantity\x18\x03 \x01(\v2\x14.google.type.DecimalR\bquantity\x12*\n" +
	"\x04side\x18\x04 \x01(\x0e2\x16.grpc.tradeapi.v1.SideR\x04side\x126\n" +
	"\x04type\x18\x05 \x01(\x0e2\".grpc.tradeapi.v1.orders.OrderTypeR\x04type\x12H\n" +
	"\rtime_in_force\x18\x06 \x01(\x0e2$.grpc.tradeapi.v1.orders.TimeInForceR\vtimeInForce\x125\n" +
	"\vlimit_price\x18\a \x01(\v2\x14.google.type.DecimalR\n" +
	"limitPrice\x123\n" +
	"\n" +
	"stop_price\x18\b \x01(\v2\x14.google.type.DecimalR\tstopPrice\x12M\n" +
	"\x0estop_condition\x18\t \x01(\x0e2&.grpc.tradeapi.v1.orders.StopConditionR\rstopCondition\x120\n" +
	"\x04legs\x18\n" +
	" \x03(\v2\x1c.grpc.tradeapi.v1.orders.LegR\x04legs\x12&\n" +
	"\x0fclient_order_id\x18\v \x01(\tR\rclientOrderId\"{\n" +
	"\x03Leg\x12\x16\n" +
	"\x06symbol\x18\x01 \x01(\tR\x06symbol\x120\n" +
	"\bquantity\x18\x02 \x01(\v2\x14.google.type.DecimalR\bquantity\x12*\n" +
	"\x04side\x18\x03 \x01(\x0e2\x16.grpc.tradeapi.v1.SideR\x04side\"\xe7\x02\n" +
	"\n" +
	"OrderState\x12\x19\n" +
	"\border_id\x18\x01 \x01(\tR\aorderId\x12\x17\n" +
	"\aexec_id\x18\x02 \x01(\tR\x06execId\x12<\n" +
	"\x06status\x18\x03 \x01(\x0e2$.grpc.tradeapi.v1.orders.OrderStatusR\x06status\x124\n" +
	"\x05order\x18\x04 \x01(\v2\x1e.grpc.tradeapi.v1.orders.OrderR\x05order\x12;\n" +
	"\vtransact_at\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"transactAt\x127\n" +
	"\taccept_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\bacceptAt\x12;\n" +
	"\vwithdraw_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"withdrawAt\".\n" +
	"\rOrdersRequest\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\"M\n" +
	"\x0eOrdersResponse\x12;\n" +
	"\x06orders\x18\x01 \x03(\v2#.grpc.tradeapi.v1.orders.OrderStateR\x06orders\"N\n" +
	"\x12CancelOrderRequest\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\x12\x19\n" +
	"\border_id\x18\x02 \x01(\tR\aorderId*\x9e\x01\n" +
	"\tOrderType\x12\x1a\n" +
	"\x16ORDER_TYPE_UNSPECIFIED\x10\x00\x12\x15\n" +
	"\x11ORDER_TYPE_MARKET\x10\x01\x12\x14\n" +
	"\x10ORDER_TYPE_LIMIT\x10\x02\x12\x13\n" +
	"\x0fORDER_TYPE_STOP\x10\x03\x12\x19\n" +
	"\x15ORDER_TYPE_STOP_LIMIT\x10\x04\x12\x18\n" +
	"\x14ORDER_TYPE_MULTI_LEG\x10\x05*\x89\x02\n" +
	"\vTimeInForce\x12\x1d\n" +
	"\x19TIME_IN_FORCE_UNSPECIFIED\x10\x00\x12\x15\n" +
	"\x11TIME_IN_FORCE_DAY\x10\x01\x12\"\n" +
	"\x1eTIME_IN_FORCE_GOOD_TILL_CANCEL\x10\x02\x12$\n" +
	" TIME_IN_FORCE_GOOD_TILL_CROSSING\x10\x03\x12\x15\n" +
	"\x11TIME_IN_FORCE_EXT\x10\x04\x12\x19\n" +
	"\x15TIME_IN_FORCE_ON_OPEN\x10\x05\x12\x1a\n" +
	"\x16TIME_IN_FORCE_ON_CLOSE\x10\x06\x12\x15\n" +
	"\x11TIME_IN_FORCE_IOC\x10\a\x12\x15\n" +
	"\x11TIME_IN_FORCE_FOK\x10\b*i\n" +
	"\rStopCondition\x12\x1e\n" +
	"\x1aSTOP_CONDITION_UNSPECIFIED\x10\x00\x12\x1a\n" +
	"\x16STOP_CONDITION_LAST_UP\x10\x01\x12\x1c\n" +
	"\x18STOP_CONDITION_LAST_DOWN\x10\x02*\xe7\x06\n" +
	"\vOrderStatus\x12\x1c\n" +
	"\x18ORDER_STATUS_UNSPECIFIED\x10\x00\x12\x14\n" +
	"\x10ORDER_STATUS_NEW\x10\x01\x12!\n" +
	"\x1dORDER_STATUS_PARTIALLY_FILLED\x10\x02\x12\x17\n" +
	"\x13ORDER_STATUS_FILLED\x10\x03\x12\x1d\n" +
	"\x19ORDER_STATUS_DONE_FOR_DAY\x10\x04\x12\x19\n" +
	"\x15ORDER_STATUS_CANCELED\x10\x05\x12\x19\n" +
	"\x15ORDER_STATUS_REPLACED\x10\x06\x12\x1f\n" +
	"\x1bORDER_STATUS_PENDING_CANCEL\x10\a\x12\x19\n" +
	"\x15ORDER_STATUS_REJECTED\x10\t\x12\x1a\n" +
	"\x16ORDER_STATUS_SUSPENDED\x10\n" +
	"\x12\x1c\n" +
	"\x18ORDER_STATUS_PENDING_NEW\x10\v\x12\x18\n" +
	"\x14ORDER_STATUS_EXPIRED\x10\r\x12\x17\n" +
	"\x13ORDER_STATUS_FAILED\x10\x10\x12\x1b\n" +
	"\x17ORDER_STATUS_FORWARDING\x10\x11\x12\x15\n" +
	"\x11ORDER_STATUS_WAIT\x10\x12\x12!\n" +
	"\x1dORDER_STATUS_DENIED_BY_BROKER\x10\x13\x12%\n" +
	"!ORDER_STATUS_REJECTED_BY_EXCHANGE\x10\x14\x12\x19\n" +
	"\x15ORDER_STATUS_WATCHING\x10\x15\x12\x19\n" +
	"\x15ORDER_STATUS_EXECUTED\x10\x16\x12\x19\n" +
	"\x15ORDER_STATUS_DISABLED\x10\x17\x12\x1a\n" +
	"\x16ORDER_STATUS_LINK_WAIT\x10\x18\x12\x1e\n" +
	"\x1aORDER_STATUS_SL_GUARD_TIME\x10\x1b\x12\x1c\n" +
	"\x18ORDER_STATUS_SL_EXECUTED\x10\x1c\x12\x1e\n" +
	"\x1aORDER_STATUS_SL_FORWARDING\x10\x1d\x12\x1e\n" +
	"\x1aORDER_STATUS_TP_GUARD_TIME\x10\x1e\x12\x1c\n" +
	"\x18ORDER_STATUS_TP_EXECUTED\x10\x1f\x12\x1e\n" +
	"\x1aORDER_STATUS_TP_CORRECTION\x10 \x12\x1e\n" +
	"\x1aORDER_STATUS_TP_FORWARDING\x10!\x12#\n" +
	"\x1fORDER_STATUS_TP_CORR_GUARD_TIME\x10\"2\xc0\x04\n" +
	"\rOrdersService\x12~\n" +
	"\n" +
	"PlaceOrder\x12\x1e.grpc.tradeapi.v1.orders.Order\x1a#.grpc.tradeapi.v1.orders.OrderState\"+\x82\xd3\xe4\x93\x02%:\x01*\" /v1/accounts/{account_id}/orders\x12\x94\x01\n" +
	"\vCancelOrder\x12+.grpc.tradeapi.v1.orders.CancelOrderRequest\x1a#.grpc.tradeapi.v1.orders.OrderState\"3\x82\xd3\xe4\x93\x02-*+/v1/accounts/{account_id}/orders/{order_id}\x12\x86\x01\n" +
	"\tGetOrders\x12&.grpc.tradeapi.v1.orders.OrdersRequest\x1a'.grpc.tradeapi.v1.orders.OrdersResponse\"(\x82\xd3\xe4\x93\x02\"\x12 /v1/accounts/{account_id}/orders\x12\x8e\x01\n" +
	"\bGetOrder\x12(.grpc.tradeapi.v1.orders.GetOrderRequest\x1a#.grpc.tradeapi.v1.orders.OrderState\"3\x82\xd3\xe4\x93\x02-\x12+/v1/accounts/{account_id}/orders/{order_id}B&P\x01Z\"trade_api/v1/orders/orders_serviceb\x06proto3"

var (
	file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescOnce sync.Once
	file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescData []byte
)

func file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescGZIP() []byte {
	file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescOnce.Do(func() {
		file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_grpc_tradeapi_v1_orders_orders_service_proto_rawDesc), len(file_grpc_tradeapi_v1_orders_orders_service_proto_rawDesc)))
	})
	return file_grpc_tradeapi_v1_orders_orders_service_proto_rawDescData
}

var file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_grpc_tradeapi_v1_orders_orders_service_proto_goTypes = []any{
	(OrderType)(0),                // 0: grpc.tradeapi.v1.orders.OrderType
	(TimeInForce)(0),              // 1: grpc.tradeapi.v1.orders.TimeInForce
	(StopCondition)(0),            // 2: grpc.tradeapi.v1.orders.StopCondition
	(OrderStatus)(0),              // 3: grpc.tradeapi.v1.orders.OrderStatus
	(*GetOrderRequest)(nil),       // 4: grpc.tradeapi.v1.orders.GetOrderRequest
	(*Order)(nil),                 // 5: grpc.tradeapi.v1.orders.Order
	(*Leg)(nil),                   // 6: grpc.tradeapi.v1.orders.Leg
	(*OrderState)(nil),            // 7: grpc.tradeapi.v1.orders.OrderState
	(*OrdersRequest)(nil),         // 8: grpc.tradeapi.v1.orders.OrdersRequest
	(*OrdersResponse)(nil),        // 9: grpc.tradeapi.v1.orders.OrdersResponse
	(*CancelOrderRequest)(nil),    // 10: grpc.tradeapi.v1.orders.CancelOrderRequest
	(*decimal.Decimal)(nil),       // 11: google.type.Decimal
	(side.Side)(0),                // 12: grpc.tradeapi.v1.Side
	(*timestamppb.Timestamp)(nil), // 13: google.protobuf.Timestamp
}
var file_grpc_tradeapi_v1_orders_orders_service_proto_depIdxs = []int32{
	11, // 0: grpc.tradeapi.v1.orders.Order.quantity:type_name -> google.type.Decimal
	12, // 1: grpc.tradeapi.v1.orders.Order.side:type_name -> grpc.tradeapi.v1.Side
	0,  // 2: grpc.tradeapi.v1.orders.Order.type:type_name -> grpc.tradeapi.v1.orders.OrderType
	1,  // 3: grpc.tradeapi.v1.orders.Order.time_in_force:type_name -> grpc.tradeapi.v1.orders.TimeInForce
	11, // 4: grpc.tradeapi.v1.orders.Order.limit_price:type_name -> google.type.Decimal
	11, // 5: grpc.tradeapi.v1.orders.Order.stop_price:type_name -> google.type.Decimal
	2,  // 6: grpc.tradeapi.v1.orders.Order.stop_condition:type_name -> grpc.tradeapi.v1.orders.StopCondition
	6,  // 7: grpc.tradeapi.v1.orders.Order.legs:type_name -> grpc.tradeapi.v1.orders.Leg
	11, // 8: grpc.tradeapi.v1.orders.Leg.quantity:type_name -> google.type.Decimal
	12, // 9: grpc.tradeapi.v1.orders.Leg.side:type_name -> grpc.tradeapi.v1.Side
	3,  // 10: grpc.tradeapi.v1.orders.OrderState.status:type_name -> grpc.tradeapi.v1.orders.OrderStatus
	5,  // 11: grpc.tradeapi.v1.orders.OrderState.order:type_name -> grpc.tradeapi.v1.orders.Order
	13, // 12: grpc.tradeapi.v1.orders.OrderState.transact_at:type_name -> google.protobuf.Timestamp
	13, // 13: grpc.tradeapi.v1.orders.OrderState.accept_at:type_name -> google.protobuf.Timestamp
	13, // 14: grpc.tradeapi.v1.orders.OrderState.withdraw_at:type_name -> google.protobuf.Timestamp
	7,  // 15: grpc.tradeapi.v1.orders.OrdersResponse.orders:type_name -> grpc.tradeapi.v1.orders.OrderState
	5,  // 16: grpc.tradeapi.v1.orders.OrdersService.PlaceOrder:input_type -> grpc.tradeapi.v1.orders.Order
	10, // 17: grpc.tradeapi.v1.orders.OrdersService.CancelOrder:input_type -> grpc.tradeapi.v1.orders.CancelOrderRequest
	8,  // 18: grpc.tradeapi.v1.orders.OrdersService.GetOrders:input_type -> grpc.tradeapi.v1.orders.OrdersRequest
	4,  // 19: grpc.tradeapi.v1.orders.OrdersService.GetOrder:input_type -> grpc.tradeapi.v1.orders.GetOrderRequest
	7,  // 20: grpc.tradeapi.v1.orders.OrdersService.PlaceOrder:output_type -> grpc.tradeapi.v1.orders.OrderState
	7,  // 21: grpc.tradeapi.v1.orders.OrdersService.CancelOrder:output_type -> grpc.tradeapi.v1.orders.OrderState
	9,  // 22: grpc.tradeapi.v1.orders.OrdersService.GetOrders:output_type -> grpc.tradeapi.v1.orders.OrdersResponse
	7,  // 23: grpc.tradeapi.v1.orders.OrdersService.GetOrder:output_type -> grpc.tradeapi.v1.orders.OrderState
	20, // [20:24] is the sub-list for method output_type
	16, // [16:20] is the sub-list for method input_type
	16, // [16:16] is the sub-list for extension type_name
	16, // [16:16] is the sub-list for extension extendee
	0,  // [0:16] is the sub-list for field type_name
}

func init() { file_grpc_tradeapi_v1_orders_orders_service_proto_init() }
func file_grpc_tradeapi_v1_orders_orders_service_proto_init() {
	if File_grpc_tradeapi_v1_orders_orders_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_grpc_tradeapi_v1_orders_orders_service_proto_rawDesc), len(file_grpc_tradeapi_v1_orders_orders_service_proto_rawDesc)),
			NumEnums:      4,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_tradeapi_v1_orders_orders_service_proto_goTypes,
		DependencyIndexes: file_grpc_tradeapi_v1_orders_orders_service_proto_depIdxs,
		EnumInfos:         file_grpc_tradeapi_v1_orders_orders_service_proto_enumTypes,
		MessageInfos:      file_grpc_tradeapi_v1_orders_orders_service_proto_msgTypes,
	}.Build()
	File_grpc_tradeapi_v1_orders_orders_service_proto = out.File
	file_grpc_tradeapi_v1_orders_orders_service_proto_goTypes = nil
	file_grpc_tradeapi_v1_orders_orders_service_proto_depIdxs = nil
}
