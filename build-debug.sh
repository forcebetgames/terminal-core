#!/bin/bash

echo "======================================"
echo "🐛 BUILD DEBUG - Terminal com DevTools"
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

echo ""
echo -e "${YELLOW}======================================"
echo "🐛 HABILITANDO DEVTOOLS"
echo "======================================${NC}"

# Habilita DevTools no .env
sed -i 's/ENABLE_DEVTOOLS=.*/ENABLE_DEVTOOLS=true/' .env

# Build versão debug VERTICAL
echo ""
echo -e "${YELLOW}🐯 BUILD DEBUG VERTICAL (Tigrinho + DevTools)${NC}"
sed -i 's/SCREEN_ORIENTATION=.*/SCREEN_ORIENTATION=portrait/' .env
go build -o terminal-debug-vertical main.go

if [ $? -eq 0 ]; then
    SIZE=$(du -h terminal-debug-vertical | cut -f1)
    echo -e "${GREEN}✅ terminal-debug-vertical compilado! ($SIZE)${NC}"
else
    echo "❌ Erro ao compilar terminal-debug-vertical"
    mv .env.backup .env
    exit 1
fi

# Build versão debug HORIZONTAL
echo ""
echo -e "${YELLOW}🏛️  BUILD DEBUG HORIZONTAL (Empire + DevTools)${NC}"
sed -i 's/SCREEN_ORIENTATION=.*/SCREEN_ORIENTATION=landscape/' .env
go build -o terminal-debug-horizontal main.go

if [ $? -eq 0 ]; then
    SIZE=$(du -h terminal-debug-horizontal | cut -f1)
    echo -e "${GREEN}✅ terminal-debug-horizontal compilado! ($SIZE)${NC}"
else
    echo "❌ Erro ao compilar terminal-debug-horizontal"
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
echo -e "${GREEN}🎉 BUILD DEBUG COMPLETO!${NC}"
echo "======================================"
echo ""
echo "📦 Binários DEBUG criados:"
echo ""
ls -lh terminal-debug-vertical terminal-debug-horizontal | awk '{print "  " $9 " → " $5}'
echo ""
echo -e "${BLUE}🐛 USO (DevTools habilitado):${NC}"
echo "  Vertical:   ./terminal-debug-vertical"
echo "  Horizontal: ./terminal-debug-horizontal"
echo ""
echo -e "${YELLOW}⚠️  IMPORTANTE:${NC}"
echo "  - DevTools abre automaticamente (F12)"
echo "  - Modo WINDOWED (não kiosk)"
echo "  - Use para DEBUG/DESENVOLVIMENTO"
echo ""
echo -e "${GREEN}💡 LOGS:${NC}"
echo "  Console do navegador (F12) mostra:"
echo "  - Erros JavaScript"
echo "  - Network requests"
echo "  - Console.log do site"
echo ""
echo "======================================"
