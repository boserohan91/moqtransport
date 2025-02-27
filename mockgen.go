//go:build gomock || generate

package moqtransport

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_stream_test.go github.com/mengelbart/moqtransport Stream"
type Stream = stream

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_receive_stream_test.go github.com/mengelbart/moqtransport ReceiveStream"
type ReceiveStream = receiveStream

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_send_stream_test.go github.com/mengelbart/moqtransport SendStream"
type SendStream = sendStream

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_connection_test.go github.com/mengelbart/moqtransport Connection"
type Connection = connection

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_parser_test.go github.com/mengelbart/moqtransport Parser"
type Parser = parser

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_parser_factory_test.go github.com/mengelbart/moqtransport ParserFactory"
type ParserFactory = parserFactory

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_sink_test.go github.com/mengelbart/moqtransport Sink"
type Sink = sink

//go:generate sh -c "go run go.uber.org/mock/mockgen -build_flags=\"-tags=gomock\" -package moqtransport -self_package github.com/mengelbart/moqtransport -destination mock_control_msg_sender_test.go github.com/mengelbart/moqtransport ControlMsgSender"
type ControlMsgSender = controlMsgSender
