#!/bin/bash
# Docker Helper Script for DKL Automationgo

set -e

# Check if command is provided
if [ $# -lt 1 ]; then
    echo "Usage: $0 <command> [arg]"
    echo "Commands: start, stop, restart, logs, build, shell, db, backup, restore"
    exit 1
fi

COMMAND=$1
ARG=$2

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH. Please install Docker."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Error: Docker Compose is not installed or not in PATH. Please install Docker Compose."
    exit 1
fi

# Function to check if containers are running
function check_containers_running() {
    local containers=$(docker ps --filter "name=dklautomationgo" --format "{{.Names}}")
    if [ -z "$containers" ]; then
        return 1
    else
        return 0
    fi
}

# Execute the requested command
case $COMMAND in
    "start")
        echo -e "\033[0;32mStarting containers...\033[0m"
        docker-compose up -d
        echo -e "\033[0;32mContainers started. Application is available at http://localhost:8080\033[0m"
        ;;
    "stop")
        echo -e "\033[0;33mStopping containers...\033[0m"
        docker-compose down
        echo -e "\033[0;33mContainers stopped.\033[0m"
        ;;
    "restart")
        echo -e "\033[0;33mRestarting containers...\033[0m"
        docker-compose down
        docker-compose up -d
        echo -e "\033[0;32mContainers restarted. Application is available at http://localhost:8080\033[0m"
        ;;
    "logs")
        if [ "$ARG" == "db" ]; then
            echo -e "\033[0;36mShowing database logs...\033[0m"
            docker-compose logs -f db
        else
            echo -e "\033[0;36mShowing application logs...\033[0m"
            docker-compose logs -f app
        fi
        ;;
    "build")
        echo -e "\033[0;32mBuilding and starting containers...\033[0m"
        docker-compose up -d --build
        echo -e "\033[0;32mContainers built and started. Application is available at http://localhost:8080\033[0m"
        ;;
    "shell")
        if ! check_containers_running; then
            echo -e "\033[0;31mError: Containers are not running. Start them first with: ./scripts/docker-helper.sh start\033[0m"
            exit 1
        fi
        echo -e "\033[0;36mOpening shell in app container...\033[0m"
        docker-compose exec app sh
        ;;
    "db")
        if ! check_containers_running; then
            echo -e "\033[0;31mError: Containers are not running. Start them first with: ./scripts/docker-helper.sh start\033[0m"
            exit 1
        fi
        echo -e "\033[0;36mOpening PostgreSQL shell...\033[0m"
        docker-compose exec db psql -U postgres -d dklautomationgo
        ;;
    "backup")
        if ! check_containers_running; then
            echo -e "\033[0;31mError: Containers are not running. Start them first with: ./scripts/docker-helper.sh start\033[0m"
            exit 1
        fi
        TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
        BACKUP_FILE="backup_$TIMESTAMP.sql"
        echo -e "\033[0;32mCreating database backup to $BACKUP_FILE...\033[0m"
        docker-compose exec -T db pg_dump -U postgres dklautomationgo > $BACKUP_FILE
        echo -e "\033[0;32mBackup created: $BACKUP_FILE\033[0m"
        ;;
    "restore")
        if ! check_containers_running; then
            echo -e "\033[0;31mError: Containers are not running. Start them first with: ./scripts/docker-helper.sh start\033[0m"
            exit 1
        fi
        if [ -z "$ARG" ]; then
            echo -e "\033[0;31mError: Please specify a backup file to restore. Usage: ./scripts/docker-helper.sh restore backup_file.sql\033[0m"
            exit 1
        fi
        if [ ! -f "$ARG" ]; then
            echo -e "\033[0;31mError: Backup file not found: $ARG\033[0m"
            exit 1
        fi
        echo -e "\033[0;33mRestoring database from $ARG...\033[0m"
        cat $ARG | docker-compose exec -T db psql -U postgres -d dklautomationgo
        echo -e "\033[0;32mDatabase restored from $ARG\033[0m"
        ;;
    *)
        echo -e "\033[0;31mError: Unknown command: $COMMAND\033[0m"
        echo "Usage: $0 <command> [arg]"
        echo "Commands: start, stop, restart, logs, build, shell, db, backup, restore"
        exit 1
        ;;
esac 