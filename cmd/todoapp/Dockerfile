FROM golang:1.26.4-bookworm AS builder                                                                   
    WORKDIR /app                                                                                                                                                                      
    COPY go.mod go.sum ./                                                                                    
    RUN --mount=type=cache,target=/go/pkg/mod \                                                              
        go mod download                                                                                            
    COPY . .                                                                                                  
    RUN --mount=type=cache,target=/go/pkg/mod \                                                              
        --mount=type=cache,target=/root/.cache/go-build \                                                    
        CGO_ENABLED=0 go build \                                                                             
        -ldflags="-s -w" \                                                                                   
        -o /app/todoapp \                                                                                    
        ./cmd/todoapp         
                                                                                         
FROM alpine:3.23                                                                                          
    RUN apk --no-cache add ca-certificates tzdata                                                                                                                  
    WORKDIR /app                                                                                               
    COPY --from=builder /app/todoapp /app/todoapp                                                           
    USER 65534:65534                                                                                            
    ENTRYPOINT ["/app/todoapp"]