FROM node

WORKDIR /app

COPY . .

RUN npm install --global corepack@latest

RUN corepack enable pnpm 

RUN pnpm install

EXPOSE 3000

CMD ["npm", "run", "dev"]
