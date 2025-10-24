#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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
    
    if command_exists "$package"; then
        log_info "$package_name is already installed"
        return 0
    fi
    
    log_info "Installing $package_name..."
    
    case "$OS-$PKG_MANAGER" in
        macos-*)
            brew install "$package"
            ;;
        linux-apt)
            sudo apt-get update
            sudo apt-get install -y "$package"
            ;;
        linux-dnf)
            sudo dnf install -y "$package"
            ;;
        linux-yum)
            sudo yum install -y "$package"
            ;;
        linux-pacman)
            sudo pacman -S --noconfirm "$package"
            ;;
        linux-zypper)
            sudo zypper install -y "$package"
            ;;
        *)
            log_error "Unsupported OS/package manager combination: $OS-$PKG_MANAGER"
            exit 1
            ;;
    esac
    
    log_info "$package_name installed successfully"
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
            log_warn "On macOS, Docker Desktop needs to be installed manually."
            log_warn "Please visit: https://www.docker.com/products/docker-desktop"
            log_warn "Or use: brew install --cask docker"
            read -p "Press Enter to continue after installing Docker Desktop, or Ctrl+C to exit..."
            ;;
        linux)
            case "$PKG_MANAGER" in
                apt)
                    sudo apt-get update
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
                    ;;
                dnf|yum)
                    sudo "$PKG_MANAGER" install -y docker
                    sudo systemctl start docker
                    sudo systemctl enable docker
                    ;;
                pacman)
                    sudo pacman -S --noconfirm docker
                    sudo systemctl start docker
                    sudo systemctl enable docker
                    ;;
                *)
                    log_error "Automatic Docker installation not supported for $PKG_MANAGER"
                    log_warn "Please install Docker manually: https://docs.docker.com/engine/install/"
                    ;;
            esac
            
            # Add current user to docker group
            if [ "$OS" = "linux" ]; then
                sudo usermod -aG docker "$USER"
                log_warn "You may need to log out and back in for Docker group membership to take effect"
            fi
            ;;
    esac
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
    
    install_uv
    echo
    
    install_docker
    echo
    
    log_info "All dependencies installed successfully!"
    echo
    log_info "Next steps:"
    echo "  1. If you installed Docker on Linux, log out and back in for group permissions"
    echo "  2. Verify installation: go version && docker --version && uv --version"
    echo "  3. Run 'make proto-all' to generate Protobuf files (uses Docker)"
    echo "  4. Run 'make build' to build the project"
    echo
    log_info "For more information, see the README.md and CLAUDE.md files"
}

# Run main function
main

