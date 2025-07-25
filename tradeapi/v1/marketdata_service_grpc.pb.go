// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.1
// source: grpc/tradeapi/v1/marketdata/marketdata_service.proto

//package marketdata_service
package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	MarketDataService_Bars_FullMethodName                  = "/grpc.tradeapi.v1.marketdata.MarketDataService/Bars"
	MarketDataService_LastQuote_FullMethodName             = "/grpc.tradeapi.v1.marketdata.MarketDataService/LastQuote"
	MarketDataService_OrderBook_FullMethodName             = "/grpc.tradeapi.v1.marketdata.MarketDataService/OrderBook"
	MarketDataService_LatestTrades_FullMethodName          = "/grpc.tradeapi.v1.marketdata.MarketDataService/LatestTrades"
	MarketDataService_SubscribeQuote_FullMethodName        = "/grpc.tradeapi.v1.marketdata.MarketDataService/SubscribeQuote"
	MarketDataService_SubscribeOrderBook_FullMethodName    = "/grpc.tradeapi.v1.marketdata.MarketDataService/SubscribeOrderBook"
	MarketDataService_SubscribeLatestTrades_FullMethodName = "/grpc.tradeapi.v1.marketdata.MarketDataService/SubscribeLatestTrades"
	MarketDataService_SubscribeBars_FullMethodName         = "/grpc.tradeapi.v1.marketdata.MarketDataService/SubscribeBars"
)

// MarketDataServiceClient is the client API for MarketDataService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Сервис рыночных данных
type MarketDataServiceClient interface {
	// Получение исторических данных по инструменту (агрегированные свечи)
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/bars?timeframe=TIME_FRAME_D&interval.start_time=2023-01-01T00:00:00Z&interval.end_time=2023-01-31T23:59:59Z
	// Authorization: <token>
	//
	// Параметры:
	// - symbol - передается в URL пути
	// - timeframe и interval - передаются как query-параметры
	Bars(ctx context.Context, in *BarsRequest, opts ...grpc.CallOption) (*BarsResponse, error)
	// Получение последней котировки по инструменту
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/quotes/latest
	// Authorization: <token>
	LastQuote(ctx context.Context, in *QuoteRequest, opts ...grpc.CallOption) (*QuoteResponse, error)
	// Получение текущего стакана по инструменту
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/orderbook
	// Authorization: <token>
	OrderBook(ctx context.Context, in *OrderBookRequest, opts ...grpc.CallOption) (*OrderBookResponse, error)
	// Получение списка последних сделок по инструменту
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/trades/latest
	// Authorization: <token>
	LatestTrades(ctx context.Context, in *LatestTradesRequest, opts ...grpc.CallOption) (*LatestTradesResponse, error)
	// Подписка на котировки по инструменту. Стрим метод
	SubscribeQuote(ctx context.Context, in *SubscribeQuoteRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeQuoteResponse], error)
	// Подписка на стакан по инструменту. Стрим метод
	SubscribeOrderBook(ctx context.Context, in *SubscribeOrderBookRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeOrderBookResponse], error)
	// Подписка на сделки по инструменту. Стрим метод
	SubscribeLatestTrades(ctx context.Context, in *SubscribeLatestTradesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeLatestTradesResponse], error)
	// Подписка на агрегированные свечи. Стрим метод
	SubscribeBars(ctx context.Context, in *SubscribeBarsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeBarsResponse], error)
}

type marketDataServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMarketDataServiceClient(cc grpc.ClientConnInterface) MarketDataServiceClient {
	return &marketDataServiceClient{cc}
}

