FROM node:22-alpine AS dev
WORKDIR /app
RUN corepack enable   # enables yarn bundled with Node >=16
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY frontend/ .
EXPOSE 5173
CMD ["yarn", "dev", "--host", "0.0.0.0"]

FROM node:20-alpine AS build
WORKDIR /app
RUN corepack enable
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY frontend/ .
RUN yarn build

FROM nginx:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80

