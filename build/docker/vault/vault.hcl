listener "tcp" {
  address          = "0.0.0.0:8200"
  tls_disable      = 0
  tls_cert_file    = "/etc/vault/certs/fullchain.pem"
  tls_key_file     = "/etc/vault/certs/privkey.pem"
}

storage "file" {
  path = "/vault/file"
}

ui = true