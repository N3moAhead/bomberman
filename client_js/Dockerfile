FROM docker.io/node:20

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

EXPOSE 8038

CMD ["npm", "start"]
