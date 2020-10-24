## Manual Deployment

```bash
cd build

# build
heroku container:push --app=$APP web --context-path=../

# release
heroku container:release --app=$APP web
```