resource "readarr_host" "test" {
  launch_browser  = true
  port            = 8989
  url_base        = ""
  bind_address    = "*"
  application_url = ""
  instance_name   = "Readarr"
  proxy = {
    enabled = false
  }
  ssl = {
    enabled                = false
    certificate_validation = "enabled"
  }
  logging = {
    log_level = "info"
  }
  backup = {
    folder    = "/backup"
    interval  = 5
    retention = 10
  }
  authentication = {
    method = "none"
  }
  update = {
    mechanism = "docker"
    branch    = "develop"
  }
}