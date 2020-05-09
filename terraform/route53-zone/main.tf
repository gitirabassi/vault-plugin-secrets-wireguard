variable domain {}

variable a_records {
  type    = map(string)
  default = {}
}

resource "aws_route53_zone" "main" {
  name = var.domain
}

resource "aws_route53_record" "main" {
  for_each = var.a_records
  zone_id  = aws_route53_zone.main.id
  name     = each.key
  type     = "A"
  ttl      = "300"
  records  = [each.value]
}
