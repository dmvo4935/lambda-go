FROM golang:1.21 as build
WORKDIR /function
# Copy dependencies list
COPY go.mod go.sum ./
# Build with optional lambda.norpc tag
COPY main.go .
ENV GOOS=linux GOARCH=amd64
RUN go build -tags lambda.norpc -o main main.go

# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=build /function/main ./main
ENTRYPOINT [ "./main" ]
