variable "{{serviceRole}}_env_vars" {
  default = "[]"
}

variable "{{serviceRole}}_docker_image" {}

module "{{serviceRole}}" {
  source = "github.com/Originate/exosphere.git//terraform//aws//worker-service?ref={{terraformCommitHash}}"

  name = "{{serviceRole}}"

  cluster_id            = "${data.terraform_remote_state.main_infrastructure.ecs_cluster_id}"
  cpu                   = "{{cpu}}"
  desired_count         = 1
  docker_image          = "${var.{{serviceRole}}_docker_image}"
  env                   = "${var.env}"
  environment_variables = "${var.{{serviceRole}}_env_vars}"
  memory_reservation    = "{{memory}}"
  region                = "${data.terraform_remote_state.main_infrastructure.region}"
}
