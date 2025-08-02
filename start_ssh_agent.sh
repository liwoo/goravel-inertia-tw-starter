#!/bin/bash

# Script to start SSH agent and add keys

# Check if SSH agent is already running
if [ -z "$SSH_AUTH_SOCK" ]; then
    echo "Starting SSH agent..."
    eval "$(ssh-agent -s)"
else
    echo "SSH agent is already running"
fi

# Add SSH keys
if [ -f ~/.ssh/id_ed25519 ]; then
    echo "Adding id_ed25519 key..."
    ssh-add ~/.ssh/id_ed25519
else
    echo "SSH key ~/.ssh/id_ed25519 not found"
fi

# List added keys
echo "Current SSH keys:"
ssh-add -l