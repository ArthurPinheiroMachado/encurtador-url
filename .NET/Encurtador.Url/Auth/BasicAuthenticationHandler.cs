using System.Security.Claims;
using System.Security.Cryptography;
using System.Text;
using System.Text.Encodings.Web;
using Microsoft.AspNetCore.Authentication;
using Microsoft.Extensions.Options;

namespace Encurtador.Url.Auth;

public class BasicAuthenticationOptions : AuthenticationSchemeOptions { }

public class BasicAuthenticationHandler : AuthenticationHandler<BasicAuthenticationOptions>
{
    public const string SchemeName = "Basic";

    public BasicAuthenticationHandler(
        IOptionsMonitor<BasicAuthenticationOptions> options,
        ILoggerFactory logger,
        UrlEncoder encoder)
        : base(options, logger, encoder) { }

    protected override Task<AuthenticateResult> HandleAuthenticateAsync()
    {
        var authHeader = Request.Headers.Authorization.FirstOrDefault();

        if (string.IsNullOrEmpty(authHeader) || !authHeader.StartsWith("Basic ", StringComparison.OrdinalIgnoreCase))
        {
            return Task.FromResult(AuthenticateResult.Fail("Missing or invalid Authorization header"));
        }

        var token = authHeader["Basic ".Length..].Trim();
        string decoded;

        try
        {
            var bytes = Convert.FromBase64String(token);
            decoded = Encoding.UTF8.GetString(bytes);
        }
        catch
        {
            return Task.FromResult(AuthenticateResult.Fail("Invalid Base64 encoding"));
        }

        var parts = decoded.Split(':', 2);
        if (parts.Length != 2)
        {
            return Task.FromResult(AuthenticateResult.Fail("Invalid Authorization format"));
        }

        var username = parts[0];
        var password = parts[1];

        var expectedUser = Environment.GetEnvironmentVariable("USER") ?? "user";
        var expectedPass = Environment.GetEnvironmentVariable("PASS") ?? "pass123";

        var userBytes = Encoding.UTF8.GetBytes(username);
        var expectedUserBytes = Encoding.UTF8.GetBytes(expectedUser);
        var passBytes = Encoding.UTF8.GetBytes(password);
        var expectedPassBytes = Encoding.UTF8.GetBytes(expectedPass);

        var userOk = userBytes.Length == expectedUserBytes.Length &&
            CryptographicOperations.FixedTimeEquals(userBytes, expectedUserBytes);
        var passOk = passBytes.Length == expectedPassBytes.Length &&
            CryptographicOperations.FixedTimeEquals(passBytes, expectedPassBytes);

        if (!userOk || !passOk)
        {
            return Task.FromResult(AuthenticateResult.Fail("Invalid credentials"));
        }

        var claims = new[] { new Claim(ClaimTypes.Name, username) };
        var identity = new ClaimsIdentity(claims, Scheme.Name);
        var principal = new ClaimsPrincipal(identity);
        var ticket = new AuthenticationTicket(principal, Scheme.Name);

        return Task.FromResult(AuthenticateResult.Success(ticket));
    }
}
