locals {
  cluster_name       = "sandbox-core"
  env                = "sandbox"
  service_name       = "dashboard-api"
  root_domain        = "sandbox.string-api.xyz"
  domain             = "dashboard-api"
  container_port     = "3000"
  origin_id          = "dashboard-api"
  desired_task_count = "1"
  db_port            = "5432"
  redis_port         = "6379"
  memory             = 512
  cpu                = 256
  region             = "us-west-2"
}

variable "versioning" {
  type    = string
  default = "v1.1.0"
}

locals {
  task_definition = jsonencode([
    {
      name      = local.service_name
      image     = "${aws_ecr_repository.repo.repository_url}:${var.versioning}"
      essential = true,
      dockerLabels = {
        "com.datadoghq.ad.instances" : "[{\"host\":\"%%host%%\"}]",
        "com.datadoghq.ad.check_names" : "[\"${local.service_name}\"]",
     },
      portMappings = [
        {
          containerPort = 3000
        }
      ],
      secrets = [
        {
          name      = "STRING_ENCRYPTION_KEY"
          valueFrom = data.aws_ssm_parameter.string_encryption_secret.arn
        },
        {
          name      = "JWT_SECRET_KEY"
          valueFrom = data.aws_ssm_parameter.jwt_secret.arn
        },
        {
          name      = "SENDGRID_API_KEY"
          valuefrom = data.aws_ssm_parameter.sendgrid_api_key.arn
        },
        {
          name      = "DB_USERNAME"
          valueFrom = data.aws_ssm_parameter.db_username.arn
        },
        {
          name      = "DB_PASSWORD"
          valueFrom = data.aws_ssm_parameter.db_password.arn
        },
        {
          name      = "DB_HOST"
          valueFrom = data.aws_ssm_parameter.db_host.arn
        },
        {
          name      = "DB_NAME"
          valueFrom = data.aws_ssm_parameter.db_name.arn
        },
        {
          name      = "REDIS_HOST",
          valuefrom = data.aws_ssm_parameter.redis_host_url.arn
        },
        {
          name      = "REDIS_PASSWORD",
          valuefrom = data.aws_ssm_parameter.redis_auth_token.arn
        }
      ]
      environment = [
        {
          name  = "BASE_DASHBOARD_URL"
          value = "https://sandbox.string.xyz"
        },
        {
          name  = "MEMBER_ROLE_OWNER_ID"
          value = data.aws_ssm_parameter.member_role_owner_id.value
        },
        {
          name  = "MEMBER_ROLE_ADMIN_ID"
          value = data.aws_ssm_parameter.member_role_admin_id.value
        },
        {
          name  = "MEMBER_ROLE_MEMBER_ID"
          value = data.aws_ssm_parameter.member_role_member_id.value
        },
        {
          name  = "PORT"
          value = local.container_port
        },
        {
          name  = "REDIS_PORT"
          value = local.redis_port
        },
        {
          name  = "DB_PORT",
          value = local.db_port
        },
        {
          name  = "ENV"
          value = local.env
        },
        {
          name  = "AWS_REGION"
          value = local.region
        },
        {
          name  = "AWS_KMS_KEY_ID"
          value = data.aws_kms_key.kms_key.key_id
        },
        {
          name = "AUTH_EMAIL_ADDRESS"
          value = "auth@string.xyz"
        }
      ]

      logConfiguration = {
        logDriver = "awsfirelens"
        secretOptions = [{
          name      = "apiKey",
          valueFrom = data.aws_ssm_parameter.datadog.arn
        }]
        options = {
          Name             = "datadog"
          "dd_service"     = local.service_name
          "Host"           = "http-intake.logs.datadoghq.com"
          "dd_source"      = local.service_name
          "dd_message_key" = "log"
          "dd_tags"        = "project:${local.service_name}"
          "TLS"            = "on"
          "provider"       = "ecs"
        }
      }
    },
    {
      name      = "datadog-agent"
      image     = "public.ecr.aws/datadog/agent:latest"
      essential = true
      portMappings = [{
        hostPort      = 8126,
        protocol      = "tcp",
        containerPort = 8126
      }
      ]
      secrets = [{
        name      = "DD_API_KEY"
        valueFrom = data.aws_ssm_parameter.datadog.arn
      }]
      environment = [
        {
          name  = "DD_VERSION"
         value = var.versioning
        },
        {
          name  = "DD_ENV"
          value = local.env
       }]
  },
  {
    name      = "log_router"
    image     = "public.ecr.aws/aws-observability/aws-for-fluent-bit:stable"
    essential = true
    firelensConfiguration = {
      type = "fluentbit"
      options = {
        "enable-ecs-log-metadata" = "true"
      }
    }
  }
])
}
