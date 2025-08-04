#!/bin/bash

echo "🔨 Build..."
GOOS=linux GOARCH=amd64 go build -o monad-server ./cmd/server

echo "🛑 Stop service..."
ssh ubuntu@54.38.183.183 "sudo systemctl stop monad-portfolio-dev"

echo "📤 Upload..."
scp monad-server ubuntu@54.38.183.183:~/

echo "🔧 Rendre exécutable..."
ssh ubuntu@54.38.183.183 "chmod +x ~/monad-server"

echo "🔄 Restart..."
ssh ubuntu@54.38.183.183 "sudo systemctl restart monad-portfolio-dev"

echo "✅ Status:"
ssh ubuntu@54.38.183.183 "sudo systemctl status monad-portfolio-dev --no-pager"

echo "📋 Derniers logs:"
ssh ubuntu@54.38.183.183 "sudo journalctl -u monad-portfolio-dev -n 20 --no-pager"

echo "Déployé !"