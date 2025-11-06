#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Testing Gin CRUD API...${NC}"
echo "========================"

# Base URL
BASE_URL="http://localhost:8080/api/v1"

# Test 1: Health Check
echo -e "\n${GREEN}1. Testing Health Check${NC}"
curl -s http://localhost:8080/health | jq .

# Test 2: List Departments (should be empty)
echo -e "\n${GREEN}2. List Departments (empty)${NC}"
curl -s $BASE_URL/departments | jq .

# Test 3: Create Department
echo -e "\n${GREEN}3. Create Department (Engineering)${NC}"
DEPT_RESPONSE=$(curl -s -X POST $BASE_URL/departments \
  -H "Content-Type: application/json" \
  -d '{"name": "Engineering"}')
echo "$DEPT_RESPONSE" | jq .
DEPT_ID=$(echo "$DEPT_RESPONSE" | jq -r .id)
echo "Department ID: $DEPT_ID"

# Test 4: Create Another Department
echo -e "\n${GREEN}4. Create Department (Marketing)${NC}"
DEPT2_RESPONSE=$(curl -s -X POST $BASE_URL/departments \
  -H "Content-Type: application/json" \
  -d '{"name": "Marketing"}')
echo "$DEPT2_RESPONSE" | jq .
DEPT2_ID=$(echo "$DEPT2_RESPONSE" | jq -r .id)

# Test 5: List All Departments
echo -e "\n${GREEN}5. List All Departments${NC}"
curl -s $BASE_URL/departments | jq .

# Test 6: Get Specific Department
echo -e "\n${GREEN}6. Get Department by ID${NC}"
curl -s $BASE_URL/departments/$DEPT_ID | jq .

# Test 7: Update Department
echo -e "\n${GREEN}7. Update Department Name${NC}"
curl -s -X PUT $BASE_URL/departments/$DEPT_ID \
  -H "Content-Type: application/json" \
  -d '{"name": "Software Engineering"}' | jq .

# Test 8: List Employees (should be empty)
echo -e "\n${GREEN}8. List Employees (empty)${NC}"
curl -s $BASE_URL/employees | jq .

# Test 9: Create Employee
echo -e "\n${GREEN}9. Create Employee (John Doe)${NC}"
EMP_RESPONSE=$(curl -s -X POST $BASE_URL/employees \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"John Doe\", \"email\": \"john@example.com\", \"department_id\": \"$DEPT_ID\"}")
echo "$EMP_RESPONSE" | jq .
EMP_ID=$(echo "$EMP_RESPONSE" | jq -r .id)

# Test 10: Create Another Employee
echo -e "\n${GREEN}10. Create Employee (Jane Smith)${NC}"
EMP2_RESPONSE=$(curl -s -X POST $BASE_URL/employees \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Jane Smith\", \"email\": \"jane@example.com\", \"department_id\": \"$DEPT_ID\"}")
echo "$EMP2_RESPONSE" | jq .
EMP2_ID=$(echo "$EMP2_RESPONSE" | jq -r .id)

# Test 11: Create Employee in Marketing
echo -e "\n${GREEN}11. Create Employee (Bob Wilson in Marketing)${NC}"
EMP3_RESPONSE=$(curl -s -X POST $BASE_URL/employees \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Bob Wilson\", \"email\": \"bob@example.com\", \"department_id\": \"$DEPT2_ID\"}")
echo "$EMP3_RESPONSE" | jq .
EMP3_ID=$(echo "$EMP3_RESPONSE" | jq -r .id)

# Test 12: List All Employees
echo -e "\n${GREEN}12. List All Employees${NC}"
curl -s $BASE_URL/employees | jq .

# Test 13: Get Specific Employee
echo -e "\n${GREEN}13. Get Employee by ID${NC}"
curl -s $BASE_URL/employees/$EMP_ID | jq .

# Test 14: Update Employee
echo -e "\n${GREEN}14. Update Employee (change department)${NC}"
curl -s -X PUT $BASE_URL/employees/$EMP_ID \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"John Doe\", \"email\": \"john.doe@example.com\", \"department_id\": \"$DEPT2_ID\"}" | jq .

# Test 15: Try to update employee with invalid department
echo -e "\n${GREEN}15. Update Employee with Invalid Department (should fail)${NC}"
curl -s -X PUT $BASE_URL/employees/$EMP2_ID \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Smith", "email": "jane@example.com", "department_id": "invalid-dept-id"}' | jq .

# Test 16: Delete Single Employee
echo -e "\n${GREEN}16. Delete Employee (Jane)${NC}"
curl -s -X DELETE $BASE_URL/employees/$EMP2_ID
echo "Status: Deleted (204 No Content expected)"

# Test 17: Verify Employee Deleted
echo -e "\n${GREEN}17. Try to Get Deleted Employee (should fail)${NC}"
curl -s $BASE_URL/employees/$EMP2_ID | jq .

# Test 18: List Employees After Delete
echo -e "\n${GREEN}18. List Employees After Delete${NC}"
curl -s $BASE_URL/employees | jq .

# Test 19: Delete Department with Employees (CASCADE DELETE)
echo -e "\n${GREEN}19. Delete Department with Employees (Marketing - CASCADE)${NC}"
echo "Marketing department has 2 employees (John and Bob)"
curl -s -X DELETE $BASE_URL/departments/$DEPT2_ID
echo "Status: Deleted (204 No Content expected)"

# Test 20: Verify Cascade Delete
echo -e "\n${GREEN}20. List All Employees After Cascade Delete${NC}"
echo "John and Bob should be deleted (only Engineering department employees remain)"
curl -s $BASE_URL/employees | jq .

# Test 21: List Departments After Delete
echo -e "\n${GREEN}21. List Departments After Delete${NC}"
curl -s $BASE_URL/departments | jq .

# Test 22: Try to get non-existent department
echo -e "\n${GREEN}22. Get Non-existent Department (should fail)${NC}"
curl -s $BASE_URL/departments/non-existent-id | jq .

# Test 23: Try to get non-existent employee
echo -e "\n${GREEN}23. Get Non-existent Employee (should fail)${NC}"
curl -s $BASE_URL/employees/non-existent-id | jq .

# Test 24: Create employee with invalid email format
echo -e "\n${GREEN}24. Create Employee with Invalid Email (should fail)${NC}"
curl -s -X POST $BASE_URL/employees \
  -H "Content-Type: application/json" \
  -d "{\"name\": \"Invalid User\", \"email\": \"not-an-email\", \"department_id\": \"$DEPT_ID\"}" | jq .

# Test 25: Create employee without required fields
echo -e "\n${GREEN}25. Create Employee without Name (should fail)${NC}"
curl -s -X POST $BASE_URL/employees \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"test@example.com\", \"department_id\": \"$DEPT_ID\"}" | jq .

echo -e "\n${YELLOW}Testing Complete!${NC}"