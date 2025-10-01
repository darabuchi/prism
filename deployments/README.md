# 部署配置

## 说明

包含各种部署方式的配置文件。

## 子目录

### docker/
- `Dockerfile`: Docker 镜像构建
- `docker-compose.yml`: Docker Compose 配置

### systemd/
- `prism.service`: Systemd 服务配置

### kubernetes/
- `deployment.yaml`: K8s Deployment
- `service.yaml`: K8s Service

## 部署方式

### Docker
```bash
cd deployments/docker
docker-compose up -d
```

### Systemd
```bash
sudo cp deployments/systemd/prism.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl start prism
sudo systemctl enable prism
```

### Kubernetes
```bash
kubectl apply -f deployments/kubernetes/
```

## 参考文档

- [部署与运维](../docs/详细设计/部署与运维.md)
