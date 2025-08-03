#!/bin/bash

echo "🔨 Build..."
GOOS=linux GOARCH=amd64 go build -o monad-server ./cmd/server

echo "📤 Upload..."
scp monad-server ubuntu@54.38.183.183:~/

echo "🔧 Rendre exécutable..."
ssh ubuntu@54.38.183.183 "chmod +x ~/monad-server"

echo "🔄 Restart..."
ssh ubuntu@54.38.183.183 "sudo systemctl restart monad-portfolio-dev"

echo "✅ Status:"
ssh ubuntu@54.38.183.183 "sudo systemctl status monad-portfolio-dev --no-pager"

echo "Déployé !"