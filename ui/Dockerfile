# Frontend Dockerfile
FROM node:18.16.0-alpine

WORKDIR /app

COPY package.json .

RUN npm install

RUN npm install -g serve

COPY ./build ./build

EXPOSE 3000

CMD ["serve", "-s", "build"]
