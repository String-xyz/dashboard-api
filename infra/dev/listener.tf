module "alb_listener_acm" {
  source      = "git@github.com:String-xyz/terraform-modules.git//acm"
  domain_name = "${local.service_name}.${local.root_domain}"
  aws_region  = "us-west-2"
  zone_id     = data.aws_route53_zone.root.zone_id
  tags = {
    Environment = local.env
    Name = "${local.service_name}.${local.root_domain}"
  }
}

data "aws_ssm_parameter" "alb_listener" {
  name = "string-api-alb-istener-arn"
}

resource "aws_alb_listener_certificate" "cert" {
  listener_arn    = data.aws_ssm_parameter.alb_listener.value
  certificate_arn = module.alb_listener_acm.arn
}

resource "aws_alb_target_group" "ecs_task_target_group" {
  name        = "${local.service_name}-tg"
  port        = local.container_port
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
  target_type = "ip"
  protocol    = "HTTP"

  lifecycle {
    create_before_destroy = true
  }

  health_check {
    path                = "/heartbeat"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 60
    timeout             = 30
    unhealthy_threshold = "3"
    healthy_threshold   = "3"
  }

  tags = {
    Name = "${local.service_name}-tg"
  }
}

resource "aws_alb_listener_rule" "ecs_alb_listener_rule" {
  listener_arn = data.aws_ssm_parameter.alb_listener.value
  priority     = 99
  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.ecs_task_target_group.arn
  }

  condition {
    host_header {
      values = ["${local.service_name}.${local.root_domain}"]
    }
  }
}

