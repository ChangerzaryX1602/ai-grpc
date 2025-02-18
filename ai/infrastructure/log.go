package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ANSI color codes
const (
	colorReset   = "\033[0m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorRed     = "\033[31m"
	colorCyan    = "\033[36m"
	colorMagenta = "\033[35m"
)

// LogResponsesInterceptor logs gRPC requests and responses with color output.
func LogResponsesInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	fmt.Println(colorCyan + "===== Incoming gRPC Request =====" + colorReset)
	fmt.Printf(colorBlue+"Method: %s\n"+colorReset, info.FullMethod)

	// Log metadata if present
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Println(colorMagenta + "Metadata:" + colorReset)
		for key, value := range md {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	// Log the request payload
	fmt.Println(colorYellow + "Request Payload:" + colorReset)
	logJSON(req)

	// Call the handler to get the response or error
	res, err := handler(ctx, req)

	// Log the response or error
	if err != nil {
		grpcErr, _ := status.FromError(err)
		fmt.Println(colorRed + "===== gRPC Response Error =====" + colorReset)
		fmt.Printf("Error Code: %s\n", grpcErr.Code())
		fmt.Printf("Error Message: %s\n", grpcErr.Message())
	} else {
		fmt.Println(colorGreen + "===== gRPC Response =====" + colorReset)
		logJSON(res)
	}

	fmt.Println(colorCyan + "===============================" + colorReset)
	return res, err
}

// logJSON pretty-prints any object as JSON
func logJSON(v interface{}) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(colorRed+"Error formatting JSON:"+colorReset, err)
		return
	}
	fmt.Println(string(jsonData))
}
