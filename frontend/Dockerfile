FROM node:lts-alpine AS build

WORKDIR /app

COPY package*.json .

RUN npm install 

COPY . .

RUN npm run build

FROM nginx

COPY --from=build /app/dist/ /usr/share/nginx/html/

COPY nginx.conf /etc/nginx/conf.d/default.conf