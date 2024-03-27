#!/bin/bash

function install_deps {
    yum install -y yum-utils shadow-utils
    yum update -y
    yum-config-manager --add-repo https://rpm.releases.hashicorp.com/AmazonLinux/hashicorp.repo
    yum -y install vault-enterprise <"/dev/null"
    yum -y install jq <"/dev/null"
    yum -y install nc <"/dev/null"
    yum install awscli -y
    INSTANCE_ID=$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
    PRIVATE_IP=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)
    REGION=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone | sed 's/[a-z]$//')
    TAG_NAME=$(aws ec2 describe-tags --filters "Name=resource-id,Values=$INSTANCE_ID" "Name=key,Values=Name" --region $REGION --output text | cut -f5)
    echo "export PS1='[$TAG_NAME@$PRIVATE_IP \W]\$ '" >> /home/ec2-user/.bash_profile
}

# initializes a single server vault instance raft
function init_vault {
    echo "PASTE_LICENSE_HERE" >/etc/vault.d/vault.hclic
    cat <<EOF1 >/etc/vault.d/vault.hcl
storage "raft" {
  path    = "/opt/vault/data"
  node_id = "$(hostname)"
}

listener "tcp" {
  address         = "0.0.0.0:8200"
  tls_disable     = true
}

license_path = "/etc/vault.d/vault.hclic"
api_addr = "http://$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4):8200"
cluster_addr = "http://$(hostname):8201"
log_level = "trace"
EOF1
    echo 'export VAULT_ADDR=http://127.0.0.1:8200' >>/etc/environment
    echo "export AWS_DEFAULT_REGION=$(curl -s http://169.254.169.254/latest/dynamic/instance-identity/document | /usr/bin/jq -r '.region')" >>/etc/environment
    export VAULT_ADDR=http://127.0.0.1:8200
    systemctl start vault
    vault operator init -key-shares=1 -key-threshold=1 >/home/ec2-user/keys
    echo $(grep 'Key 1:' /home/ec2-user/keys | awk '{print $NF}') >/home/ec2-user/unseal
    vault operator unseal $(cat home/ec2-user/unseal)
    echo $(grep 'Initial Root Token:' /home/ec2-user/keys | awk '{print $NF}') >/home/ec2-user/root
    rm /home/ec2-user/keys
    cat <<EOF2 >>/home/ec2-user/.bash_profile
function login () {
    vault login -<root
}
login
EOF2
}

install_deps
init_vault