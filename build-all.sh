#!/bin/bash

echo "======================================"
echo "🏗️  BUILD ALL - Terminal Versions"
echo "======================================"
echo ""

# Cores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Backup do .env original
echo -e "${BLUE}📦 Fazendo backup do .env...${NC}"
cp .env .env.backup

# Build 1: Versão VERTICAL (Portrait - Tigrinho)
echo ""
echo -e "${YELLOW}======================================"
echo "🐯 BUILD 1: VERTICAL (Tigrinho)"
echo "======================================${NC}"
sed -i 's/SCREEN_ORIENTATION=.*/SCREEN_ORIENTATION=portrait/' .env
go build -o terminal-vertical main.go

if [ $? -eq 0 ]; then
    SIZE_V=$(du -h terminal-vertical | cut -f1)
    echo -e "${GREEN}✅ terminal-vertical compilado com sucesso! ($SIZE_V)${NC}"
else
    echo "❌ Erro ao compilar terminal-vertical"
    mv .env.backup .env
    exit 1
fi

# Build 2: Versão HORIZONTAL (Landscape - Empire)
echo ""
echo -e "${YELLOW}======================================"
echo "🏛️  BUILD 2: HORIZONTAL (Empire)"
echo "======================================${NC}"
sed -i 's/SCREEN_ORIENTATION=.*/SCREEN_ORIENTATION=landscape/' .env
go build -o terminal-horizontal main.go

if [ $? -eq 0 ]; then
    SIZE_H=$(du -h terminal-horizontal | cut -f1)
    echo -e "${GREEN}✅ terminal-horizontal compilado com sucesso! ($SIZE_H)${NC}"
else
    echo "❌ Erro ao compilar terminal-horizontal"
    mv .env.backup .env
    exit 1
fi

# Restaurar .env original
echo ""
echo -e "${BLUE}🔄 Restaurando .env original...${NC}"
mv .env.backup .env

# Resumo final
echo ""
echo "======================================"
echo -e "${GREEN}🎉 BUILD COMPLETO!${NC}"
echo "======================================"
echo ""
echo "📦 Binários criados:"
echo ""
ls -lh terminal-vertical terminal-horizontal | awk '{print "  " $9 " → " $5}'
echo ""
echo -e "${BLUE}📱 USO:${NC}"
echo "  Vertical (Tigrinho):   ./terminal-vertical"
echo "  Horizontal (Empire):   ./terminal-horizontal"
echo ""
echo -e "${YELLOW}📋 INFORMAÇÕES:${NC}"
echo "  Vertical:   1080x1920 (portrait)  🐯"
echo "  Horizontal: 1920x1080 (landscape) 🏛️"
echo ""
echo "======================================"
