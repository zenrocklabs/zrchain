#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Parse command line arguments
for arg in "$@"; do
    case $arg in
        -h|--help)
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  -h, --help  Show this help message"
            echo ""
            echo "This script installs all development dependencies for zrchain."
            echo "Go tool versions are pinned to what's specified in go.mod."
            exit 0
            ;;
        *)
            echo "Unknown option: $arg"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Logging functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            OS="macos"
            ;;
        Linux*)
            OS="linux"
            ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
    log_info "Detected OS: $OS"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install uv (Python package manager)
install_uv() {
    if command_exists uv; then
        log_info "uv is already installed ($(uv --version))"
        return 0
    fi

    log_info "Installing uv..."

    if command_exists curl; then
        curl -LsSf https://astral.sh/uv/install.sh | sh
    elif command_exists wget; then
        wget -qO- https://astral.sh/uv/install.sh | sh
    else
        log_error "Neither curl nor wget is available. Please install one of them first."
        exit 1
    fi

    log_info "uv installed successfully"
}

# Install just (command runner)
install_just() {
    if command_exists just; then
        log_info "just is already installed ($(just --version))"
        return 0
    fi

    log_info "Installing just..."
    
    case "$OS" in
        macos)
            if ! command_exists brew; then
                install_homebrew
            fi
            brew install just
            ;;
        linux)
            case "$PKG_MANAGER" in
                apt)
                    install_package "just" "just"
                    ;;
                dnf)
                    install_package "just" "just"
                    ;;
                yum)
                    log_info "Attempting to install just via yum..."
                    if ! sudo yum install -y just; then
                        log_warn "just not found in yum repositories. Trying cargo install..."
                        cargo install just
                    fi
                    ;;
                pacman)
                    install_package "just" "just"
                    ;;
                zypper)
                    install_package "just" "just"
                    ;;
                *)
                    log_warn "No direct package available for $PKG_MANAGER. Installing via cargo..."
                    cargo install just
                    ;;
            esac
            ;;
    esac

    log_info "just installed successfully ($(just --version))"
}

# Install Homebrew on macOS
install_homebrew() {
    log_info "Homebrew is required to install dependencies on macOS."
    read -p "Would you like to install Homebrew? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        
        # Add Homebrew to PATH for Apple Silicon Macs
        if [[ $(uname -m) == "arm64" ]]; then
            echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
            eval "$(/opt/homebrew/bin/brew shellenv)"
        fi
        
        log_info "Homebrew installed successfully"
    else
        log_error "Homebrew is required to proceed. Exiting."
        exit 1
    fi
}

# Detect Linux package manager
detect_linux_package_manager() {
    local managers=()
    
    if command_exists apt-get; then
        managers+=("apt")
    fi
    if command_exists dnf; then
        managers+=("dnf")
    fi
    if command_exists yum; then
        managers+=("yum")
    fi
    if command_exists pacman; then
        managers+=("pacman")
    fi
    if command_exists zypper; then
        managers+=("zypper")
    fi
    
    if [ ${#managers[@]} -eq 0 ]; then
        log_error "No supported package manager found"
        exit 1
    elif [ ${#managers[@]} -eq 1 ]; then
        PKG_MANAGER="${managers[0]}"
        log_info "Detected package manager: $PKG_MANAGER"
    else
        log_info "Multiple package managers detected: ${managers[*]}"
        echo "Please choose your preferred package manager:"
        select pm in "${managers[@]}"; do
            if [ -n "$pm" ]; then
                PKG_MANAGER="$pm"
                log_info "Selected package manager: $PKG_MANAGER"
                break
            else
                log_error "Invalid selection"
            fi
        done
    fi
}

# Install package using the appropriate package manager
install_package() {
    local package=$1
    local package_name=${2:-$package}
    
    if command_exists "$package" && [ "$UPDATE_MODE" = false ]; then
        log_info "$package_name is already installed"
        return 0
    fi
    
    if [ "$UPDATE_MODE" = true ] && command_exists "$package"; then
        log_info "Updating $package_name..."
    else
        log_info "Installing $package_name..."
    fi
    
    case "$OS-$PKG_MANAGER" in
        macos-*)
            if [ "$UPDATE_MODE" = true ] && command_exists "$package"; then
                brew upgrade "$package"
            else
                brew install "$package"
            fi
            ;;
        linux-apt)
            sudo apt-get update
            if [ "$UPDATE_MODE" = true ] && command_exists "$package"; then
                sudo apt-get install -y --only-upgrade "$package"
            else
                sudo apt-get install -y "$package"
            fi
            ;;
        linux-dnf)
            if [ "$UPDATE_MODE" = true ] && command_exists "$package"; then
                sudo dnf upgrade -y "$package"
            else
                sudo dnf install -y "$package"
            fi
            ;;
        linux-yum)
            if [ "$UPDATE_MODE" = true ] && command_exists "$package"; then
                sudo yum update -y "$package"
            else
                sudo yum install -y "$package"
            fi
            ;;
        linux-pacman)
            if [ "$UPDATE_MODE" = true ] && command_exists "$package"; then
                sudo pacman -Syu --noconfirm "$package"
            else
                sudo pacman -S --noconfirm "$package"
            fi
            ;;
        linux-zypper)
            if [ "$UPDATE_MODE" = true ] && command_exists "$package"; then
                sudo zypper update -y "$package"
            else
                sudo zypper install -y "$package"
            fi
            ;;
        *)
            log_error "Unsupported OS/package manager combination: $OS-$PKG_MANAGER"
            exit 1
            ;;
    esac
    
    if [ "$UPDATE_MODE" = true ]; then
        log_info "$package_name updated successfully"
    else
        log_info "$package_name installed successfully"
    fi
}

