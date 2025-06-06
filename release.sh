#!/bin/bash
# Automated release script for Go-OPML
# This script handles the complete release process including building binaries for all platforms

set -e  # Exit on any error

echo "🚀 Go-OPML Automated Release Script"
echo "===================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "cmd/go-opml" ]; then
    log_error "This script must be run from the Go-OPML project root directory"
    exit 1
fi

# Get current git status
CURRENT_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v1.0.2")
BRANCH=$(git branch --show-current)

log_info "Current branch: $BRANCH"
log_info "Latest tag: $LATEST_TAG"

# Determine next version
if [ "$LATEST_TAG" = "v1.0.2" ]; then
    NEXT_VERSION="v1.0.3"
elif [ "$LATEST_TAG" = "v1.0.3" ]; then
    NEXT_VERSION="v1.0.4"
else
    # Extract version numbers and increment patch
    VERSION_NUMBER=$(echo $LATEST_TAG | sed 's/v//')
    IFS='.' read -r major minor patch <<< "$VERSION_NUMBER"
    NEXT_VERSION="v$major.$minor.$((patch + 1))"
fi

echo ""
log_info "Next version will be: $NEXT_VERSION"
echo ""

# Ask for confirmation
read -p "Do you want to proceed with release $NEXT_VERSION? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    log_warning "Release cancelled by user"
    exit 0
fi

echo ""
log_info "Starting release process for $NEXT_VERSION..."
echo ""

# Step 1: Clean workspace
log_info "1. Cleaning workspace..."
rm -rf build/
rm -f coverage.out coverage.html
mkdir -p build/

# Step 2: Check dependencies
log_info "2. Checking dependencies..."
go mod tidy
go mod verify

# Step 3: Check for dependency updates
log_info "3. Checking for dependency updates..."
UPDATES=$(go list -u -m all | grep '\[.*\]' | wc -l)
if [ $UPDATES -gt 0 ]; then
    log_warning "$UPDATES dependency updates available (will be noted in changelog)"
    go list -u -m all | grep '\[.*\]' | head -5
else
    log_success "All dependencies are up to date"
fi

# Step 4: Build binaries for all platforms
log_info "4. Building binaries for all platforms..."

# Build flags for optimization (corrected syntax)
BUILD_FLAGS="-ldflags=-s -w"

# Linux 64-bit
log_info "   Building for Linux 64-bit..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/Go-OPML-linux cmd/go-opml/main.go

# Windows 64-bit
log_info "   Building for Windows 64-bit..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/Go-OPML.exe cmd/go-opml/main.go

# macOS Intel 64-bit
log_info "   Building for macOS Intel..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/Go-OPML-mac-intel cmd/go-opml/main.go

# macOS Apple Silicon (ARM64)
log_info "   Building for macOS Apple Silicon..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/Go-OPML-mac-arm64 cmd/go-opml/main.go

# Build for current platform (for testing)
log_info "   Building for current platform..."
go build -ldflags="-s -w" -o build/Go-OPML cmd/go-opml/main.go

log_success "All binaries built successfully!"
echo ""
log_info "📦 Built binaries:"
ls -la build/
echo ""

# Step 5: Test the binary
log_info "5. Testing the built binary..."
if [ -f "examples/sample.opml" ]; then
    if ./build/Go-OPML -input examples/sample.opml -output test-output.json 2>/dev/null; then
        log_success "Binary test successful!"
        rm -f test-output.json
    else
        log_error "Binary test failed!"
        exit 1
    fi
else
    log_warning "No sample.opml found, skipping binary test"
fi

# Step 6: Calculate binary sizes and checksums
log_info "6. Generating checksums..."
cd build
sha256sum * > checksums.txt
cd ..
log_success "Checksums generated in build/checksums.txt"

# Step 7: Update README with new version (if needed)
log_info "7. Updating README version references..."
if grep -q "v1.0.3 (Ready)" README.md; then
    sed -i "s/v1.0.3 (Ready)/$NEXT_VERSION (Released)/g" README.md
    log_success "README updated with new version"
fi

# Step 8: Commit changes
log_info "8. Committing changes..."
git add .
git commit -m "Release $NEXT_VERSION: Enhanced build pipeline, dependency updates, modernized website

- Updated dependencies for better security and performance
- Enhanced build pipeline with intelligent version detection  
- Fixed sample.opml by removing broken RSS feed
- Modernized website with improved UX and design
- Added automated release script for all platforms
- Cross-platform binaries: Linux, Windows, macOS (Intel & ARM64)

