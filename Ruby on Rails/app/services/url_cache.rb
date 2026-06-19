require "singleton"

class UrlCache
  include Singleton

  def initialize
    @mutex = Monitor.new
    @urls = {}
  end

  def load
    @mutex.synchronize do
      @urls = Url.pluck(:id, :original, :accesses).each_with_object({}) do |(id, original, accesses), hash|
        hash[id] = { "original" => original, "accesses" => accesses }
      end
    end
  end

  def get_all
    @mutex.synchronize { @urls.deep_dup }
  end

  def get(id)
    @mutex.synchronize { @urls[id]&.dup }
  end

  def exists?(id)
    @mutex.synchronize { @urls.key?(id) }
  end

  def set(id, info)
    @mutex.synchronize { @urls[id] = info }
  end

  def delete(id)
    @mutex.synchronize { @urls.delete(id) }
  end

  def increment_accesses(id)
    @mutex.synchronize do
      if @urls[id]
        @urls[id]["accesses"] += 1
        @urls[id]["accesses"]
      else
        0
      end
    end
  end
end