# Install Go
install_go() {
    if command_exists go; then
        log_info "Go is already installed ($(go version))"
        return 0
    fi

    log_info "Installing Go..."

    case "$OS" in
        macos)
            if ! command_exists brew; then
                install_homebrew
            fi
            brew install go
            ;;
        linux)
            case "$PKG_MANAGER" in
                apt)
                    install_package "golang-go" "Go"
                    ;;
                *)
                    install_package "go" "Go"
                    ;;
            esac
            ;;
    esac

    log_info "Go installed successfully ($(go version))"
}

# Install Docker
install_docker() {
    if command_exists docker; then
        log_info "Docker is already installed ($(docker --version))"
        return 0
    fi

    log_info "Installing Docker..."
    
    case "$OS" in
        macos)
            if [ "$UPDATE_MODE" = true ] && command_exists docker; then
                log_warn "On macOS, Docker Desktop should be updated manually."
                log_warn "Please visit: https://www.docker.com/products/docker-desktop"
                log_warn "Or use: brew upgrade --cask docker"
            else
                log_warn "On macOS, Docker Desktop needs to be installed manually."
                log_warn "Please visit: https://www.docker.com/products/docker-desktop"
                log_warn "Or use: brew install --cask docker"
            fi
            read -p "Press Enter to continue after installing/updating Docker Desktop, or Ctrl+C to exit..."
            ;;
        linux)
            case "$PKG_MANAGER" in
                apt)
                    sudo apt-get update
                    if [ "$UPDATE_MODE" = true ] && command_exists docker; then
                        sudo apt-get install -y --only-upgrade docker-ce docker-ce-cli containerd.io || {
                            # If upgrade fails, ensure repo is set up and try again
                            sudo apt-get install -y ca-certificates curl gnupg
                            sudo install -m 0755 -d /etc/apt/keyrings
                            curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
                            sudo chmod a+r /etc/apt/keyrings/docker.gpg
                            echo \
                              "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
                              "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
                              sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
                            sudo apt-get update
                            sudo apt-get install -y --only-upgrade docker-ce docker-ce-cli containerd.io
                        }
                    else
                        sudo apt-get install -y ca-certificates curl gnupg
                        sudo install -m 0755 -d /etc/apt/keyrings
                        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
                        sudo chmod a+r /etc/apt/keyrings/docker.gpg
                        echo \
                          "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
                          "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
                          sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
                        sudo apt-get update
                        sudo apt-get install -y docker-ce docker-ce-cli containerd.io
                    fi
                    ;;
                dnf|yum)
                    if [ "$UPDATE_MODE" = true ] && command_exists docker; then
                        sudo "$PKG_MANAGER" upgrade -y docker
                    else
                        sudo "$PKG_MANAGER" install -y docker
                        sudo systemctl start docker
                        sudo systemctl enable docker
                    fi
                    ;;
                pacman)
                    if [ "$UPDATE_MODE" = true ] && command_exists docker; then
                        sudo pacman -Syu --noconfirm docker
                    else
                        sudo pacman -S --noconfirm docker
                        sudo systemctl start docker
                        sudo systemctl enable docker
                    fi
                    ;;
                *)
                    log_error "Automatic Docker installation not supported for $PKG_MANAGER"
                    log_warn "Please install Docker manually: https://docs.docker.com/engine/install/"
                    ;;
            esac
            
            # Add current user to docker group (only if not already added)
            if [ "$OS" = "linux" ] && ! groups "$USER" | grep -q docker; then
                sudo usermod -aG docker "$USER"
                log_warn "You may need to log out and back in for Docker group membership to take effect"
            fi
            ;;
    esac
    
    if [ "$UPDATE_MODE" = true ]; then
        log_info "Docker updated successfully ($(docker --version))"
    else
        log_info "Docker installed successfully"
    fi
}

