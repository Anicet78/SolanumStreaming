FROM node:slim AS dev

WORKDIR /app

COPY ./package*.json ./
RUN npm install

COPY . .

RUN npm run build

CMD ["npm", "start"]

FROM node:slim AS prod

WORKDIR /app

COPY ./package*.json ./
RUN npm install

COPY . .

RUN npm run build

EXPOSE 8084

CMD ["npm", "start"]
