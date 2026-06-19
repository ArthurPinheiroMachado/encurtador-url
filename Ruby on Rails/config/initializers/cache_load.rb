Rails.application.config.after_initialize do
  UrlCache.instance.load
rescue => e
  Rails.logger.warn "UrlCache could not be loaded: #{e.message}"
end
