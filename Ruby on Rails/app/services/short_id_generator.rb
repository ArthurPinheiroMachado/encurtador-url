require "securerandom"

module ShortIdGenerator
  CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
  MAX_ATTEMPTS = 100

  def self.generate(length: 8, exists: nil)
    MAX_ATTEMPTS.times do
      id = length.times.map { CHARSET[SecureRandom.random_number(CHARSET.length)] }.join
      return id if exists.nil? || !exists.call(id)
    end

    raise "failed to generate unique ID after #{MAX_ATTEMPTS} attempts"
  end
end
