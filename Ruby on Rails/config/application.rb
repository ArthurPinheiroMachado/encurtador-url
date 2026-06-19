require "rails"
require "active_record/railtie"
require "action_controller/railtie"

module Encurtador
  class Application < Rails::Application
    config.load_defaults 7.1
    config.api_only = true

    config.active_record.primary_key_prefix_type = :table_name
  end
end
