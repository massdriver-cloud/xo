resource "local_file" "foo" {
  content     = <<EOF
# ${var.content.header}

---
${var.content.body}
EOF
  filename = var.path
}