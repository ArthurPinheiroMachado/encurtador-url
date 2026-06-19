using System.Text.Json.Serialization;

namespace Encurtador.Url.Models;

public record UrlInfo(
    [property: JsonPropertyName("original")] string Original,
    [property: JsonPropertyName("accesses")] long Accesses
);

public record UrlCreateInput(
    [property: JsonPropertyName("url")] string Url
);

public record UrlCreated(
    [property: JsonPropertyName("id")] string Id,
    [property: JsonPropertyName("url")] string Url
);
