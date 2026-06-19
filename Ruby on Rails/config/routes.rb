Rails.application.routes.draw do
  prefix = ENV.fetch("HTTP_BASE", "/api/").delete_suffix("/")

  scope path: prefix do
    get   "urls",     to: "api/urls#index"
    post  "urls",     to: "api/urls#create"
    get   "urls/:id", to: "api/urls#show",  as: :url_info
    get   ":id",      to: "api/urls#redirect", as: :url_redirect
    delete ":id",     to: "api/urls#destroy"
  end
end
