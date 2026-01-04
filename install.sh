#!/bin/bash

# ReboloLang Framework Installation Script
# Installs Rebolo CLI globally

set -e

echo "🚀 Installing ReboloLang Framework..."
echo "Inspired by Rebolo, Barranquilla, Colombia 🇨🇴"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first: https://golang.org/dl/"
    exit 1
fi

# Check if Bun is installed
if ! command -v bun &> /dev/null; then
    echo "⚠️  Bun.js is not installed. Installing Bun..."
    curl -fsSL https://bun.sh/install | bash
    export PATH="$HOME/.bun/bin:$PATH"
fi

# Install Rebolo CLI
echo "📦 Installing Rebolo CLI..."
go install github.com/Palaciodiego008/rebololang/cmd/rebolo@latest

# Verify installation
if command -v rebolo &> /dev/null; then
    echo "✅ ReboloLang installed successfully!"
    echo ""
    echo "🎉 Get started:"
    echo "   rebolo new myapp"
    echo "   cd myapp"
    echo "   rebolo dev"
    echo ""
    echo "🔧 Generate resources:"
    echo "   rebolo generate resource posts title:string content:text"
    echo ""
    echo "🗃️ Database operations:"
    echo "   rebolo db migrate"
    echo ""
    echo "📚 Documentation: https://github.com/Palaciodiego008/rebololang"
else
    echo "❌ Installation failed. Make sure $GOPATH/bin is in your PATH"
    exit 1
fi
