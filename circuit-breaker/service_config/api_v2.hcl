service {
  name = "api"
  id = "api-v2"
  address = "10.5.0.5"
  port = 9092

  tags = ["v2"]
  meta = {
    version = "2"
  }

  connect {
    sidecar_service {
      port = 20000

      check {
        name     = "Connect Envoy Sidecar"
        tcp      = "10.5.0.5:20000"
        interval = "10s"
      }
    }
  }
}
