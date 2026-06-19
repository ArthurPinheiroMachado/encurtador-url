using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace Encurtador.Url.Models;

[Table("url")]
public class UrlItem
{
    [Key]
    [Column("id")]
    public string Id { get; set; } = string.Empty;

    [Required]
    [Column("original")]
    public string Original { get; set; } = string.Empty;

    [Column("accesses")]
    public long Accesses { get; set; }
}
