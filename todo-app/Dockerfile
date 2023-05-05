FROM node:20-bullseye-slim

WORKDIR /app

# install dependencies
COPY package.json package-lock.json .
RUN --mount=type=cache,target=/root/.npm \
  npm ci

# copy app source
COPY . .

CMD ["npm", "run", "dev", "--", "--host=0.0.0.0"]