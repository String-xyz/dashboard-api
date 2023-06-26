
data "aws_ssm_parameter" "datadog" {
  name = "datadog-key"
}

data "aws_ssm_parameter" "string_encryption_secret" {
  name = "string-encryption-secret"  
}

data "aws_ssm_parameter" "member_role_owner_id" {
  name = "platform-member-role-owner-id"
}

data "aws_ssm_parameter" "member_role_member_id" {
  name = "platform-member-role-member-id"
}

data "aws_ssm_parameter" "member_role_admin_id" {
  name = "platform-member-role-admin-id"
}

data "aws_ssm_parameter" "jwt_secret" {
  name = "dashboard-jwt-secret"
}

data "aws_ssm_parameter" "fingerprint_api_key" { 
  name = "fingerprint-api-key"
}

data "aws_ssm_parameter" "sendgrid_api_key" {
  name = "sendgrid-api-key"
}

data "aws_ssm_parameter" "db_password" {
  name = "string-pg-db-password"
}

data "aws_ssm_parameter" "db_username" {
  name = "string-pg-db-username"
}

data "aws_ssm_parameter" "db_name" {
  name = "string-pg-db-name"
}

data "aws_ssm_parameter" "db_host" {
  name = "pg-cluster-write-host-url"
}

data "aws_ssm_parameter" "redis_auth_token" {
  name = "redis-auth-token"
}

data "aws_ssm_parameter" "redis_host_url" {
  name  = "redis-host-url"
}

data "aws_kms_key" "kms_key" {
  key_id = "alias/main-kms-key"
}
