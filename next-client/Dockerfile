FROM node:20-alpine

WORKDIR /app

RUN apk update

COPY . .

RUN mv ./.env.docker ./.env
RUN yarn install
RUN yarn build

CMD yarn start
 