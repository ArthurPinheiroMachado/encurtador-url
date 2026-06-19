using Microsoft.EntityFrameworkCore;
using Encurtador.Url.Models;

namespace Encurtador.Url.Data;

public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options) : base(options) { }

    public DbSet<UrlItem> Urls => Set<UrlItem>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<UrlItem>(entity =>
        {
            entity.ToTable("url");
            entity.HasKey(e => e.Id);
            entity.HasIndex(e => e.Original).IsUnique();
        });
    }
}
