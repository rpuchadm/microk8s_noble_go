# microk8s_noble_go

go mod init go-nombre
go mod tidy

docker build -t go-nombre-service:latest .
docker tag go-nombre-service:latest localhost:32000/go-nombre-service:latest
docker push localhost:32000/go-nombre-service:latest

microk8s kubectl rollout restart deploy go-nombre-deployment -n mytestns
