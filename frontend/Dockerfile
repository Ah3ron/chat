FROM node:23 AS builder

WORKDIR /app

RUN npm install -g pnpm

COPY package.json pnpm-lock.yaml ./

ENV NODE_ENV=development
RUN  pnpm install --frozen-lockfile
RUN pnpm add vite --save-dev

COPY . .

RUN pnpm run build

# --- production image ---
# --- production image ---
FROM nginx:alpine

# Указываем рабочую директорию для nginx
WORKDIR /usr/share/nginx/html

# Удаляем стандартные статику по умолчанию
RUN rm -rf ./*

# Копируем собранный билд из предыдущего этапа
# Предполагается, что на предыдущем этапе сборки (builder)
# файлы собраны и находятся в директории /app/build
COPY --from=builder /app/build/ ./

# Если требуется использовать кастомный конфиг nginx, раскомментируйте и настройте строку ниже:
# COPY nginx.conf /etc/nginx/nginx.conf

# Открываем порт 80
EXPOSE 80

# Запуск nginx в foreground
ENTRYPOINT ["nginx", "-g", "daemon off;"]
