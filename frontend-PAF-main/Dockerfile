FROM node:lts-alpine

RUN npm install -g http-server

RUN mkdir -p /app

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

EXPOSE 3001

CMD [ "http-server", "dist" ]
