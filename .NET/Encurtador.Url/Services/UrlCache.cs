using System.Collections.Concurrent;
using Encurtador.Url.Models;

namespace Encurtador.Url.Services;

public class UrlCache
{
    private readonly ConcurrentDictionary<string, UrlInfo> _urls = new();

    public void Load(IEnumerable<UrlItem> items)
    {
        _urls.Clear();
        foreach (var item in items)
        {
            _urls.TryAdd(item.Id, new UrlInfo(item.Original, item.Accesses));
        }
    }

    public Dictionary<string, UrlInfo> GetAll()
    {
        return _urls.ToDictionary(kvp => kvp.Key, kvp => kvp.Value);
    }

    public UrlInfo? Get(string id)
    {
        return _urls.TryGetValue(id, out var info) ? info : null;
    }

    public bool Exists(string id)
    {
        return _urls.ContainsKey(id);
    }

    public void Set(string id, UrlInfo info)
    {
        _urls[id] = info;
    }

    public void Delete(string id)
    {
        _urls.TryRemove(id, out _);
    }

    public long IncrementAccesses(string id)
    {
        var newInfo = _urls.AddOrUpdate(
            id,
            _ => new UrlInfo("", 1),
            (_, existing) => new UrlInfo(existing.Original, existing.Accesses + 1)
        );
        return newInfo.Accesses;
    }
}
