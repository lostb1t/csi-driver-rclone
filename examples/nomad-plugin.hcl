job "plugin-rclone-csi" {
  datacenters = ["local"]
  type = "system"


  group "plugin-rclone-csi" {
    count = 1
    task "plugin" {
      driver = "docker"

      config {
        image      = "ghcr.io/lostb1t/csi-driver-rclone"
        force_pull = true
        args = [
            "-logtostderr=true",
            "--endpoint=unix:///csi/csi.sock",
            "--nodeid=${node.unique.id}",
        ]
      }

      resources {
        cpu    = 100
        memory = 250
      }

      csi_plugin {
        id = "csi-rclone"
        type = "monolith"
        mount_dir      = "/csi"
        health_timeout = "30s"
      }

       }
  }
}
