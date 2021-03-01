mkdir -p .git/hooks
cp misc/git/commit-msg .git/hooks/commit-msg
chmod +x .git/hooks/commit-msg
#git config core.hooksPath .git/hooks
