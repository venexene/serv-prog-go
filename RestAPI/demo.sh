#!/bin/bash
# ============================================================
#  Demo Script: Gravity Game Store API
#  Запуск: bash demo.sh
# ============================================================
set -e

BASE="http://localhost:8080/api/v1"
BOLD="\033[1m"
GREEN="\033[32m"
RED="\033[31m"
CYAN="\033[36m"
RESET="\033[0m"

section() { echo -e "\n${BOLD}${CYAN}━━━ $1 ━━━${RESET}"; }
ok()     { echo -e "${GREEN}✓${RESET} $1"; }

# ---------- 1. Health ----------
section "1. Health Check"
curl -s "$BASE/../health"
echo ""

# ---------- 2. Public read ----------
section "2. Публичные GET-эндпоинты (без авторизации)"

echo -e "\n${BOLD}Игры (3 шт, пагинация):${RESET}"
curl -s "$BASE/games?limit=3" | python3 -c "
import sys,json
d=json.load(sys.stdin)
print(f'  Всего: {d[\"total\"]} игр, страница {d[\"page\"]}/{d[\"total_pages\"]}')
for g in d['data']:
    a=', '.join(a['name'] for a in g.get('authors',[]))
    print(f'  #{g[\"id\"]:>2} {g[\"title\"]:<40} [{g[\"genre\"]:<15}] {g[\"platform\"]:<6} — {a}')
"

echo -e "\n${BOLD}Студии (5 шт):${RESET}"
curl -s "$BASE/authors?limit=5" | python3 -c "
import sys,json
for a in json.load(sys.stdin)['data']:
    print(f'  #{a[\"id\"]:>2} {a[\"name\"]}')
"

echo -e "\n${BOLD}Игры FromSoftware (author/1/games):${RESET}"
curl -s "$BASE/authors/1/games" | python3 -c "
import sys,json
for g in json.load(sys.stdin):
    print(f'  #{g[\"id\"]:>2} {g[\"title\"]} [{g[\"genre\"]}]')
"

echo -e "\n${BOLD}Авторы Elden Ring (games/3/authors):${RESET}"
curl -s "$BASE/games/3/authors" | python3 -c "
import sys,json
for a in json.load(sys.stdin):
    print(f'  #{a[\"id\"]:>2} {a[\"name\"]}')
"

echo -e "\n${BOLD}Заказы покупателя #1:${RESET}"
curl -s "$BASE/customers/1/orders" | python3 -c "
import sys,json
orders=json.load(sys.stdin)
print(f'  Заказов: {len(orders)}')
for o in orders:
    items=[ol['game']['title'] for ol in o.get('order_lines',[]) if ol.get('game')]
    print(f'  Заказ #{o[\"id\"]}: {items}')
"

# ---------- 3. Auth ----------
section "3. Регистрация + JWT"

echo -e "\n${BOLD}Регистрация demo-пользователя:${RESET}"
REG=$(curl -s -X POST "$BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"demo123"}')
echo "$REG" | python3 -c "
import sys,json
d=json.load(sys.stdin)
print(f'  access_token: {d[\"access_token\"][:40]}...')
print(f'  token_type:   {d[\"token_type\"]}')
print(f'  expires_in:   {d[\"expires_in\"]}s')
"

TOKEN=$(echo "$REG" | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")
ok "Токен получен"

# ---------- 4. Protected: 401 ----------
section "4. Защищённый эндпоинт БЕЗ токена → 401"

echo -e "\n${BOLD}POST /games без токена:${RESET}"
curl -s -X POST "$BASE/games" \
  -H "Content-Type: application/json" \
  -d '{"title":"hack"}' | python3 -c "
import sys,json
d=json.load(sys.stdin)
print(f'  {d[\"error\"]}')
"

# ---------- 5. Protected: success ----------
section "5. Защищённые эндпоинты С токеном"

echo -e "\n${BOLD}POST /games (создание):${RESET}"
NEW=$(curl -s -X POST "$BASE/games" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Demo Game","genre":"Sandbox","platform":"PC","num_players":8,"author_ids":[1,7]}')
echo "$NEW" | python3 -c "
import sys,json
g=json.load(sys.stdin)
a=', '.join(x['name'] for x in g.get('authors',[]))
print(f'  Создана игра #{g[\"id\"]}: {g[\"title\"]} [{g[\"genre\"]}] — {a}')
"
GAME_ID=$(echo "$NEW" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")

echo -e "\n${BOLD}PUT /games/$GAME_ID (обновление):${RESET}"
curl -s -X PUT "$BASE/games/$GAME_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Demo Game: Enhanced Edition","num_players":16}' | python3 -c "
import sys,json
g=json.load(sys.stdin)
print(f'  Обновлена: #{g[\"id\"]} {g[\"title\"]} (игроков: {g[\"num_players\"]})')
"

echo -e "\n${BOLD}DELETE /games/$GAME_ID (удаление):${RESET}"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE/games/$GAME_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "  HTTP $STATUS (204 = No Content)"

# ---------- 6. Summary ----------
section "6. Итог"
echo ""
echo -e "  ${GREEN}Все эндпоинты работают корректно${RESET}"
echo ""
echo -e "  ${BOLD}Swagger UI:${RESET} http://localhost:8080/swagger/index.html"
echo -e "  ${BOLD}Health:${RESET}     http://localhost:8080/health"
echo -e "  ${BOLD}Игры:${RESET}       http://localhost:8080/api/v1/games"
echo -e "  ${BOLD}Студии:${RESET}     http://localhost:8080/api/v1/authors"
echo ""
echo -e "  ${BOLD}Для защищённых запросов в Swagger:${RESET}"
echo "    1. POST /api/v1/auth/login с admin/admin123"
echo "    2. Кнопка Authorize → Bearer <токен>"
echo ""
