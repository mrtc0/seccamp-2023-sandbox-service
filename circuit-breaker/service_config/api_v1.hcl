service {
  name = "api"
  id = "api-v1"
  address = "10.5.0.4"
  port = 9092

  tags = ["v1"]
  meta = {
    version = "1"
  }

  connect {
    sidecar_service {
      port = 20000

      check {
        name     = "Connect Envoy Sidecar"
        tcp      = "10.5.0.4:20000"
        interval = "10s"
      }
    }
  }
}
