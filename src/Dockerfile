 
FROM golang:latest as build
  
# Move to working directory /build
WORKDIR /src
ADD . /src
# Copy and download dependency using go mod
RUN go get -d -v ./...

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
RUN go build -o /dist main.go
 
FROM gcr.io/distroless/base
COPY --from=build /dist /
ENTRYPOINT [ "/dist" ]