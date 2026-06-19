module Api
  class UrlsController < ApplicationController
    def index
      render json: UrlCache.instance.get_all
    end

    def create
      url_param = params.require(:url)
      return render json: { detail: "Invalid URL" }, status: :bad_request unless valid_url?(url_param)

      existing = Url.find_by(original: url_param)
      if existing
        return render json: { id: existing.id, url: url_param }, status: :ok
      end

      short_id = ShortIdGenerator.generate(exists: ->(id) { UrlCache.instance.exists?(id) })

      Url.create!(id: short_id, original: url_param, accesses: 0)
      UrlCache.instance.set(short_id, { "original" => url_param, "accesses" => 0 })

      render json: { id: short_id, url: url_param }, status: :created
    end

    def show
      info = UrlCache.instance.get(params[:id])
      return render json: { detail: "URL not found" }, status: :bad_request unless info

      render json: info
    end

    def redirect
      info = UrlCache.instance.get(params[:id])
      return render json: { detail: "URL not found" }, status: :not_found unless info

      new_accesses = UrlCache.instance.increment_accesses(params[:id])
      Url.where(id: params[:id]).update_all(accesses: new_accesses)

      redirect_to info["original"], status: :found
    end

    def destroy
      return render json: { detail: "URL not found" }, status: :bad_request unless UrlCache.instance.exists?(params[:id])

      Url.where(id: params[:id]).delete_all
      UrlCache.instance.delete(params[:id])

      head :ok
    end

    private

    def valid_url?(url)
      uri = URI.parse(url)
      uri.is_a?(URI::HTTP) || uri.is_a?(URI::HTTPS)
    rescue URI::InvalidURIError
      false
    end
  end
end
