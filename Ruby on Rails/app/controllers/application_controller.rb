class ApplicationController < ActionController::API
  before_action :authenticate

  private

  def authenticate
    authenticate_or_request_with_http_basic do |username, password|
      expected_user = ENV.fetch("USER", "user")
      expected_pass = ENV.fetch("PASS", "pass123")
      ActiveSupport::SecurityUtils.secure_compare(username, expected_user) &
        ActiveSupport::SecurityUtils.secure_compare(password, expected_pass)
    end
  end
end
