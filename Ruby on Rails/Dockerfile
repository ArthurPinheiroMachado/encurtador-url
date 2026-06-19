FROM ruby:3.2-slim

RUN apt-get update && apt-get install -y --no-install-recommends gcc libpq-dev && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY Gemfile ./

RUN bundle install

COPY . .

EXPOSE 6060

CMD rails db:migrate && rails server -b 0.0.0.0 -p 6060
