FROM node

WORKDIR /app

RUN npm install --global corepack@latest
RUN corepack enable pnpm 

COPY package.json pnpm-lock.yaml ./
RUN pnpm install
COPY . .

EXPOSE 3000

CMD ["npm", "run", "dev"]
