using Encurtador.Url.Auth;
using Encurtador.Url.Data;
using Encurtador.Url.Models;
using Encurtador.Url.Services;
using Microsoft.AspNetCore.Authentication;
using Microsoft.EntityFrameworkCore;

var builder = WebApplication.CreateBuilder(args);

var dbHost = Environment.GetEnvironmentVariable("DB_HOST") ?? "0.0.0.0";
var dbPort = Environment.GetEnvironmentVariable("DB_PORT") ?? "5432";
var dbName = Environment.GetEnvironmentVariable("DB_NAME") ?? "encurtador";
var dbUser = Environment.GetEnvironmentVariable("DB_USER") ?? "postgres";
var dbPass = Environment.GetEnvironmentVariable("DB_PASS") ?? "postgres";
var httpBase = Environment.GetEnvironmentVariable("HTTP_BASE") ?? "/api/";
var httpPort = Environment.GetEnvironmentVariable("HTTP_PORT") ?? "6060";

var connectionString = $"Host={dbHost};Port={dbPort};Database={dbName};Username={dbUser};Password={dbPass}";

builder.WebHost.UseUrls($"http://0.0.0.0:{httpPort}");

builder.Services.AddDbContext<AppDbContext>(options =>
    options.UseNpgsql(connectionString));

builder.Services.AddAuthentication(BasicAuthenticationHandler.SchemeName)
    .AddScheme<BasicAuthenticationOptions, BasicAuthenticationHandler>(
        BasicAuthenticationHandler.SchemeName, null);

builder.Services.AddAuthorization();
builder.Services.AddSingleton<UrlCache>();

var app = builder.Build();

app.UseAuthentication();
app.UseAuthorization();

using (var scope = app.Services.CreateScope())
{
    var db = scope.ServiceProvider.GetRequiredService<AppDbContext>();
    db.Database.EnsureCreated();

    var cache = scope.ServiceProvider.GetRequiredService<UrlCache>();
    cache.Load(db.Urls.AsEnumerable());
}

var prefix = httpBase.TrimEnd('/');

app.MapGet($"{prefix}/urls", async (UrlCache cache) =>
{
    return Results.Ok(cache.GetAll());
}).RequireAuthorization();

app.MapPost($"{prefix}/urls", async (UrlCreateInput input, AppDbContext db, UrlCache cache) =>
{
    if (!Uri.TryCreate(input.Url, UriKind.Absolute, out var uri) ||
        (uri.Scheme != "http" && uri.Scheme != "https"))
    {
        return Results.BadRequest(new { detail = "Invalid URL" });
    }

    var existing = await db.Urls.FirstOrDefaultAsync(u => u.Original == input.Url);
    if (existing != null)
    {
        return Results.Ok(new UrlCreated(existing.Id, input.Url));
    }

    var shortId = ShortIdGenerator.Generate(8, id => cache.Exists(id));

    db.Urls.Add(new UrlItem { Id = shortId, Original = input.Url, Accesses = 0 });
    await db.SaveChangesAsync();

    cache.Set(shortId, new UrlInfo(input.Url, 0));

    return Results.Json(new UrlCreated(shortId, input.Url), statusCode: 201);
}).RequireAuthorization();

app.MapGet($"{prefix}/urls/{{id}}", (string id, UrlCache cache) =>
{
    var info = cache.Get(id);
    return info is null
        ? Results.BadRequest(new { detail = "URL not found" })
        : Results.Ok(info);
}).RequireAuthorization();

app.MapGet($"{prefix}/{{id}}", async (string id, AppDbContext db, UrlCache cache) =>
{
    var info = cache.Get(id);
    if (info is null)
    {
        return Results.NotFound(new { detail = "URL not found" });
    }

    var newAccesses = cache.IncrementAccesses(id);
    await db.Urls.Where(u => u.Id == id)
                 .ExecuteUpdateAsync(setters => setters.SetProperty(u => u.Accesses, newAccesses));

    return Results.Redirect(info.Original, false);
}).RequireAuthorization();

app.MapDelete($"{prefix}/{{id}}", async (string id, AppDbContext db, UrlCache cache) =>
{
    if (!cache.Exists(id))
    {
        return Results.BadRequest(new { detail = "URL not found" });
    }

    await db.Urls.Where(u => u.Id == id).ExecuteDeleteAsync();
    cache.Delete(id);

    return Results.Ok();
}).RequireAuthorization();

Console.WriteLine($"Starting ENCURTADOR at port {httpPort}");
app.Run();
