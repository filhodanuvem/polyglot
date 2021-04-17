from golang:1.16-alpine as build
workdir /project 
copy . .
run go build -o polyglot . 

from alpine:latest
copy --from=build /project/polyglot .
cmd ./polyglot -s
expose 8080
