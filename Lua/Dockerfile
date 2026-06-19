FROM lua:5.4

RUN apt-get update && apt-get install -y --no-install-recommends luarocks && rm -rf /var/lib/apt/lists/*

RUN luarocks install luasocket && \
    luarocks install pgmoon && \
    luarocks install dkjson

WORKDIR /app
COPY . .

EXPOSE 6060
CMD lua main.lua