func (c *marketDataServiceClient) Bars(ctx context.Context, in *BarsRequest, opts ...grpc.CallOption) (*BarsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BarsResponse)
	err := c.cc.Invoke(ctx, MarketDataService_Bars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketDataServiceClient) LastQuote(ctx context.Context, in *QuoteRequest, opts ...grpc.CallOption) (*QuoteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuoteResponse)
	err := c.cc.Invoke(ctx, MarketDataService_LastQuote_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketDataServiceClient) OrderBook(ctx context.Context, in *OrderBookRequest, opts ...grpc.CallOption) (*OrderBookResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OrderBookResponse)
	err := c.cc.Invoke(ctx, MarketDataService_OrderBook_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketDataServiceClient) LatestTrades(ctx context.Context, in *LatestTradesRequest, opts ...grpc.CallOption) (*LatestTradesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LatestTradesResponse)
	err := c.cc.Invoke(ctx, MarketDataService_LatestTrades_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketDataServiceClient) SubscribeQuote(ctx context.Context, in *SubscribeQuoteRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeQuoteResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MarketDataService_ServiceDesc.Streams[0], MarketDataService_SubscribeQuote_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeQuoteRequest, SubscribeQuoteResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeQuoteClient = grpc.ServerStreamingClient[SubscribeQuoteResponse]

func (c *marketDataServiceClient) SubscribeOrderBook(ctx context.Context, in *SubscribeOrderBookRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeOrderBookResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MarketDataService_ServiceDesc.Streams[1], MarketDataService_SubscribeOrderBook_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeOrderBookRequest, SubscribeOrderBookResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeOrderBookClient = grpc.ServerStreamingClient[SubscribeOrderBookResponse]

func (c *marketDataServiceClient) SubscribeLatestTrades(ctx context.Context, in *SubscribeLatestTradesRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeLatestTradesResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MarketDataService_ServiceDesc.Streams[2], MarketDataService_SubscribeLatestTrades_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeLatestTradesRequest, SubscribeLatestTradesResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeLatestTradesClient = grpc.ServerStreamingClient[SubscribeLatestTradesResponse]

func (c *marketDataServiceClient) SubscribeBars(ctx context.Context, in *SubscribeBarsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[SubscribeBarsResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &MarketDataService_ServiceDesc.Streams[3], MarketDataService_SubscribeBars_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SubscribeBarsRequest, SubscribeBarsResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeBarsClient = grpc.ServerStreamingClient[SubscribeBarsResponse]

// MarketDataServiceServer is the server API for MarketDataService service.
// All implementations must embed UnimplementedMarketDataServiceServer
// for forward compatibility.
//
// Сервис рыночных данных
type MarketDataServiceServer interface {
	// Получение исторических данных по инструменту (агрегированные свечи)
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/bars?timeframe=TIME_FRAME_D&interval.start_time=2023-01-01T00:00:00Z&interval.end_time=2023-01-31T23:59:59Z
	// Authorization: <token>
	//
	// Параметры:
	// - symbol - передается в URL пути
	// - timeframe и interval - передаются как query-параметры
	Bars(context.Context, *BarsRequest) (*BarsResponse, error)
	// Получение последней котировки по инструменту
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/quotes/latest
	// Authorization: <token>
	LastQuote(context.Context, *QuoteRequest) (*QuoteResponse, error)
	// Получение текущего стакана по инструменту
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/orderbook
	// Authorization: <token>
	OrderBook(context.Context, *OrderBookRequest) (*OrderBookResponse, error)
	// Получение списка последних сделок по инструменту
	// Пример HTTP запроса:
	// GET /v1/instruments/SBER@MISX/trades/latest
	// Authorization: <token>
	LatestTrades(context.Context, *LatestTradesRequest) (*LatestTradesResponse, error)
	// Подписка на котировки по инструменту. Стрим метод
	SubscribeQuote(*SubscribeQuoteRequest, grpc.ServerStreamingServer[SubscribeQuoteResponse]) error
	// Подписка на стакан по инструменту. Стрим метод
	SubscribeOrderBook(*SubscribeOrderBookRequest, grpc.ServerStreamingServer[SubscribeOrderBookResponse]) error
	// Подписка на сделки по инструменту. Стрим метод
	SubscribeLatestTrades(*SubscribeLatestTradesRequest, grpc.ServerStreamingServer[SubscribeLatestTradesResponse]) error
	// Подписка на агрегированные свечи. Стрим метод
	SubscribeBars(*SubscribeBarsRequest, grpc.ServerStreamingServer[SubscribeBarsResponse]) error
	mustEmbedUnimplementedMarketDataServiceServer()
}

// UnimplementedMarketDataServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMarketDataServiceServer struct{}

func (UnimplementedMarketDataServiceServer) Bars(context.Context, *BarsRequest) (*BarsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bars not implemented")
}
func (UnimplementedMarketDataServiceServer) LastQuote(context.Context, *QuoteRequest) (*QuoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LastQuote not implemented")
}
func (UnimplementedMarketDataServiceServer) OrderBook(context.Context, *OrderBookRequest) (*OrderBookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OrderBook not implemented")
}
func (UnimplementedMarketDataServiceServer) LatestTrades(context.Context, *LatestTradesRequest) (*LatestTradesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LatestTrades not implemented")
}
func (UnimplementedMarketDataServiceServer) SubscribeQuote(*SubscribeQuoteRequest, grpc.ServerStreamingServer[SubscribeQuoteResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeQuote not implemented")
}
func (UnimplementedMarketDataServiceServer) SubscribeOrderBook(*SubscribeOrderBookRequest, grpc.ServerStreamingServer[SubscribeOrderBookResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeOrderBook not implemented")
}
func (UnimplementedMarketDataServiceServer) SubscribeLatestTrades(*SubscribeLatestTradesRequest, grpc.ServerStreamingServer[SubscribeLatestTradesResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeLatestTrades not implemented")
}
func (UnimplementedMarketDataServiceServer) SubscribeBars(*SubscribeBarsRequest, grpc.ServerStreamingServer[SubscribeBarsResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeBars not implemented")
}
func (UnimplementedMarketDataServiceServer) mustEmbedUnimplementedMarketDataServiceServer() {}
func (UnimplementedMarketDataServiceServer) testEmbeddedByValue()                           {}

// UnsafeMarketDataServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MarketDataServiceServer will
// result in compilation errors.
type UnsafeMarketDataServiceServer interface {
	mustEmbedUnimplementedMarketDataServiceServer()
}

func RegisterMarketDataServiceServer(s grpc.ServiceRegistrar, srv MarketDataServiceServer) {
	// If the following call pancis, it indicates UnimplementedMarketDataServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MarketDataService_ServiceDesc, srv)
}

func _MarketDataService_Bars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BarsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketDataServiceServer).Bars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MarketDataService_Bars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketDataServiceServer).Bars(ctx, req.(*BarsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketDataService_LastQuote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketDataServiceServer).LastQuote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MarketDataService_LastQuote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketDataServiceServer).LastQuote(ctx, req.(*QuoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketDataService_OrderBook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderBookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketDataServiceServer).OrderBook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MarketDataService_OrderBook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketDataServiceServer).OrderBook(ctx, req.(*OrderBookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketDataService_LatestTrades_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LatestTradesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketDataServiceServer).LatestTrades(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MarketDataService_LatestTrades_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketDataServiceServer).LatestTrades(ctx, req.(*LatestTradesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketDataService_SubscribeQuote_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeQuoteRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MarketDataServiceServer).SubscribeQuote(m, &grpc.GenericServerStream[SubscribeQuoteRequest, SubscribeQuoteResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeQuoteServer = grpc.ServerStreamingServer[SubscribeQuoteResponse]

func _MarketDataService_SubscribeOrderBook_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeOrderBookRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MarketDataServiceServer).SubscribeOrderBook(m, &grpc.GenericServerStream[SubscribeOrderBookRequest, SubscribeOrderBookResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeOrderBookServer = grpc.ServerStreamingServer[SubscribeOrderBookResponse]

func _MarketDataService_SubscribeLatestTrades_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeLatestTradesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MarketDataServiceServer).SubscribeLatestTrades(m, &grpc.GenericServerStream[SubscribeLatestTradesRequest, SubscribeLatestTradesResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeLatestTradesServer = grpc.ServerStreamingServer[SubscribeLatestTradesResponse]

func _MarketDataService_SubscribeBars_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeBarsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MarketDataServiceServer).SubscribeBars(m, &grpc.GenericServerStream[SubscribeBarsRequest, SubscribeBarsResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type MarketDataService_SubscribeBarsServer = grpc.ServerStreamingServer[SubscribeBarsResponse]

// MarketDataService_ServiceDesc is the grpc.ServiceDesc for MarketDataService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MarketDataService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.tradeapi.v1.marketdata.MarketDataService",
	HandlerType: (*MarketDataServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bars",
			Handler:    _MarketDataService_Bars_Handler,
		},
		{
			MethodName: "LastQuote",
			Handler:    _MarketDataService_LastQuote_Handler,
		},
		{
			MethodName: "OrderBook",
			Handler:    _MarketDataService_OrderBook_Handler,
		},
		{
			MethodName: "LatestTrades",
			Handler:    _MarketDataService_LatestTrades_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeQuote",
			Handler:       _MarketDataService_SubscribeQuote_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeOrderBook",
			Handler:       _MarketDataService_SubscribeOrderBook_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeLatestTrades",
			Handler:       _MarketDataService_SubscribeLatestTrades_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeBars",
			Handler:       _MarketDataService_SubscribeBars_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc/tradeapi/v1/marketdata/marketdata_service.proto",
}
