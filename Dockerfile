FROM node:18-alpine

WORKDIR /app

COPY ./swagger/ .

RUN if [ -f package.json ]; then npm install; fi

EXPOSE 8002

CMD ["node", "main.js"]
