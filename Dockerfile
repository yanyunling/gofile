FROM node:24-alpine AS nodejs
COPY ./web /build/web
WORKDIR /build/web
# 淘宝npm镜像
# RUN npm config set registry https://registry.npmmirror.com
RUN npm install
RUN npm run build

FROM golang:1.26-alpine AS go
# 七牛云镜像
# ENV GOPROXY=https://goproxy.cn,direct
COPY ./ /build/gofile
COPY --from=nodejs /build/web/dist /build/gofile/web/dist
WORKDIR /build/gofile
RUN go build

FROM alpine:latest
COPY --from=go /build/gofile/gofile /gofile/
ENV pass=Admin123
ENV limit=0
EXPOSE 9300
RUN chmod +x /gofile/gofile
CMD /gofile/gofile -p 9300 -log /gofile/logs -data /gofile/data -pass=${pass} -limit=${limit}