# Install Go tools for protobuf generation
install_go_tools() {
    if ! command_exists go; then
        log_error "Go is not installed. Please run install_go first."
        exit 1
    fi

    log_info "Installing Go tools for protobuf generation (versions from go.mod)..."

    # Extract exact versions from go.mod using go list -m
    # This ensures we always install the versions pinned in go.mod
    local protoc_go_grpc_version=$(go list -m google.golang.org/grpc/cmd/protoc-gen-go-grpc 2>/dev/null | awk '{print $2}')
    local buf_version=$(go list -m github.com/bufbuild/buf 2>/dev/null | awk '{print $2}')
    local cosmos_proto_version=$(go list -m github.com/cosmos/cosmos-proto 2>/dev/null | awk '{print $2}')
    local gogoproto_version=$(go list -m github.com/cosmos/gogoproto 2>/dev/null | awk '{print $2}')
    local grpc_gateway_v2_version=$(go list -m github.com/grpc-ecosystem/grpc-gateway/v2 2>/dev/null | awk '{print $2}')
    local goimports_version=$(go list -m golang.org/x/tools 2>/dev/null | awk '{print $2}')
    local protoc_go_version=$(go list -m google.golang.org/protobuf 2>/dev/null | awk '{print $2}')

    # Array of tools to install - Format: "binary_name|package_path|version"
    local tools=(
        "buf|github.com/bufbuild/buf/cmd/buf|${buf_version}"
        "protoc-gen-go-pulsar|github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar|${cosmos_proto_version}"
        "protoc-gen-gocosmos|github.com/cosmos/gogoproto/protoc-gen-gocosmos|${gogoproto_version}"
        "protoc-gen-openapiv2|github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2|${grpc_gateway_v2_version}"
        "goimports|golang.org/x/tools/cmd/goimports|${goimports_version}"
        "protoc-gen-go-grpc|google.golang.org/grpc/cmd/protoc-gen-go-grpc|${protoc_go_grpc_version}"
        "protoc-gen-go|google.golang.org/protobuf/cmd/protoc-gen-go|${protoc_go_version}"
    )

    for tool_entry in "${tools[@]}"; do
        local tool_name="${tool_entry%%|*}"
        local tool_package="${tool_entry#*|}"
        tool_package="${tool_package%|*}"
        local tool_version="${tool_entry##*|}"

        if [ -z "$tool_version" ]; then
            log_warn "$tool_name version not found in go.mod, skipping"
            continue
        fi

        log_info "Installing $tool_name@$tool_version..."
        go install "${tool_package}@${tool_version}"
        if [ $? -eq 0 ]; then
            log_info "$tool_name@$tool_version installed successfully"
        else
            log_error "Failed to install $tool_name@$tool_version"
            exit 1
        fi
    done

    # Verify GOPATH/bin is in PATH
    local gopath_bin="${GOPATH:-$HOME/go}/bin"
    if [[ ":$PATH:" != *":$gopath_bin:"* ]]; then
        log_warn "Warning: $gopath_bin is not in your PATH"
        log_warn "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"\$PATH:$gopath_bin\""
    fi

    log_info "All Go tools installed successfully (matching go.mod versions)"
}

# Install build essentials (gcc, make, git)
install_build_essentials() {
    log_info "Installing build essentials..."
    
    case "$OS" in
        macos)
            # Check if Xcode Command Line Tools are installed
            if ! xcode-select -p &>/dev/null; then
                log_info "Installing Xcode Command Line Tools..."
                xcode-select --install
                log_warn "Please complete the Xcode Command Line Tools installation and re-run this script"
                exit 0
            else
                log_info "Xcode Command Line Tools already installed"
            fi

            # Install make if needed
            if ! command_exists make; then
                brew install make
            fi
            ;;
        linux)
            case "$PKG_MANAGER" in
                apt)
                    sudo apt-get update
                    sudo apt-get install -y build-essential git
                    ;;
                dnf|yum)
                    sudo "$PKG_MANAGER" groupinstall -y "Development Tools"
                    sudo "$PKG_MANAGER" install -y git
                    ;;
                pacman)
                    sudo pacman -S --noconfirm base-devel git
                    ;;
                zypper)
                    sudo zypper install -y -t pattern devel_basis
                    sudo zypper install -y git
                    ;;
            esac
            ;;
    esac

    log_info "Build essentials installed successfully"
}

# Main installation flow
main() {
    log_info "Starting zrchain dependencies installation..."
    echo
    
    # Detect OS
    detect_os
    echo
    
    # Detect package manager for Linux
    if [ "$OS" = "linux" ]; then
        detect_linux_package_manager
        echo
    fi
    
    # Check for Homebrew on macOS
    if [ "$OS" = "macos" ] && ! command_exists brew; then
        install_homebrew
        echo
    fi
    
    # Set PKG_MANAGER for macOS (for consistency in install_package function)
    if [ "$OS" = "macos" ]; then
        PKG_MANAGER="brew"
    fi
    
    # Install dependencies
    install_build_essentials
    echo
    
    install_go
    echo
    
    install_go_tools
    echo
    
    install_just
    echo
    
    install_uv
    echo
    
    install_docker
    echo

    log_info "All dependencies installed successfully!"
    echo
    log_info "Next steps:"
    echo "  1. If you installed Docker on Linux, log out and back in for group permissions"
    echo "  2. Verify installation: go version && docker --version && just --version"
    echo "  3. Verify Go tools: buf --version && protoc-gen-gocosmos --version"
    echo "  4. Run 'just' to see available commands or 'make proto-all' to generate Protobuf files"
    echo "  5. Run 'make build' or 'just build' to build the project"
    echo
    log_info "For more information, see the README.md and CLAUDE.md files"
}

# Run main function
main