Binaries:
- Go-OPML-linux (Linux 64-bit)
- Go-OPML.exe (Windows 64-bit) 
- Go-OPML-mac-intel (macOS Intel 64-bit)
- Go-OPML-mac-arm64 (macOS Apple Silicon)"

# Step 9: Create and push tag
log_info "9. Creating and pushing tag $NEXT_VERSION..."
git tag $NEXT_VERSION
<<<<<<< HEAD
git push origin $BRANCH
=======
git push origin main
>>>>>>> 74b842961efc81c74660e67b987e41ff88f1e1d3
git push origin $NEXT_VERSION

log_success "Tag $NEXT_VERSION created and pushed!"

# Step 10: Create release notes
log_info "10. Generating release notes..."
cat > release-notes.md << EOF
# Go-OPML $NEXT_VERSION

## 🚀 What's New

- **Enhanced Build Pipeline**: Automated release script with intelligent version detection
- **Dependency Updates**: Updated dependencies for better security and performance  
- **Sample File Fix**: Removed broken RSS feed from sample.opml
- **Modernized Website**: Improved UX, modern design, and better documentation
- **Cross-Platform Binaries**: Native binaries for all major platforms

## 📦 Available Downloads

### Windows
- **Go-OPML.exe** - Windows 64-bit

### macOS  
- **Go-OPML-mac-intel** - macOS Intel 64-bit
- **Go-OPML-mac-arm64** - macOS Apple Silicon (M1/M2)

### Linux
- **Go-OPML-linux** - Linux 64-bit

## 🔧 Quick Start

1. Download the appropriate binary for your platform
2. Make it executable: \`chmod +x Go-OPML-*\`
3. Run: \`./Go-OPML-* -input your_file.opml -output result.json\`

## 📋 Checksums

SHA256 checksums are provided in \`checksums.txt\` for verification.

## 🔗 Links

- **Documentation**: [go-opml.madi.se](https://go-opml.madi.se)
- **Go Package**: [pkg.go.dev/github.com/jadmadi/Go-OPML](https://pkg.go.dev/github.com/jadmadi/Go-OPML)
- **Source Code**: [github.com/jadmadi/Go-OPML](https://github.com/jadmadi/Go-OPML)

---

**Full Changelog**: https://github.com/jadmadi/Go-OPML/compare/$LATEST_TAG...$NEXT_VERSION
EOF

log_success "Release notes generated in release-notes.md"

# Step 11: Instructions for GitHub release
echo ""
log_success "🎉 Release $NEXT_VERSION is ready!"
echo ""
log_info "📋 Next steps to complete the GitHub release:"
echo ""
echo "1. Go to: https://github.com/jadmadi/Go-OPML/releases/new"
echo "2. Choose tag: $NEXT_VERSION"
echo "3. Release title: Go-OPML $NEXT_VERSION"
echo "4. Copy content from release-notes.md as description"
echo "5. Upload binaries by dragging files from build/ directory:"
echo "   📁 build/Go-OPML.exe (Windows 64-bit)"
echo "   📁 build/Go-OPML-mac-intel (macOS Intel)" 
echo "   📁 build/Go-OPML-mac-arm64 (macOS Apple Silicon)"
echo "   📁 build/Go-OPML-linux (Linux 64-bit)"
echo "   📁 build/checksums.txt (SHA256 checksums)"
echo "6. Click 'Publish release'"
echo ""
echo "💡 Tip: You can drag and drop all files at once into the GitHub release page"
echo ""

# Show file sizes and checksums for reference
echo "📊 Files ready for upload:"
ls -lah build/
echo ""
echo "📋 SHA256 Checksums:"
cat build/checksums.txt
echo ""

log_success "Automated release process completed successfully! 🚀"
log_warning "📤 Don't forget to upload the binaries to complete the GitHub release."

# Optional: Open GitHub releases page
if command -v xdg-open &> /dev/null; then
    read -p "Open GitHub releases page in browser? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        xdg-open "https://github.com/jadmadi/Go-OPML/releases/new?tag=$NEXT_VERSION"
    fi
elif command -v open &> /dev/null; then
    read -p "Open GitHub releases page in browser? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        open "https://github.com/jadmadi/Go-OPML/releases/new?tag=$NEXT_VERSION"
    fi
fi
