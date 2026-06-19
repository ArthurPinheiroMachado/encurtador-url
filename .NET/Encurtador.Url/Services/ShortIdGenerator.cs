using System.Security.Cryptography;

namespace Encurtador.Url.Services;

public static class ShortIdGenerator
{
    private const string Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    private const int MaxAttempts = 100;

    public static string Generate(int length = 8, Func<string, bool>? exists = null)
    {
        for (int i = 0; i < MaxAttempts; i++)
        {
            var id = new char[length];
            var bytes = RandomNumberGenerator.GetBytes(length);

            for (int j = 0; j < length; j++)
            {
                id[j] = Charset[bytes[j] % Charset.Length];
            }

            var result = new string(id);
            if (exists == null || !exists(result))
            {
                return result;
            }
        }

        throw new InvalidOperationException($"Failed to generate unique ID after {MaxAttempts} attempts");
    }
}
