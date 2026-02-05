#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
BUILD_MODE="optimized"
COMMAND="build"
OUTPUT_DIR="./bin"
BINARY_NAME="byebye"

# Parse arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            debug)
                BUILD_MODE="debug"
                shift
                ;;
            optimized)
                BUILD_MODE="optimized"
                shift
                ;;
            install)
                COMMAND="install"
                shift
                ;;
            test)
                COMMAND="test"
                shift
                ;;
            clean)
                COMMAND="clean"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                show_help
                exit 1
                ;;
        esac
    done
}

show_help() {
    cat << EOF
${BLUE}byebye Build Script${NC}

Usage: ./build.sh [mode] [command]

${YELLOW}Modes:${NC}
    debug       - Build with debug symbols and no optimizations (default: optimized)
    optimized   - Build with optimizations and strip symbols (default)

${YELLOW}Commands:${NC}
    build       - Build binary (default)
    install     - Install binary to system
    test        - Run tests
    clean       - Clean build artifacts

${YELLOW}Examples:${NC}
    ./build.sh                  # Build optimized binary
    ./build.sh debug            # Build debug binary
    ./build.sh debug test       # Run tests in debug mode
    ./build.sh optimized install # Install optimized binary

EOF
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Clean build artifacts
clean() {
    print_info "Cleaning build artifacts..."
    rm -rf "${OUTPUT_DIR}"
    go clean
    print_success "Clean complete"
}

# Build the binary
build() {
    print_info "Building ${BUILD_MODE} binary..."
    
    mkdir -p "${OUTPUT_DIR}"
    
    if [ "$BUILD_MODE" = "debug" ]; then
        print_info "Debug build: symbols included, no optimizations"
        go build \
            -v \
            -gcflags="all=-N -l" \
            -o "${OUTPUT_DIR}/${BINARY_NAME}" \
            .
    else
        print_info "Optimized build: stripped, maximum optimizations"
        go build \
            -v \
            -ldflags="-s -w" \
            -tags=release \
            -trimpath \
            -o "${OUTPUT_DIR}/${BINARY_NAME}" \
            .
    fi
    
    BINARY_PATH="${OUTPUT_DIR}/${BINARY_NAME}"
    SIZE=$(du -h "${BINARY_PATH}" | cut -f1)
    print_success "Binary built: ${BINARY_PATH} (${SIZE})"
}

# Run tests
run_tests() {
    print_info "Running tests in ${BUILD_MODE} mode..."
    
    if [ "$BUILD_MODE" = "debug" ]; then
        go test -v -race -cover ./...
    else
        go test -v -race -cover -run=. ./... || true
        go test -bench=. -benchmem ./... || true
    fi
    
    print_success "Tests complete"
}

# Install binary
install_binary() {
    # First build it
    build
    
    print_info "Installing ${BINARY_NAME}..."
    
    INSTALL_PATH="${HOME}/.local/bin/${BINARY_NAME}"
    mkdir -p "${HOME}/.local/bin"
    
    cp "${OUTPUT_DIR}/${BINARY_NAME}" "${INSTALL_PATH}"
    chmod +x "${INSTALL_PATH}"
    
    if command -v sudo &> /dev/null && [ "$EUID" -ne 0 ]; then
        print_warn "Running install as non-root. Use 'sudo ./build.sh optimized install' for system-wide install"
    fi
    
    print_success "Installed to: ${INSTALL_PATH}"
    print_info "Make sure ${HOME}/.local/bin is in your PATH"
}

# Download/verify dependencies
verify_deps() {
    print_info "Verifying dependencies..."
    go mod download
    go mod tidy
    print_success "Dependencies verified"
}

# Main execution
main() {
    parse_args "$@"
    
    print_info "Build Mode: ${BUILD_MODE}"
    print_info "Command: ${COMMAND}"
    
    # Verify deps before any command
    verify_deps
    
    case $COMMAND in
        build)
            build
            ;;
        test)
            run_tests
            ;;
        install)
            install_binary
            ;;
        clean)
            clean
            ;;
        *)
            print_error "Unknown command: $COMMAND"
            show_help
            exit 1
            ;;
    esac
    
    print_success "Done!"
}

# Run main
main "$@"
