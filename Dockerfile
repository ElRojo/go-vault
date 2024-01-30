FROM golang AS build
WORKDIR /src
ADD ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o go-vault

FROM alpine
WORKDIR /src
COPY --from=build /src/go-vault /src/
CMD "./go-vault"



